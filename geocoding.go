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

// More information about Google Geocoding API is available on
// https://developers.google.com/maps/documentation/geocoding

package maps

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

var geocodingAPI = &apiConfig{
	host:             "https://maps.googleapis.com",
	path:             "/maps/api/geocode/json",
	acceptsClientID:  true,
	acceptsSignature: false,
}

// Geocode makes a Geocoding API request
func (c *Client) Geocode(ctx context.Context, r *GeocodingRequest) ([]GeocodingResult, error) {
	if r.Address == "" && len(r.Components) == 0 && r.LatLng == nil {
		return nil, errors.New("maps: address, components and LatLng are all missing")
	}

	var response struct {
		Results []GeocodingResult `json:"results"`
		commonResponse
	}

	if err := c.getJSON(ctx, geocodingAPI, r, &response); err != nil {
		return nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// ReverseGeocode makes a Reverse Geocoding API request
func (c *Client) ReverseGeocode(ctx context.Context, r *GeocodingRequest) ([]GeocodingResult, error) {
	// Since Geocode() does not allow a nil LatLng, whereas it is allowed here
	if r.LatLng == nil && r.PlaceID == "" {
		return nil, errors.New("maps: LatLng and PlaceID are both missing")
	}

	var response struct {
		Results []GeocodingResult `json:"results"`
		commonResponse
	}

	if err := c.getJSON(ctx, geocodingAPI, r, &response); err != nil {
		return nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, err
	}

	return response.Results, nil

}

func (r *GeocodingRequest) params() url.Values {
	q := make(url.Values)

	for k, v := range r.Custom {
		q[k] = v
	}

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
	if r.PlaceID != "" {
		q.Set("place_id", r.PlaceID)
	}
	if r.Language != "" {
		q.Set("language", r.Language)
	}

	return q
}

// GeocodeAccuracy is the type of a location result from the Geocoding API.
type GeocodeAccuracy string

const (
	// GeocodeAccuracyRooftop restricts the results to addresses for which Google has
	// location information accurate down to street address precision.
	GeocodeAccuracyRooftop = GeocodeAccuracy("ROOFTOP")
	// GeocodeAccuracyRangeInterpolated restricts the results to those that reflect an
	// approximation interpolated between two precise points.
	GeocodeAccuracyRangeInterpolated = GeocodeAccuracy("RANGE_INTERPOLATED")
	// GeocodeAccuracyGeometricCenter restricts the results to geometric centers of a
	// location such as a polyline or polygon.
	GeocodeAccuracyGeometricCenter = GeocodeAccuracy("GEOMETRIC_CENTER")
	// GeocodeAccuracyApproximate restricts the results to those that are characterized
	// as approximate.
	GeocodeAccuracyApproximate = GeocodeAccuracy("APPROXIMATE")
)

// GeocodingRequest is the request structure for Geocoding API
type GeocodingRequest struct {
	// Geocoding fields

	// Address is the street address that you want to geocode, in the format used by
	// the national postal service of the country concerned.
	Address string
	// Components is a component filter for which you wish to obtain a geocode. Either
	// Address or Components is required in a geocoding request. For more detail on
	// Component Filtering please see
	// https://developers.google.com/maps/documentation/geocoding/intro#ComponentFiltering
	Components map[Component]string
	// Bounds is the bounding box of the viewport within which to bias geocode results
	// more prominently. Optional.
	Bounds *LatLngBounds
	// Region is the region code, specified as a ccTLD two-character value. Optional.
	Region string

	// Reverse geocoding fields

	// LatLng is the textual latitude/longitude value for which you wish to obtain the
	// closest, human-readable address. Either LatLng or PlaceID is required for
	// Reverse Geocoding.
	LatLng *LatLng
	// ResultType is an array of one or more address types. Optional.
	ResultType []string
	// LocationType is an array of one or more geocoding accuracy types. Optional.
	LocationType []GeocodeAccuracy
	// PlaceID is a string which contains the place_id, which can be used for reverse
	// geocoding requests. Either LatLng or PlaceID is required for Reverse Geocoding.
	PlaceID string

	// Language is the language in which to return results. Optional.
	Language string

	// Custom allows passing through custom parameters to the Geocoding back end.
	// Use with caution. For more detail on why this is required, please see
	// https://googlegeodevelopers.blogspot.com/2016/11/address-geocoding-in-google-maps-apis.html
	Custom url.Values
}

// GeocodingResult is a single geocoded address
type GeocodingResult struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          AddressGeometry    `json:"geometry"`
	Types             []string           `json:"types"`
	PlaceID           string             `json:"place_id"`

	// PartialMatch indicates that the geocoder did not return an exact match for
	// the original request, though it was able to match part of the requested address.
	// You may wish to examine the original request for misspellings and/or an incomplete address.
	// Partial matches most often occur for street addresses that do not exist within the
	// locality you pass in the request.
	// Partial matches may also be returned when a request matches two or more locations in
	// the same locality. For example, "21 Henr St, Bristol, UK" will return a partial match
	// for both Henry Street and Henrietta Street.
	// Note that if a request includes a misspelled address component, the geocoding service may
	// suggest an alternative address.
	// Suggestions triggered in this way will also be marked as a partial match.
	PartialMatch bool `json:"partial_match"`

	// PlusCode (see https://en.wikipedia.org/wiki/Open_Location_Code and https://plus.codes/)
	// is an encoded location reference, derived from latitude and longitude coordinates,
	// that represents an area: 1/8000th of a degree by 1/8000th of a degree (about 14m x 14m at the equator)
	// or smaller.
	//
	// Plus codes can be used as a replacement for street addresses in places where they do not exist
	// (where buildings are not numbered or streets are not named).
	// The plus code is formatted as a global code and a compound code:
	// Typically, both the global code and compound code are returned.
	// However, if the result is in a remote location (for example, an ocean or desert)
	// only the global code may be returned.
	PlusCode AddressPlusCode `json:"plus_code"`
}

// AddressPlusCode (see https://en.wikipedia.org/wiki/Open_Location_Code and https://plus.codes/)
// is an encoded location reference, derived from latitude and longitude coordinates,
// that represents an area: 1/8000th of a degree by 1/8000th of a degree (about 14m x 14m at the equator)
// or smaller.
//
// Plus codes can be used as a replacement for street addresses in places where they do not exist
// (where buildings are not numbered or streets are not named).
// The plus code is formatted as a global code and a compound code:
// Typically, both the global code and compound code are returned.
// However, if the result is in a remote location (for example, an ocean or desert)
// only the global code may be returned.
type AddressPlusCode struct {
	// GlobalCode is a 4 character area code and 6 character or longer local code (849VCWC8+R9).
	GlobalCode string `json:"global_code"`
	// CompoundCode is a 6 character or longer local code with an explicit location (CWC8+R9, Mountain View, CA, USA).
	CompoundCode string `json:"compound_code"`
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
	Bounds       LatLngBounds `json:"bounds"`
	Viewport     LatLngBounds `json:"viewport"`
	Types        []string     `json:"types"`
}
