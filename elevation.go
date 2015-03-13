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

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

// Get makes an Elevation API request
func (r *ElevationRequest) Get(ctx context.Context) ([]ElevationResult, error) {
	var response elevationResponse

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/elevation/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("key", internal.APIKey(ctx))

	if len(r.Path) == 0 && len(r.Locations) == 0 {
		return nil, errors.New("elevation: Provide either Path or Locations")
	}

	if len(r.Path) > 0 {
		// Sampled path request
		if r.Samples == 0 {
			return nil, errors.New("elevation: Sampled Path Request requires Samples to be specifed")
		}
		q.Set("path", fmt.Sprintf("enc:%s", Encode(r.Path).Points))
		q.Set("samples", strconv.Itoa(r.Samples))
	}

	if len(r.Locations) > 0 {
		q.Set("locations", fmt.Sprintf("enc:%s", Encode(r.Locations).Points))
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
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf("distancematrix: %s - %s", response.Status, response.ErrorMessage)
	}

	return response.Results, nil
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

type elevationResponse struct {
	// Results is the Elevation results array
	Results []ElevationResult `json:"results"`

	// Status indicating if this request was successful
	Status string `json:"status"`
	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
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
