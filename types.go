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

// Mode is for specifying travel mode.
type Mode string

// Avoid is for specifying routes that avoid certain features.
type Avoid string

// Units specifies which units system to return human readable results in.
type Units string

// TransitMode is for specifying a transit mode for a request
type TransitMode string

// TransitRoutingPreference biases which routes are returned
type TransitRoutingPreference string

const (
	// TravelModeDriving is for specifying driving as travel mode
	TravelModeDriving = Mode("driving")
	// TravelModeWalking is for specifying walking as travel mode
	TravelModeWalking = Mode("walking")
	// TravelModeBicycling is for specifying bicycling as travel mode
	TravelModeBicycling = Mode("bicycling")
	// TravelModeTransit is for specifying transit as travel mode
	TravelModeTransit = Mode("transit")

	// AvoidTolls is for specifying routes that avoid tolls
	AvoidTolls = Avoid("tolls")
	// AvoidHighways is for specifying routes that avoid highways
	AvoidHighways = Avoid("highways")
	// AvoidFerries is for specifying routes that avoid ferries
	AvoidFerries = Avoid("ferries")

	// UnitsMetric specifies usage of the metric units system
	UnitsMetric = Units("metric")
	// UnitsImperial specifies usage of the Imperial (English) units system
	UnitsImperial = Units("imperial")

	// TransitModeBus is for specifying a transit mode of bus
	TransitModeBus = TransitMode("bus")
	// TransitModeSubway is for specifying a transit mode of subway
	TransitModeSubway = TransitMode("subway")
	// TransitModeTrain is for specifying a transit mode of train
	TransitModeTrain = TransitMode("train")
	// TransitModeTram is for specifying a transit mode of tram
	TransitModeTram = TransitMode("tram")
	// TransitModeRail is for specifying a transit mode of rail
	TransitModeRail = TransitMode("rail")

	// TransitRoutingPreferenceLessWalking indicates that the calculated route should prefer limited amounts of walking
	TransitRoutingPreferenceLessWalking = TransitRoutingPreference("less_walking")
	// TransitRoutingPreferenceFewerTransfers indicates that the calculated route should prefer a limited number of transfers
	TransitRoutingPreferenceFewerTransfers = TransitRoutingPreference("fewer_transfers")
)

// Distance is the API representation for a distance between two points.
type Distance struct {
	// HumanReadable is the human friendly distance. This is rounded and in an appropriate unit for the
	// request. The units can be overriden with a request parameter.
	HumanReadable string `json:"text"`
	// Meters is the numeric distance, always in meters. This is intended to be used only in
	// algorithmic situations, e.g. sorting results by some user specified metric.
	Meters int `json:"value"`
}
