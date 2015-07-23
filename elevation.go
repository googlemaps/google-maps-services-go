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

package maps

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
)

// GetElevation makes an Elevation API request
func (c *Client) GetElevation(ctx context.Context, r *ElevationRequest) ([]ElevationResult, error) {

	if len(r.Path) == 0 && len(r.Locations) == 0 {
		return nil, errors.New("elevation: Provide either Path or Locations")
	}

	// Sampled path request
	if len(r.Path) > 0 && r.Samples == 0 {
		return nil, errors.New("elevation: Sampled Path Request requires Samples to be specifed")
	}

	chResult := make(chan elevationResultWithError)

	go func() {
		elevations, err := c.doGetElevation(r)
		chResult <- elevationResultWithError{elevations, err}
	}()

	select {
	case resp := <-chResult:
		return resp.elevations, resp.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type elevationResultWithError struct {
	elevations []ElevationResult
	err        error
}

func (c *Client) doGetElevation(r *ElevationRequest) ([]ElevationResult, error) {
	baseURL := "https://maps.googleapis.com/"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}

	req, err := http.NewRequest("GET", baseURL+"/maps/api/elevation/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()

	if len(r.Path) > 0 {
		q.Set("path", fmt.Sprintf("enc:%s", Encode(r.Path)))
		q.Set("samples", strconv.Itoa(r.Samples))
	}

	if len(r.Locations) > 0 {
		q.Set("locations", fmt.Sprintf("enc:%s", Encode(r.Locations)))
	}
	query, err := c.generateAuthQuery(req.URL.Path, q, true)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query

	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response elevationResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Status != "OK" {
		err = fmt.Errorf("distancematrix: %s - %s", response.Status, response.ErrorMessage)
		return nil, err
	}

	return response.Results, nil
}

type elevationResponse struct {
	// Results is the Elevation results array
	Results []ElevationResult `json:"results"`
	// Status indicating if this request was successful
	Status string `json:"status"`
	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// ElevationRequest is the request structure for Elevation API. Either Locations or Path must be set.
type ElevationRequest struct {
	// Locations defines the location(s) on the earth from which to return elevation data.
	Locations []LatLng
	// Path defines a path on the earth for which to return elevation data.
	Path []LatLng
	// Samples specifies the number of sample points along a path for which to return elevation data. Required if Path is supplied.
	Samples int
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
