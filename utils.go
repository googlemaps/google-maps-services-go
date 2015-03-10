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

// Package maps contains TODO(brettmorgan)
//
// More information about Google Distance Matrix API is available on
// https://developers.google.com/maps/documentation/distancematrix/
package maps // import "google.golang.org/maps"

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func rawService(ctx context.Context) *http.Client {
	return internal.Service(ctx, "directions", func(hc *http.Client) interface{} {
		// TODO(brettmorgan): Introduce a rate limiting wrapper for hc here.
		return hc
	}).(*http.Client)
}
