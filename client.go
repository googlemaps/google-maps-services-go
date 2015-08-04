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

// More information about Google Directions API is available on
// https://developers.google.com/maps/documentation/directions/

package maps

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/googlemaps/google-maps-services-go/internal"
	"golang.org/x/net/context"
)

// Client may be used to make requests to the Google Maps WebService APIs
type Client struct {
	httpClient        *http.Client
	apiKey            string
	baseURL           string
	clientID          string
	signature         []byte
	requestsPerSecond int
	rateLimiter       chan int
}

// ClientOption is the type of constructor options for NewClient(...).
type ClientOption func(*Client) error

var defaultRequestsPerSecond = 10

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

	// Implement a bursty rate limiter.
	// Allow up to 1 second worth of requests to be made at once.
	c.rateLimiter = make(chan int, c.requestsPerSecond)
	// Prefill rateLimiter with 1 seconds worth of requests.
	for i := 0; i < c.requestsPerSecond; i++ {
		c.rateLimiter <- 1
	}
	go func() {
		// Refill rateLimiter continuously
		for range time.Tick(time.Second / time.Duration(c.requestsPerSecond)) {
			c.rateLimiter <- 1
		}
	}()

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
	return func(client *Client) error {
		client.apiKey = apiKey
		return nil
	}
}

// WithClientIDAndSignature configures a Maps API client for a Maps for Work application
// The signature is assumed to be URL modified Base64 encoded
func WithClientIDAndSignature(clientID, signature string) ClientOption {
	return func(client *Client) error {
		client.clientID = clientID
		decoded, err := base64.URLEncoding.DecodeString(signature)
		if err != nil {
			return err
		}
		client.signature = decoded
		return nil
	}
}

// WithRateLimit configures the rate limit for back end requests.
// Default is to limit to 10 requests per second.
func WithRateLimit(requestsPerSecond int) func(*Client) error {
	return func(client *Client) error {
		client.requestsPerSecond = requestsPerSecond
		return nil
	}
}

func (client *Client) httpDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	// TODO: change this code to use ctxhttp when it is released.
	// https://go-review.googlesource.com/#/c/12755/
	type httpResponse struct {
		response *http.Response
		err      error
	}

	<-client.rateLimiter
	c := make(chan httpResponse)

	go func() {
		resp, err := client.httpClient.Do(req)
		c <- httpResponse{resp, err}
	}()

	select {
	case resp := <-c:
		return resp.response, resp.err
	case <-ctx.Done():
		client.httpClient.Transport.(*transport).Base.(*http.Transport).CancelRequest(req)
		return nil, ctx.Err()
	}
}

const userAgent = "GoogleGeoApiClientGo/0.1"

// Transport is an http.RoundTripper that appends
// Google Cloud client's user-agent to the original
// request's user-agent header.
type transport struct {
	// Base represents the actual http.RoundTripper
	// the requests will be delegated to.
	Base http.RoundTripper
}

// RoundTrip appends a user-agent to the existing user-agent
// header and delegates the request to the base http.RoundTripper.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	ua := req.Header.Get("User-Agent")
	if ua == "" {
		ua = userAgent
	} else {
		ua = fmt.Sprintf("%s;%s", ua, userAgent)
	}
	req.Header.Set("User-Agent", ua)
	return t.Base.RoundTrip(req)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

func (client *Client) generateAuthQuery(path string, q url.Values, acceptClientID bool) (string, error) {
	if client.apiKey != "" {
		q.Set("key", client.apiKey)
		return q.Encode(), nil
	}
	if acceptClientID {
		query, err := internal.SignURL(path, client.clientID, client.signature, q)
		if err != nil {
			return "", err
		}
		return query, nil
	}
	return "", errors.New("maps: API Key missing")
}

func (client *Client) getBaseURL(baseURL string) string {
	if client.baseURL != "" {
		return client.baseURL
	}
	return baseURL
}
