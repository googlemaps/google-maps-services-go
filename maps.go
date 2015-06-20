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

// Package maps contains Google Maps API Web Services related types
// and common functions.
package maps // import "google.golang.org/maps"

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

const baseURL = "https://maps.googleapis.com/"

// NewContext returns a new context that uses the provided http.Client.
// It mutates the client's original Transport to append the cloud
// package's user-agent to the outgoing requests.
// You can obtain the API Key from the Google Developers Console,
// https://console.developers.google.com.
func NewContext(apiKey string, c *http.Client) context.Context {
	return newContextWithBaseURL(apiKey, c, baseURL)
}

// newContextWithBaseURL returns a new context in a similar way NewContext does,
// but with a specified baseURL. Useful for testing.
func newContextWithBaseURL(apiKey string, c *http.Client, baseURL string) context.Context {
	if c == nil {
		panic("invalid nil *http.Client passed to NewContext")
	}
	return internalWithConstructor(context.Background(), apiKey, c, baseURL)
}

// WithContext returns a new context in a similar way NewContext does,
// but initiates the new context with the specified parent.
func WithContext(parent context.Context, apiKey string, c *http.Client) context.Context {
	return internalWithConstructor(parent, apiKey, c, baseURL)
}

func internalWithConstructor(parent context.Context, apiKey string, c *http.Client, baseURL string) context.Context {
	// TODO(bradfitz): delete internal.Transport. It's too wrappy for what it does.
	// Do User-Agent some other way.
	if _, ok := c.Transport.(*internal.Transport); !ok {
		c.Transport = &internal.Transport{Base: c.Transport}
	}
	return internal.WithContext(parent, apiKey, c, baseURL)
}
