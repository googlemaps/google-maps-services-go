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

// Package maps contains TODO(brettmorgan)
//
// More information about Google Distance Matrix API is available on
// https://developers.google.com/maps/documentation/distancematrix/
package maps // import "google.golang.org/maps"
import "time"

const (
	// ModeDriving is for specifying driving as travel mode
	ModeDriving = "driving"
	// ModeWalking is for specifying walking as travel mode
	ModeWalking = "walking"
	// ModeBicycling is for specifying bicycling as travel mode
	ModeBicycling = "bicycling"
	// ModeTransit is for specifying transit as travel mode
	ModeTransit = "transit"

	// AvoidTolls is for specifying routes that avoid tolls
	AvoidTolls = "tolls"
	// AvoidHighways is for specifying routes that avoid highways
	AvoidHighways = "highways"
	// AvoidFerries is for specifying routes that avoid ferries
	AvoidFerries = "ferries"

	// UnitsMetric specifies usage of the metric units system
	UnitsMetric = "metric"
	// UnitsImperial specifies usage of the Imperial (English) units system
	UnitsImperial = "imperial"

	// TransitModeBus is for specifying a transit mode of bus
	TransitModeBus = "bus"
	// TransitModeSubway is for specifying a transit mode of subway
	TransitModeSubway = "subway"
	// TransitModeTrain is for specifying a transit mode of train
	TransitModeTrain = "train"
	// TransitModeTram is for specifying a transit mode of tram
	TransitModeTram = "tram"
	// TransitModeRail is for specifying a transit mode of rail
	TransitModeRail = "rail"

	// TransitRoutingPreferenceLessWalking indicates that the calculated route should prefer limited amounts of walking
	TransitRoutingPreferenceLessWalking = "less_walking"
	// TransitRoutingPreferenceFewerTransfers indicates that the calculated route should prefer a limited number of transfers
	TransitRoutingPreferenceFewerTransfers = "fewer_transfers"
)

// DistanceMatrixRequest is the functional options struct for DistanceMatrixGet
type DistanceMatrixRequest struct {
	// Origins is a list of addresses and/or textual latitude/longitude values from which to calculate distance and time. Required.
	Origins []string
	// Destinations is a list of addresses and/or textual latitude/longitude values to which to calculate distance and time. Required.
	Destinations []string
	// Mode specifies the mode of transport to use when calculating distance. Valid values are `ModeDriving`, `ModeWalking`, `ModeBicycling`
	// and `ModeTransit`. Optional.
	Mode string
	// Language in which to return results. Optional.
	Language string
	// Avoid introduces restrictions to the route. Valid values are `AvoidTolls`, `AvoidHighways` and `AvoidFerries`. Optional.
	Avoid string
	// Units Specifies the unit system to use when expressing distance as text. Valid values are `UnitsMetric` and `UnitsImperial`. Optional.
	Units string
	// DepartureTime is the desired time of departure. You can specify the time as an integer in seconds since midnight, January 1, 1970 UTC.
	// Alternatively, you can specify a value of `"now"``. Optional.
	DepartureTime string
	// ArrivalTime specifies the desired time of arrival for transit requests, in seconds since midnight, January 1, 1970 UTC. You cannot
	// specify both `DepartureTime` and `ArrivalTime`. Optional.
	ArrivalTime string
	// TransitMode specifies one or more preferred modes of transit. This parameter may only be specified for requests where the mode is
	// `transit`. Valid values are `TransitModeBus`, `TransitModeSubway`, `TransitModeTrain`, `TransitModeTram`, and `TransitModeRail`.
	// Optional.
	TransitMode string
	// TransitRoutingPreference Specifies preferences for transit requests. Valid values are `TransitRoutingPreferenceLessWalking` and
	// `TransitRoutingPreferenceFewerTransfers`. Optional.
	TransitRoutingPreference string
}

// DistanceMatrixResponse represents a Distance Matrix API response.
type DistanceMatrixResponse struct {

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

// Distance is the API representation for a distance between two points.
type Distance struct {
	// Text is the distance in a human displayable form. The style of display can be changed by setting `units`.
	Text string `json:"text"`
	// Value is the distance in meters.
	Value int `json:"value"`
}
