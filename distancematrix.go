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
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

// Get makes a Distance Matrix API request
func (r *DistanceMatrixRequest) Get(ctx context.Context) (DistanceMatrixResponse, error) {
	var raw rawDistanceMatrixResponse
	var response DistanceMatrixResponse

	if len(r.Origins) == 0 {
		return response, errors.New("distancematrix: Origins must contain at least one start address")
	}
	if len(r.Destinations) == 0 {
		return response, errors.New("distancematrix: Destinations must contain at least one end address")
	}
	if r.Mode != "" && ModeDriving != r.Mode && ModeWalking != r.Mode && ModeBicycling != r.Mode && ModeTransit != r.Mode {
		return response, fmt.Errorf("distancematrix: unknown Mode: '%s'", r.Mode)
	}
	if r.Avoid != "" && r.Avoid != AvoidTolls && r.Avoid != AvoidHighways && r.Avoid != AvoidFerries {
		return response, fmt.Errorf("distancematrix: Unknown Avoid restriction '%s'", r.Avoid)
	}
	if r.Units != "" && r.Units != UnitsMetric && r.Units != UnitsImperial {
		return response, fmt.Errorf("distancematrix: Unknown Units '%s'", r.Units)
	}
	if r.TransitMode != "" && r.TransitMode != TransitModeBus && r.TransitMode != TransitModeSubway && r.TransitMode != TransitModeTrain && r.TransitMode != TransitModeTram && r.TransitMode != TransitModeRail {
		return response, fmt.Errorf("distancematrix: Unknown TransitMode '%s'", r.TransitMode)
	}
	if r.TransitRoutingPreference != "" && r.TransitRoutingPreference != TransitRoutingPreferenceLessWalking && r.TransitRoutingPreference != TransitRoutingPreferenceFewerTransfers {
		return response, fmt.Errorf("distancematrix: Unknown TransitRoutingPreference '%s'", r.TransitRoutingPreference)
	}
	if r.DepartureTime != "" && r.ArrivalTime != "" {
		return response, errors.New("distancematrix: must not specify both DepartureTime and ArrivalTime")
	}

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/distancematrix/json", nil)
	if err != nil {
		return response, err
	}
	q := req.URL.Query()
	q.Set("origins", strings.Join(r.Origins, "|"))
	q.Set("destinations", strings.Join(r.Destinations, "|"))
	q.Set("key", internal.APIKey(ctx))
	if r.Mode != "" {
		q.Set("mode", string(r.Mode))
	}
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	if r.Avoid != "" {
		q.Set("avoid", string(r.Avoid))
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
	if r.TransitMode != "" {
		q.Set("transit_mode", string(r.TransitMode))
	}
	if r.TransitRoutingPreference != "" {
		q.Set("transit_routing_preference", string(r.TransitRoutingPreference))
	}

	req.URL.RawQuery = q.Encode()

	log.Println("Request:", req)

	err = httpDo(ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response, err
	}
	if raw.Status != "OK" {
		return response, fmt.Errorf("distancematrix: %s - %s", raw.Status, raw.ErrorMessage)
	}

	response.DestinationAddresses = raw.DestinationAddresses
	response.OriginAddresses = raw.OriginAddresses
	response.Rows = raw.Rows
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
	TransitMode transitMode
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
