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

// More information about Google Distance Matrix API is available on
// https://developers.google.com/maps/documentation/distancematrix/

package maps // import "google.golang.org/maps"

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

// Get makes an Elevation API request
func (eReq *ElevationRequest) Get(ctx context.Context) (ElevationResponse, error) {
	var response ElevationResponse

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/elevation/json", nil)
	if err != nil {
		return response, err
	}
	q := req.URL.Query()
	q.Set("key", internal.APIKey(ctx))

	if len(eReq.Path) > 0 {
		// Sampled path request
		if eReq.Samples == 0 {
			return response, errors.New("elevation: Sampled Path Request requires Samples to be specifed")
		}
		var l []string
		for _, ll := range eReq.Path {
			l = append(l, fmt.Sprintf("%g,%g", ll.Lat, ll.Lng))
		}
		q.Set("path", strings.Join(l, "|"))
		q.Set("samples", strconv.Itoa(eReq.Samples))
	}

	if len(eReq.Locations) > 0 {
		var l []string
		for _, ll := range eReq.Locations {
			l = append(l, fmt.Sprintf("%g,%g", ll.Lat, ll.Lng))
		}
		q.Set("locations", strings.Join(l, "|"))
	}

	req.URL.RawQuery = q.Encode()

	log.Println("Request:", req)

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return err
		}
		return nil
	})
	// httpDo waits for the closure we provided to return, so it's safe to
	// read response here.
	return response, err
}

// ElevationRequest is the request structure for Elevation API
type ElevationRequest struct {
	// Locations defines the location(s) on the earth from which to return elevation data.
	Locations []LatLng
	// Path defines a path on the earth for which to return elevation data.
	Path []LatLng
	// Samples specifies the number of sample points along a path for which to return elevation data.
	Samples int
}

// ElevationResponse is the response structure for Elevation API
type ElevationResponse struct {
	// Status indicating if this request was successful
	Status string `json:"status"`
	// Results is the Elevation results array
	Results []ElevationResult `json:"results"`
}

// ElevationResult is a single elevation at a specific location
type ElevationResult struct {
	// Location is the position for which elevation data is being computed.
	Location *LatLng `json:"location"`
	// Elevation indicates the elevation of the location in meters
	Elevation float64 `json:"elevation"`
	// Resolution indicates the maximum distance between data points from which the elevation was interpolated, in meters
	Resolution float64 `json:"resolution"`
}
