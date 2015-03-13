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

type mode string
type avoid string
type units string
type transitMode string
type transitRoutingPreference string

const (
	// TravelModeDriving is for specifying driving as travel mode
	TravelModeDriving = mode("driving")
	// TravelModeWalking is for specifying walking as travel mode
	TravelModeWalking = mode("walking")
	// TravelModeBicycling is for specifying bicycling as travel mode
	TravelModeBicycling = mode("bicycling")
	// TravelModeTransit is for specifying transit as travel mode
	TravelModeTransit = mode("transit")

	// AvoidTolls is for specifying routes that avoid tolls
	AvoidTolls = avoid("tolls")
	// AvoidHighways is for specifying routes that avoid highways
	AvoidHighways = avoid("highways")
	// AvoidFerries is for specifying routes that avoid ferries
	AvoidFerries = avoid("ferries")

	// UnitsMetric specifies usage of the metric units system
	UnitsMetric = units("metric")
	// UnitsImperial specifies usage of the Imperial (English) units system
	UnitsImperial = units("imperial")

	// TransitModeBus is for specifying a transit mode of bus
	TransitModeBus = transitMode("bus")
	// TransitModeSubway is for specifying a transit mode of subway
	TransitModeSubway = transitMode("subway")
	// TransitModeTrain is for specifying a transit mode of train
	TransitModeTrain = transitMode("train")
	// TransitModeTram is for specifying a transit mode of tram
	TransitModeTram = transitMode("tram")
	// TransitModeRail is for specifying a transit mode of rail
	TransitModeRail = transitMode("rail")

	// TransitRoutingPreferenceLessWalking indicates that the calculated route should prefer limited amounts of walking
	TransitRoutingPreferenceLessWalking = transitRoutingPreference("less_walking")
	// TransitRoutingPreferenceFewerTransfers indicates that the calculated route should prefer a limited number of transfers
	TransitRoutingPreferenceFewerTransfers = transitRoutingPreference("fewer_transfers")
)

// Distance is the API representation for a distance between two points.
type Distance struct {
	// Text is the distance in a human displayable form. The style of display can be changed by setting `units`.
	Text string `json:"text"`
	// Value is the distance in meters.
	Value int `json:"value"`
}
