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
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

// Get makes a Geocoding API request
func (r *GeocodingRequest) Get(ctx context.Context) ([]GeocodingResult, error) {
	var response geocodingResponse

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/geocode/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("key", internal.APIKey(ctx))

	if r.Address != "" {
		q.Set("address", r.Address)
	}
	// TODO: fill in the rest of the params

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
		return nil, fmt.Errorf("geocoding: %s - %s", response.Status, response.ErrorMessage)
	}

	return response.Results, nil

}

type componentFilter string

const (
	// ComponentRoute matches long or short name of a route
	ComponentRoute = componentFilter("route")
	// ComponentLocality matches against both locality and sublocality types
	ComponentLocality = componentFilter("locality")
	// ComponentAdministrativeArea matches all the administrative_area levels
	ComponentAdministrativeArea = componentFilter("administrative_area")
	// ComponentPostalCode matches postal_code and postal_code_prefix
	ComponentPostalCode = componentFilter("postal_code")
	// ComponentCounty matches a country name or a two letter ISO 3166-1 country code
	ComponentCounty = componentFilter("country")
)

// GeocodingRequest is the request structure for Geocoding API
type GeocodingRequest struct {
	// Address is the street address that you want to geocode, in the format used by the national postal service of the country concerned.
	Address string
	// Components is a component filter for which you wish to obtain a geocode. Either Address or Components is required in a geocoding request.
	Components map[componentFilter]string
	// Bounds is the bounding box of the viewport within which to bias geocode results more prominently. Optional.
	Bounds Bounds
	// Language is the language in which to return results. Optional.
	Language string
	// Region is the region code, specified as a ccTLD two-character value. Optional.
	Region string
}

type geocodingResponse struct {
	Results []GeocodingResult `json:"results"`

	// Status indicating if this request was successful
	Status string `json:"status"`
	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// GeocodingResult is a single geocoded address
type GeocodingResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          AddressGeometry    `json:"geometry"`
}

// AddressComponent is a part of an address
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// AddressGeometry is the location of a an address
type AddressGeometry struct {
	Location     LatLng   `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Bounds   `json:"viewport"`
	Types        []string `json:"types"`
}
