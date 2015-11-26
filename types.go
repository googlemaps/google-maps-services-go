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

import "time"

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

// Travel mode preferences.
const (
	TravelModeDriving   = Mode("driving")
	TravelModeWalking   = Mode("walking")
	TravelModeBicycling = Mode("bicycling")
	TravelModeTransit   = Mode("transit")
)

// Features to avoid.
const (
	AvoidTolls    = Avoid("tolls")
	AvoidHighways = Avoid("highways")
	AvoidFerries  = Avoid("ferries")
)

// Units to use on human readable distances.
const (
	UnitsMetric   = Units("metric")
	UnitsImperial = Units("imperial")
)

// Transit mode of directions or distance matrix request.
const (
	TransitModeBus    = TransitMode("bus")
	TransitModeSubway = TransitMode("subway")
	TransitModeTrain  = TransitMode("train")
	TransitModeTram   = TransitMode("tram")
	TransitModeRail   = TransitMode("rail")
)

// Transit Routing preferences for transit mode requests
const (
	TransitRoutingPreferenceLessWalking    = TransitRoutingPreference("less_walking")
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

// TrafficModel specifies traffic prediction model when requesting future directions.
type TrafficModel string

// Traffic prediction model when requesting future directions.
const (
	TrafficModelBestGuess   = TrafficModel("best_guess")
	TrafficModelOptimistic  = TrafficModel("optimistic")
	TrafficModelPessimistic = TrafficModel("pessimistic")
)

// PriceLevel is the Price Levels for Places API
type PriceLevel string

// Price Levels for the Places API
const (
	PriceLevelFree          = PriceLevel("0")
	PriceLevelInexpensive   = PriceLevel("1")
	PriceLevelModerate      = PriceLevel("2")
	PriceLevelExpensive     = PriceLevel("3")
	PriceLevelVeryExpensive = PriceLevel("4")
)

// OpeningHours describes the opening hours for a Place Details result.
type OpeningHours struct {
	// OpenNow is a boolean value indicating if the place is open at the current time. Please note, this field will be null if it isn't present in the response.
	OpenNow *bool `json:"open_now"`
	// Periods is an array of opening periods covering seven days, starting from Sunday, in chronological order.
	Periods []OpeningHoursPeriod `json:"periods"`
	// weekdayText is an array of seven strings representing the formatted opening hours for each day of the week, for example "Monday: 8:30 am – 5:30 pm".
	WeekdayText []string `json:"weekday_text"`
	// PermanentlyClosed indicates that the place has permanently shut down. Please note, this field will be null if it isn't present in the response.
	PermanentlyClosed *bool `json:"permanently_closed"`
}

// OpeningHoursPeriod is a single OpeningHours day describing when the place opens and closes.
type OpeningHoursPeriod struct {
	// Open is when the place opens.
	Open OpeningHoursOpenClose `json:"open"`
	// Close is when the place closes.
	Close OpeningHoursOpenClose `json:"close"`
}

// OpeningHoursOpenClose describes when the place is open.
type OpeningHoursOpenClose struct {
	// Day is a number from 0–6, corresponding to the days of the week, starting on Sunday. For example, 2 means Tuesday.
	Day time.Weekday `json:"day"`
	// Time contains a time of day in 24-hour hhmm format. Values are in the range 0000–2359. The time will be reported in the place’s time zone.
	Time string `json:"time"`
}

// Photo describes a photo available with a Search Result.
type Photo struct {
	// PhotoReference is used to identify the photo when you perform a Photo request.
	PhotoReference string `json:"photo_reference"`
	// Height is the maximum height of the image.
	Height int `json:"height"`
	// Width is the maximum width of the image.
	Width int `json:"width"`
	// htmlAttributions contains any required attributions.
	HTMLAttributions []string `json:"html_attributions"`
}
