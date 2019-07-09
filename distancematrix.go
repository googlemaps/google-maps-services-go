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
	"context"
	"errors"
	"net/url"
	"strings"
	"time"
)

var distanceMatrixAPI = &apiConfig{
	host:             "https://maps.googleapis.com",
	path:             "/maps/api/distancematrix/json",
	acceptsClientID:  true,
	acceptsSignature: false,
}

// DistanceMatrix makes a Distance Matrix API request
func (c *Client) DistanceMatrix(ctx context.Context, r *DistanceMatrixRequest) (*DistanceMatrixResponse, error) {

	if len(r.Origins) == 0 {
		return nil, errors.New("maps: origins empty")
	}
	if len(r.Destinations) == 0 {
		return nil, errors.New("maps: destinations empty")
	}
	if r.DepartureTime != "" && r.ArrivalTime != "" {
		return nil, errors.New("maps: DepartureTime and ArrivalTime both specified")
	}
	if len(r.TransitMode) != 0 && r.Mode != TravelModeTransit {
		return nil, errors.New("maps: TransitMode specified while Mode != TravelModeTransit")
	}
	if r.TransitRoutingPreference != "" && r.Mode != TravelModeTransit {
		return nil, errors.New("maps: mode of transit '" + string(r.Mode) + "' invalid for TransitRoutingPreference")
	}
	if r.Mode == TravelModeTransit && r.TrafficModel != "" {
		return nil, errors.New("maps: cannot specify transit mode and traffic model together")
	}

	var response struct {
		commonResponse
		DistanceMatrixResponse
	}

	if err := c.getJSON(ctx, distanceMatrixAPI, r, &response); err != nil {
		return nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, err
	}

	return &response.DistanceMatrixResponse, nil
}

func (r *DistanceMatrixRequest) params() url.Values {
	q := make(url.Values)
	q.Set("origins", strings.Join(r.Origins, "|"))
	q.Set("destinations", strings.Join(r.Destinations, "|"))
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
	if r.TrafficModel != "" {
		q.Set("traffic_model", string(r.TrafficModel))
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
	return q
}

// DistanceMatrixRequest is the request struct for Distance Matrix APi
type DistanceMatrixRequest struct {
	// Origins is a list of addresses and/or textual latitude/longitude values
	// from which to calculate distance and time. Required.
	Origins []string
	// Destinations is a list of addresses and/or textual latitude/longitude values
	// to which to calculate distance and time. Required.
	Destinations []string
	// Mode specifies the mode of transport to use when calculating distance.
	// Valid values are `ModeDriving`, `ModeWalking`, `ModeBicycling`
	// and `ModeTransit`. Optional.
	Mode Mode
	// Language in which to return results. Optional.
	Language string
	// Avoid introduces restrictions to the route. Valid values are `AvoidTolls`,
	// `AvoidHighways` and `AvoidFerries`. Optional.
	Avoid Avoid
	// Units Specifies the unit system to use when expressing distance as text.
	// Valid values are `UnitsMetric` and `UnitsImperial`. Optional.
	Units Units
	// DepartureTime is the desired time of departure. You can specify the time as
	// an integer in seconds since midnight, January 1, 1970 UTC. Alternatively,
	// you can specify a value of `"now"``. Optional.
	DepartureTime string
	// ArrivalTime specifies the desired time of arrival for transit requests,
	// in seconds since midnight, January 1, 1970 UTC. You cannot specify
	// both `DepartureTime` and `ArrivalTime`. Optional.
	ArrivalTime string
	// TrafficModel determines the type of model that will be used when determining
	// travel time when using depature times in the future. Options are
	// `TrafficModelBestGuess`, `TrafficModelOptimistic`` or `TrafficModelPessimistic`.
	// Optional. Default is `TrafficModelBestGuess``
	TrafficModel TrafficModel
	// TransitMode specifies one or more preferred modes of transit. This parameter
	// may only be specified for requests where the mode is `transit`. Valid values
	// are `TransitModeBus`, `TransitModeSubway`, `TransitModeTrain`, `TransitModeTram`,
	// and `TransitModeRail`. Optional.
	TransitMode []TransitMode
	// TransitRoutingPreference Specifies preferences for transit requests. Valid
	// values are `TransitRoutingPreferenceLessWalking` and
	// `TransitRoutingPreferenceFewerTransfers`. Optional.
	TransitRoutingPreference TransitRoutingPreference
}

// DistanceMatrixResponse represents a Distance Matrix API response.
type DistanceMatrixResponse struct {

	// OriginAddresses contains an array of addresses as returned by the API from
	// your original request.
	OriginAddresses []string `json:"origin_addresses"`
	// DestinationAddresses contains an array of addresses as returned by the API
	// from your original request.
	DestinationAddresses []string `json:"destination_addresses"`
	// Rows contains an array of elements.
	Rows []DistanceMatrixElementsRow `json:"rows"`
}

// DistanceMatrixElementsRow is a row of distance elements.
type DistanceMatrixElementsRow struct {
	Elements []*DistanceMatrixElement `json:"elements"`
}

// DistanceMatrixElement is the travel distance and time for a pair of origin
// and destination.
type DistanceMatrixElement struct {
	Status string `json:"status"`
	// Duration is the length of time it takes to travel this route.
	Duration time.Duration `json:"duration"`
	// DurationInTraffic is the length of time it takes to travel this route
	// considering traffic.
	DurationInTraffic time.Duration `json:"duration_in_traffic"`
	// Distance is the total distance of this route.
	Distance Distance `json:"distance"`
}
