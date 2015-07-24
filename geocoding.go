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
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

// Geocode makes a Geocoding API request
func (c *Client) Geocode(ctx context.Context, r *GeocodingRequest) ([]GeocodingResult, error) {

	if r.Address == "" && len(r.Components) == 0 && r.LatLng == nil {
		return nil, errors.New("geocoding: You must specify at least one of Address or Components for a geocoding request, or LatLng for a reverse geocoding request")
	}

	chResult := make(chan geocodingResultWithError)

	go func() {
		results, err := c.doGetGeocoding(r)
		chResult <- geocodingResultWithError{results, err}
	}()

	select {
	case resp := <-chResult:
		return resp.Results, resp.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type geocodingResultWithError struct {
	Results []GeocodingResult
	err     error
}

func (r *GeocodingRequest) request(c *Client) (*http.Request, error) {
	baseURL := c.getBaseURL("https://maps.googleapis.com/")

	req, err := http.NewRequest("GET", baseURL+"/maps/api/geocode/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()

	if r.Address != "" {
		q.Set("address", r.Address)
	}
	var cf []string
	for c, f := range r.Components {
		cf = append(cf, string(c)+":"+f)
	}
	if len(cf) > 0 {
		q.Set("components", strings.Join(cf, "|"))
	}
	if r.Bounds != nil {
		q.Set("bounds", r.Bounds.String())
	}
	if r.Region != "" {
		q.Set("region", r.Region)
	}
	if r.LatLng != nil {
		q.Set("latlng", r.LatLng.String())
	}
	if len(r.ResultType) > 0 {
		q.Set("result_type", strings.Join(r.ResultType, "|"))
	}
	if len(r.LocationType) > 0 {
		var lt []string
		for _, l := range r.LocationType {
			lt = append(lt, string(l))
		}
		q.Set("location_type", strings.Join(lt, "|"))
	}
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	query, err := c.generateAuthQuery(req.URL.Path, q, true)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query
	return req, nil
}

func (c *Client) doGetGeocoding(r *GeocodingRequest) ([]GeocodingResult, error) {
	var response geocodingResponse

	req, err := r.request(c)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		err = errors.New("geocoding: " + response.Status + " - " + response.ErrorMessage)
		return nil, err
	}

	return response.Results, nil
}

// ComponentFilter enables filtering of returned results
type ComponentFilter string

const (
	// ComponentRoute matches long or short name of a route
	ComponentRoute = ComponentFilter("route")
	// ComponentLocality matches against both locality and sublocality types
	ComponentLocality = ComponentFilter("locality")
	// ComponentAdministrativeArea matches all the administrative_area levels
	ComponentAdministrativeArea = ComponentFilter("administrative_area")
	// ComponentPostalCode matches postal_code and postal_code_prefix
	ComponentPostalCode = ComponentFilter("postal_code")
	// ComponentCounty matches a country name or a two letter ISO 3166-1 country code
	ComponentCounty = ComponentFilter("country")
)

// LocationType restricts the results to this type
type LocationType string

const (
	// LocationTypeRooftop restricts the results to addresses for which Google has location information accurate down to street address precision
	LocationTypeRooftop = LocationType("ROOFTOP")
	// LocationTypeRangeInterpolated restricts the results to those that reflect an approximation interpolated between two precise points.
	LocationTypeRangeInterpolated = LocationType("RANGE_INTERPOLATED")
	// LocationTypeGeometricCenter restricts the results to geometric centers of a location such as a polyline or polygon.
	LocationTypeGeometricCenter = LocationType("GEOMETRIC_CENTER")
	// LocationTypeApproximate restricts the results to those that are characterized as approximate.
	LocationTypeApproximate = LocationType("APPROXIMATE")
)

// GeocodingRequest is the request structure for Geocoding API
type GeocodingRequest struct {
	// Geocoding fields

	// Address is the street address that you want to geocode, in the format used by the national postal service of the country concerned.
	Address string
	// Components is a component filter for which you wish to obtain a geocode. Either Address or Components is required in a geocoding request.
	// For more detail on Component Filtering please see https://developers.google.com/maps/documentation/geocoding/#ComponentFiltering
	Components map[ComponentFilter]string
	// Bounds is the bounding box of the viewport within which to bias geocode results more prominently. Optional.
	Bounds *LatLngBounds
	// Region is the region code, specified as a ccTLD two-character value. Optional.
	Region string

	// Reverse geocoding fields

	// LatLng is the textual latitude/longitude value for which you wish to obtain the closest, human-readable address. Required for reverse geocoding.
	LatLng *LatLng
	// ResultType is an array of one or more address types. Optional.
	ResultType []string
	// LocationType is an array of One or more location types. Optional.
	LocationType []LocationType

	// Language is the language in which to return results. Optional.
	Language string
}

type geocodingResponse struct {
	Results []GeocodingResult `json:"results"`

	// Status indicating if this request was successful
	Status string `json:"status"`
	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`

	err error
}

// GeocodingResult is a single geocoded address
type GeocodingResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          AddressGeometry    `json:"geometry"`
	Types             []string           `json:"types"`
	PlaceID           string             `json:"place_id"`
}

// AddressComponent is a part of an address
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// AddressGeometry is the location of a an address
type AddressGeometry struct {
	Location     LatLng       `json:"location"`
	LocationType string       `json:"location_type"`
	Viewport     LatLngBounds `json:"viewport"`
	Types        []string     `json:"types"`
}
