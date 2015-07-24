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
	"time"

	"golang.org/x/net/context"
)

type distanceMatrixResponse struct {
	matrix *DistanceMatrixResponse
	err    error
}

// DistanceMatrix makes a Distance Matrix API request
func (c *Client) DistanceMatrix(ctx context.Context, r *DistanceMatrixRequest) (*DistanceMatrixResponse, error) {

	if len(r.Origins) == 0 {
		return nil, errors.New("distancematrix: Origins must contain at least one start address")
	}
	if len(r.Destinations) == 0 {
		return nil, errors.New("distancematrix: Destinations must contain at least one end address")
	}
	if r.DepartureTime != "" && r.ArrivalTime != "" {
		return nil, errors.New("distancematrix: must not specify both DepartureTime and ArrivalTime")
	}
	if len(r.TransitMode) != 0 && r.Mode != TravelModeTransit {
		return nil, errors.New("distancematrix: must specify mode of transit when specifying transitMode")
	}
	if r.TransitRoutingPreference != "" && r.Mode != TravelModeTransit {
		return nil, errors.New("distancematrix: must specify mode of transit when specifying transitRoutingPreference")
	}

	chResult := make(chan distanceMatrixResponse)

	go func() {
		matrix, err := c.doGetDistanceMatrix(r)
		chResult <- distanceMatrixResponse{matrix, err}
	}()

	select {
	case resp := <-chResult:
		return resp.matrix, resp.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *DistanceMatrixRequest) request(c *Client) (*http.Request, error) {
	baseURL := c.getBaseURL("https://maps.googleapis.com/")

	req, err := http.NewRequest("GET", baseURL+"/maps/api/distancematrix/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("origins", strings.Join(r.Origins, "|"))
	q.Set("destinations", strings.Join(r.Destinations, "|"))
	if r.Mode != "" {
		q.Set("mode", string(r.Mode))
	}
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	if len(r.Avoid) > 0 {
		var avoid []string
		for _, a := range r.Avoid {
			avoid = append(avoid, string(a))
		}
		q.Set("avoid", strings.Join(avoid, "|"))
	}
	if r.Units != "" {
		q.Set("units", string(r.Units))
	}
	if r.DepartureTime != "" {
		q.Set("departure_time", r.DepartureTime)
	}
	if r.ArrivalTime != "" {
		q.Set("arrival_time", r.ArrivalTime)
	}
	if len(r.TransitMode) != 0 {
		var transitMode []string
		for _, t := range r.TransitMode {
			transitMode = append(transitMode, string(t))
		}
		q.Set("transit_mode", strings.Join(transitMode, "|"))
	}
	if r.TransitRoutingPreference != "" {
		q.Set("transit_routing_preference", string(r.TransitRoutingPreference))
	}
	query, err := c.generateAuthQuery(req.URL.Path, q, true)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query
	return req, nil
}

func (c *Client) doGetDistanceMatrix(r *DistanceMatrixRequest) (*DistanceMatrixResponse, error) {
	req, err := r.request(c)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw rawDistanceMatrixResponse
	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return nil, err
	}
	if raw.Status != "OK" {
		err = errors.New("distancematrix: " + raw.Status + " - " + raw.ErrorMessage)
		return nil, err
	}

	response := &DistanceMatrixResponse{
		DestinationAddresses: raw.DestinationAddresses,
		OriginAddresses:      raw.OriginAddresses,
		Rows:                 raw.Rows,
	}

	return response, nil
}

// DistanceMatrixRequest is the request struct for Distance Matrix APi
type DistanceMatrixRequest struct {
	// Origins is a list of addresses and/or textual latitude/longitude values from which to calculate distance and time. Required.
	Origins []string
	// Destinations is a list of addresses and/or textual latitude/longitude values to which to calculate distance and time. Required.
	Destinations []string
	// Mode specifies the mode of transport to use when calculating distance. Valid values are `ModeDriving`, `ModeWalking`, `ModeBicycling`
	// and `ModeTransit`. Optional.
	Mode mode
	// Language in which to return results. Optional.
	Language string
	// Avoid introduces restrictions to the route. Valid values are `AvoidTolls`, `AvoidHighways` and `AvoidFerries`. Optional.
	Avoid avoid
	// Units Specifies the unit system to use when expressing distance as text. Valid values are `UnitsMetric` and `UnitsImperial`. Optional.
	Units units
	// DepartureTime is the desired time of departure. You can specify the time as an integer in seconds since midnight, January 1, 1970 UTC.
	// Alternatively, you can specify a value of `"now"``. Optional.
	DepartureTime string
	// ArrivalTime specifies the desired time of arrival for transit requests, in seconds since midnight, January 1, 1970 UTC. You cannot
	// specify both `DepartureTime` and `ArrivalTime`. Optional.
	ArrivalTime string
	// TransitMode specifies one or more preferred modes of transit. This parameter may only be specified for requests where the mode is
	// `transit`. Valid values are `TransitModeBus`, `TransitModeSubway`, `TransitModeTrain`, `TransitModeTram`, and `TransitModeRail`.
	// Optional.
	TransitMode []transitMode
	// TransitRoutingPreference Specifies preferences for transit requests. Valid values are `TransitRoutingPreferenceLessWalking` and
	// `TransitRoutingPreferenceFewerTransfers`. Optional.
	TransitRoutingPreference transitRoutingPreference
}

type rawDistanceMatrixResponse struct {
	// OriginAddresses contains an array of addresses as returned by the API from your original request.
	OriginAddresses []string `json:"origin_addresses"`
	// DestinationAddresses contains an array of addresses as returned by the API from your original request.
	DestinationAddresses []string `json:"destination_addresses"`
	// Rows contains an array of elements.
	Rows []DistanceMatrixElementsRow `json:"rows"`

	// Status contains the status of the request, and may contain
	// debugging information to help you track down why the Directions
	// service failed.
	// See https://developers.google.com/maps/documentation/distancematrix/#StatusCodes
	Status string `json:"status"`

	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// DistanceMatrixResponse represents a Distance Matrix API response.
type DistanceMatrixResponse struct {

	// OriginAddresses contains an array of addresses as returned by the API from your original request.
	OriginAddresses []string
	// DestinationAddresses contains an array of addresses as returned by the API from your original request.
	DestinationAddresses []string
	// Rows contains an array of elements.
	Rows []DistanceMatrixElementsRow
}

// DistanceMatrixElementsRow is a row of distance elements.
type DistanceMatrixElementsRow struct {
	Elements []*DistanceMatrixElement `json:"elements"`
}

// DistanceMatrixElement is the travel distance and time for a pair of origin and destination.
type DistanceMatrixElement struct {
	Status string `json:"status"`
	// Duration is the length of time it takes to travel this route.
	Duration time.Duration `json:"duration"`
	// Distance is the total distance of this route.
	Distance Distance `json:"distance"`
}
