// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maps

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/time/rate"
	"googlemaps.github.io/maps/internal"
)

// Client may be used to make requests to the Google Maps WebService APIs
type Client struct {
	httpClient        *http.Client
	apiKey            string
	baseURL           string
	clientID          string
	signature         []byte
	requestsPerSecond int
	rateLimiter       *rate.Limiter
	channel           string
}

// ClientOption is the type of constructor options for NewClient(...).
type ClientOption func(*Client) error

var defaultRequestsPerSecond = 50

// NewClient constructs a new Client which can make requests to the Google Maps WebService APIs.
func NewClient(options ...ClientOption) (*Client, error) {
	c := &Client{requestsPerSecond: defaultRequestsPerSecond}
	WithHTTPClient(&http.Client{})(c)
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	if c.apiKey == "" && (c.clientID == "" || len(c.signature) == 0) {
		return nil, errors.New("maps: API Key or Maps for Work credentials missing")
	}

	if c.requestsPerSecond > 0 {
		c.rateLimiter = rate.NewLimiter(rate.Limit(c.requestsPerSecond), c.requestsPerSecond)
	}

	return c, nil
}

// WithHTTPClient configures a Maps API client with a http.Client to make requests over.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		if _, ok := c.Transport.(*transport); !ok {
			t := c.Transport
			if t != nil {
				c.Transport = &transport{Base: t}
			} else {
				c.Transport = &transport{Base: http.DefaultTransport}
			}
		}
		client.httpClient = c
		return nil
	}
}

// WithAPIKey configures a Maps API client with an API Key
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

// WithBaseURL configures a Maps API client with a custom base url
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// WithChannel configures a Maps API client with a Channel
func WithChannel(channel string) ClientOption {
	return func(c *Client) error {
		c.channel = channel
		return nil
	}
}

// WithClientIDAndSignature configures a Maps API client for a Maps for Work application
// The signature is assumed to be URL modified Base64 encoded
func WithClientIDAndSignature(clientID, signature string) ClientOption {
	return func(c *Client) error {
		c.clientID = clientID
		decoded, err := base64.URLEncoding.DecodeString(signature)
		if err != nil {
			return err
		}
		c.signature = decoded
		return nil
	}
}

// WithRateLimit configures the rate limit for back end requests. Default is to
// limit to 50 requests per second. A value of zero disables rate limiting.
func WithRateLimit(requestsPerSecond int) ClientOption {
	return func(c *Client) error {
		c.requestsPerSecond = requestsPerSecond
		return nil
	}
}

type apiConfig struct {
	host            string
	path            string
	acceptsClientID bool
}

type apiRequest interface {
	params() url.Values
}

func (c *Client) awaitRateLimiter(ctx context.Context) error {
	if c.rateLimiter == nil {
		return nil
	}
	return c.rateLimiter.Wait(ctx)
}

func (c *Client) get(ctx context.Context, config *apiConfig, apiReq apiRequest) (*http.Response, error) {
	if err := c.awaitRateLimiter(ctx); err != nil {
		return nil, err
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}
	req, err := http.NewRequest("GET", host+config.path, nil)
	if err != nil {
		return nil, err
	}
	q, err := c.generateAuthQuery(config.path, apiReq.params(), config.acceptsClientID)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = q
	return ctxhttp.Do(ctx, c.httpClient, req)
}

func (c *Client) post(ctx context.Context, config *apiConfig, apiReq interface{}) (*http.Response, error) {
	if err := c.awaitRateLimiter(ctx); err != nil {
		return nil, err
	}

	host := config.host
	if c.baseURL != "" {
		host = c.baseURL
	}

	body, err := json.Marshal(apiReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", host+config.path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	q, err := c.generateAuthQuery(config.path, url.Values{}, config.acceptsClientID)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = q
	return ctxhttp.Do(ctx, c.httpClient, req)
}

func (c *Client) getJSON(ctx context.Context, config *apiConfig, apiReq apiRequest, resp interface{}) error {
	httpResp, err := c.get(ctx, config, apiReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

func (c *Client) postJSON(ctx context.Context, config *apiConfig, apiReq interface{}, resp interface{}) error {
	httpResp, err := c.post(ctx, config, apiReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	return json.NewDecoder(httpResp.Body).Decode(resp)
}

type binaryResponse struct {
	statusCode  int
	contentType string
	data        io.ReadCloser
}

func (c *Client) getBinary(ctx context.Context, config *apiConfig, apiReq apiRequest) (binaryResponse, error) {
	httpResp, err := c.get(ctx, config, apiReq)
	if err != nil {
		return binaryResponse{}, err
	}

	return binaryResponse{httpResp.StatusCode, httpResp.Header.Get("Content-Type"), httpResp.Body}, nil
}

func (c *Client) generateAuthQuery(path string, q url.Values, acceptClientID bool) (string, error) {
	if c.channel != "" {
		q.Set("channel", c.channel)
	}
	if c.apiKey != "" {
		q.Set("key", c.apiKey)
		return q.Encode(), nil
	}
	if acceptClientID {
		return internal.SignURL(path, c.clientID, c.signature, q)
	}
	return "", errors.New("maps: API Key missing")
}

// commonResponse contains the common response fields to most API calls inside
// the Google Maps APIs. This is used internally.
type commonResponse struct {
	// Status contains the status of the request, and may contain debugging
	// information to help you track down why the call failed.
	Status string `json:"status"`

	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// StatusError returns an error iff this object has a non-OK Status.
func (c *commonResponse) StatusError() error {
	if c.Status != "OK" {
		return fmt.Errorf("maps: %s - %s", c.Status, c.ErrorMessage)
	}
	return nil
}
