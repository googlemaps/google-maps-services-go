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

// Package internal provides support for the maps packages.
//
// Users should not import this package directly.
package internal

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/context"
)

type contextKey struct{}

// WithContext is the internal constructor for mapsContext.
func WithContext(parent context.Context, apiKey string, c *http.Client, baseURL, roadsBaseURL string) context.Context {
	if c == nil {
		panic("nil *http.Client passed to WithContext")
	}
	if apiKey == "" {
		panic("empty API Key passed to WithContext")
	}
	if !strings.HasPrefix(apiKey, "AIza") {
		panic("invalid API Key passed to WithContext")
	}
	if baseURL == "" {
		panic("invalid base URL passed to WithContext")
	}
	return context.WithValue(parent, contextKey{}, &mapsContext{
		APIKey:       apiKey,
		HTTPClient:   c,
		BaseURL:      baseURL,
		RoadsBaseURL: roadsBaseURL,
	})
}

const userAgent = "GoogleGeoApiClientGo/0.1"

type mapsContext struct {
	APIKey       string
	HTTPClient   *http.Client
	BaseURL      string
	RoadsBaseURL string

	mu  sync.Mutex             // guards svc
	svc map[string]interface{} // e.g. "storage" => *rawStorage.Service
}

// Service returns the result of the fill function if it's never been
// called before for the given name (which is assumed to be an API
// service name, like "directions"). If it has already been cached, the fill
// func is not run.
// It's safe for concurrent use by multiple goroutines.
func Service(ctx context.Context, name string, fill func(*http.Client) interface{}) interface{} {
	return mc(ctx).service(name, fill)
}

func (c *mapsContext) service(name string, fill func(*http.Client) interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.svc == nil {
		c.svc = make(map[string]interface{})
	} else if v, ok := c.svc[name]; ok {
		return v
	}
	v := fill(c.HTTPClient)
	c.svc[name] = v
	return v
}

// Transport is an http.RoundTripper that appends
// Google Cloud client's user-agent to the original
// request's user-agent header.
type Transport struct {
	// Base represents the actual http.RoundTripper
	// the requests will be delegated to.
	Base http.RoundTripper
}

// RoundTrip appends a user-agent to the existing user-agent
// header and delegates the request to the base http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
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

// APIKey retrieval for mapsContext
func APIKey(ctx context.Context) string {
	return mc(ctx).APIKey
}

// HTTPClient retrieval for mapsContext
func HTTPClient(ctx context.Context) *http.Client {
	return mc(ctx).HTTPClient
}

// BaseURL retrieval for mapsContext
func BaseURL(ctx context.Context) string {
	return mc(ctx).BaseURL
}

// RoadsBaseURL retrieval for mapsContext
func RoadsBaseURL(ctx context.Context) string {
	return mc(ctx).RoadsBaseURL
}

// mc returns the internal *mapsContext (cc) state for a context.Context.
// It panics if the user did it wrong.
func mc(ctx context.Context) *mapsContext {
	if c, ok := ctx.Value(contextKey{}).(*mapsContext); ok {
		return c
	}
	panic("invalid context.Context type; it should be created with maps.NewContext")
}
