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
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

type contextKey struct{}

const userAgent = "GoogleGeoApiClientGo/0.1"

type mapsContext struct {
	APIKey          string
	HTTPClient      *http.Client
	OverrideBaseURL string

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

// mc returns the internal *mapsContext (cc) state for a context.Context.
// It panics if the user did it wrong.
func mc(ctx context.Context) *mapsContext {
	if c, ok := ctx.Value(contextKey{}).(*mapsContext); ok {
		return c
	}
	panic("invalid context.Context type; it should be created with maps.NewContext")
}
