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
	"errors"
	"net/url"
	"strconv"

	"golang.org/x/net/context"
)

var elevationAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/elevation/json",
	acceptsClientID: true,
}

// Elevation makes an Elevation API request
func (c *Client) Elevation(ctx context.Context, r *ElevationRequest) ([]ElevationResult, error) {

	if len(r.Path) == 0 && len(r.Locations) == 0 {
		return nil, errors.New("maps: Path and Locations empty")
	}

	// Sampled path request
	if len(r.Path) > 0 && r.Samples == 0 {
		return nil, errors.New("maps: Samples empty")
	}

	var response struct {
		Results []ElevationResult `json:"results"`
		commonResponse
	}

	if err := c.getJSON(ctx, elevationAPI, r, &response); err != nil {
		return nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, err
	}

	return response.Results, nil
}

func (r *ElevationRequest) params() url.Values {
	q := make(url.Values)

	if len(r.Path) > 0 {
		q.Set("path", "enc:"+Encode(r.Path))
		q.Set("samples", strconv.Itoa(r.Samples))
	}

	if len(r.Locations) > 0 {
		q.Set("locations", "enc:"+Encode(r.Locations))
	}

	return q
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
