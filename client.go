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

package maps // import "google.golang.org/maps"

import (
	"fmt"
	"net/http"
)

// Client may be used to make requests to the Google Maps WebService APIs
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type clientOption func(*Client)

// NewClient constructs a new Client which can make requests to the Google Maps WebService APIs.
// The supplied http.Client is used for making requests to the Maps WebService APIs
func NewClient(options ...clientOption) (*Client, error) {
	c := &Client{
		httpClient: http.DefaultClient,
	}
	for _, option := range options {
		option(c)
	}
	// TODO(brettmorgan): extend this to handle M4B credentials
	if c.apiKey == "" {
		return nil, fmt.Errorf("maps.Client with no API Key or credentials")
	}

	return c, nil
}

// HTTPClient configures a Maps API client with a http.Client to make requests over.
func HTTPClient(c *http.Client) func(*Client) {
	return func(client *Client) {
		client.httpClient = c
	}
}

func baseURL(url string) func(*Client) {
	return func(client *Client) {
		client.baseURL = url
	}
}

// APIKey configures a Maps API client with an API Key
func APIKey(apiKey string) func(*Client) {
	return func(client *Client) {
		client.apiKey = apiKey
	}
}

func (client *Client) httpDo(req *http.Request) (*http.Response, error) {
	return client.httpClient.Do(req)
}
