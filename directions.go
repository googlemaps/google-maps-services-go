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

// More information about Google Directions API is available on
// https://developers.google.com/maps/documentation/directions/

package maps // import "google.golang.org/maps"

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/maps/internal"

	"golang.org/x/net/context"
)

// Get issues the Directions request and retrieves the Response
func (r *DirectionsRequest) Get(ctx context.Context) ([]Route, error) {
	var response DirectionsResponse

	if r.Origin == "" {
		return nil, errors.New("directions: Origin required")
	}
	if r.Destination == "" {
		return nil, errors.New("directions: Destination required")
	}
	if r.Mode != "" && TravelModeDriving != r.Mode && TravelModeWalking != r.Mode && TravelModeBicycling != r.Mode && TravelModeTransit != r.Mode {
		return nil, fmt.Errorf("directions: unknown Mode: '%s'", r.Mode)
	}
	for _, avoid := range r.Avoid {
		if avoid != AvoidTolls && avoid != AvoidHighways && avoid != AvoidFerries {
			return nil, fmt.Errorf("directions: Unknown Avoid restriction '%s'", avoid)
		}
	}
	if r.Units != "" && r.Units != UnitsMetric && r.Units != UnitsImperial {
		return nil, fmt.Errorf("directions: Unknown Units '%s'", r.Units)
	}
	for _, transitMode := range r.TransitMode {
		if transitMode != TransitModeBus && transitMode != TransitModeSubway && transitMode != TransitModeTrain && transitMode != TransitModeTram && transitMode != TransitModeRail {
			return nil, fmt.Errorf("directions: Unknown TransitMode '%s'", r.TransitMode)
		}
	}
	if r.TransitRoutingPreference != "" && r.TransitRoutingPreference != TransitRoutingPreferenceLessWalking && r.TransitRoutingPreference != TransitRoutingPreferenceFewerTransfers {
		return nil, fmt.Errorf("directions: Unknown TransitRoutingPreference '%s'", r.TransitRoutingPreference)
	}
	if r.DepartureTime != "" && r.ArrivalTime != "" {
		return nil, errors.New("directions: must not specify both DepartureTime and ArrivalTime")
	}

	if r.DepartureTime != "" && r.ArrivalTime != "" {
		return nil, errors.New("directions: must not specify both DepartureTime and ArrivalTime")
	}
	if len(r.TransitMode) != 0 && r.Mode != TravelModeTransit {
		return nil, errors.New("directions: must specify mode of transit when specifying transitMode")
	}
	if r.TransitRoutingPreference != "" && r.Mode != TravelModeTransit {
		return nil, errors.New("directions: must specify mode of transit when specifying transitRoutingPreference")
	}

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/directions/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("origin", r.Origin)
	q.Set("destination", r.Destination)
	q.Set("key", internal.APIKey(ctx))
	if r.Mode != "" {
		q.Set("mode", string(r.Mode))
	}
	if len(r.Waypoints) != 0 {
		q.Set("waypoints", strings.Join(r.Waypoints, "|"))
	}
	if r.Alternatives {
		q.Set("alternatives", "true")
	}
	if len(r.Avoid) > 0 {
		var avoid []string
		for _, a := range r.Avoid {
			avoid = append(avoid, string(a))
		}
		q.Set("avoid", strings.Join(avoid, "|"))
	}
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	if r.Units != "" {
		q.Set("units", string(r.Units))
	}
	if r.Region != "" {
		q.Set("region", r.Region)
	}
	if len(r.TransitMode) != 0 {
		var transitMode []string
		for _, t := range r.TransitMode {
			transitMode = append(transitMode, string(t))
		}
		q.Set("transit_mode", strings.Join(transitMode, "|"))
	}
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
		return nil, fmt.Errorf("directions: %s - %s", response.Status, response.ErrorMessage)
	}

	return response.Routes, nil
}

// DirectionsRequest is the functional options struct for directions.Get
type DirectionsRequest struct {
	// Origin is the address or textual latitude/longitude value from which you wish to calculate directions. Required.
	Origin string
	// Destination is the address or textual latitude/longitude value from which you wish to calculate directions. Required.
	Destination string
	// Mode specifies the mode of transport to use when calculating directions. Optional.
	Mode mode
	// DepartureTime specifies the desired time of departure. You can specify the time as an integer in seconds since midnight, January 1, 1970 UTC. Alternatively, you can specify a value of `"now"`. Optional.
	DepartureTime string
	// ArrivalTime specifies the desired time of arrival for transit directions, in seconds since midnight, January 1, 1970 UTC. Optional. You cannot specify both `DepartureTime` and `ArrivalTime`.
	ArrivalTime string
	// Waypoints specifies an array of points to add to a route. Optional.
	Waypoints []string
	// Alternatives specifies if Directions service may provide more than one route alternative in the response. Optional.
	Alternatives bool
	// Avoid indicates that the calculated route(s) should avoid the indicated features. Optional.
	Avoid []avoid
	// Language specifies the language in which to return results. Optional.
	Language string
	// Units specifies the unit system to use when displaying results. Optional.
	Units units
	// Region specifies the region code, specified as a ccTLD two-character value. Optional.
	Region string
	// TransitMode specifies one or more preferred modes of transit. This parameter may only be specified for transit directions. Optional.
	TransitMode []transitMode
	// TransitRoutingPreference specifies preferences for transit routes. Optional.
	TransitRoutingPreference transitRoutingPreference
}

// DirectionsResponse represents a Directions API response.
type DirectionsResponse struct {
	// Routes lists the found routes between origin and destination.
	Routes []Route `json:"routes"`

	// Status contains the status of the request, and may contain
	// debugging information to help you track down why the Directions
	// service failed.
	// See https://developers.google.com/maps/documentation/directions/#StatusCodes
	Status string `json:"status"`

	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`
}

// Route represents a single route between an origin and a destination.
type Route struct {
	// Summary contains a short textual description for the route, suitable for
	// naming and disambiguating the route from alternatives.
	Summary string `json:"summary"`

	// Legs contains information about a leg of the route, between two locations within the
	// given route. A separate leg will be present for each waypoint or destination specified.
	// (A route with no waypoints will contain exactly one leg within the legs array.)
	Legs []*Leg `json:"legs"`

	// WaypointOrder contains an array indicating the order of any waypoints in the calculated route.
	WaypointOrder []int `json:"waypoint_order"`

	// OverviewPolyline contains an approximate (smoothed) path of the resulting directions.
	OverviewPolyline Polyline `json:"overview_polyline"`

	// Bounds contains the viewport bounding box of the overview polyline.
	Bounds LatLngBounds `json:"bounds"`

	// Copyrights contains the copyrights text to be displayed for this route. You must handle
	// and display this information yourself.
	Copyrights string `json:"copyrights"`

	// Warnings contains an array of warnings to be displayed when showing these directions.
	// You must handle and display these warnings yourself.
	Warnings []string `json:"warnings"`
}

// Leg represents a single leg of a route.
type Leg struct {
	// Steps contains an array of steps denoting information about each separate step of the
	// leg of the journey.
	Steps []*Step `json:"steps"`

	// Distance indicates the total distance covered by this leg.
	Distance `json:"distance"`

	// Duration indicates total time required for this leg.
	time.Duration `json:"duration"`

	// ArrivalTime contains the estimated time of arrival for this leg. This property is only
	// returned for transit directions.
	ArrivalTime time.Time `json:"arrival_time"`

	// DepartureTime contains the estimated time of departure for this leg. This property is
	// only returned for transit directions.
	DepartureTime time.Time `json:"departure_time"`

	// StartLocation contains the latitude/longitude coordinates of the origin of this leg.
	StartLocation LatLng `json:"start_location"`

	// EndLocation contains the latitude/longitude coordinates of the destination of this leg.
	EndLocation LatLng `json:"end_location"`

	// StartAddress contains the human-readable address (typically a street address)
	// reflecting the start location of this leg.
	StartAddress string `json:"start_address"`

	// EndAddress contains the human-readable address (typically a street address)
	// reflecting the end location of this leg.
	EndAddress string `json:"end_address"`
}

// Step represents a single step of a leg.
type Step struct {
	// HTMLInstructions contains formatted instructions for this step, presented as an HTML text string.
	HTMLInstructions string `json:"html_instructions"`

	// Distance contains the distance covered by this step until the next step.
	Distance `json:"distance"`

	// Duration contains the typical time required to perform the step, until the next step.
	time.Duration `json:"duration"`

	// StartLocation contains the location of the starting point of this step, as a single set of lat
	// and lng fields.
	StartLocation LatLng `json:"start_location"`

	// EndLocation contains the location of the last point of this step, as a single set of lat and
	// lng fields.
	EndLocation LatLng `json:"end_location"`

	// Polyline contains a single points object that holds an encoded polyline representation of the
	// step. This polyline is an approximate (smoothed) path of the step.
	Polyline `json:"polyline"`

	// Steps contains detailed directions for walking or driving steps in transit directions. Substeps
	// are only available when travel_mode is set to "transit". The inner steps array is of the same
	// type as steps.
	Steps []*Step `json:"steps"`

	// TransitDetails contains transit specific information. This field is only returned with travel
	// mode is set to "transit".
	TransitDetails *TransitDetails `json:"transit_details"`

	// TravelMode indicates the travel mode of this step.
	TravelMode string `json:"travel_mode"`
}

// TransitDetails contains additional information about the transit stop, transit line and transit agency.
type TransitDetails struct {
	// ArrivalStop contains information about the stop/station for this part of the trip.
	ArrivalStop TransitStop `json:"arrival_stop"`
	// DepartureStop contains information about the stop/station for this part of the trip.
	DepartureStop TransitStop `json:"departure_stop"`
	// ArrivalTime contains the arrival time for this leg of the journey
	ArrivalTime time.Time `json:"arrival_time"`
	// DepartureTime contains the departure time for this leg of the journey
	DepartureTime time.Time `json:"departure_time"`
	// Headsign specifies the direction in which to travel on this line, as it is marked on the vehicle or at the departure stop.
	Headsign string `json:"headsign"`
	// Headway specifies the expected number of seconds between departures from the same stop at this time
	Headway time.Duration `json:"headway"`
	// NumStops contains the number of stops in this step, counting the arrival stop, but not the departure stop
	NumStops uint `json:"num_stops"`
	// Line contains information about the transit line used in this step
	Line TransitLine `json:"line"`
}

// TransitStop contains information about the stop/station for this part of the trip.
type TransitStop struct {
	// Location of the transit station/stop.
	Location LatLng `json:"location"`
	// Name of the transit station/stop. eg. "Union Square".
	Name string `json:"name"`
}

// TransitLine contains information about the transit line used in this step
type TransitLine struct {
	// Name contains the full name of this transit line. eg. "7 Avenue Express".
	Name string `json:"name"`
	// ShortName contains the short name of this transit line.
	ShortName string `json:"short_name"`
	// Color contains the color commonly used in signage for this transit line.
	Color string `json:"color"`
	// Agencies contains information about the operator of the line
	Agencies []*TransitAgency `json:"agencies"`
	// URL contains the URL for this transit line as provided by the transit agency
	URL *url.URL `json:"url"`
	// Icon contains the URL for the icon associated with this line
	Icon *url.URL `json:"icon"`
	// TextColor contains the color of text commonly used for signage of this line
	TextColor string `json:"text_color"`
	// Vehicle contains the type of vehicle used on this line
	Vehicle TransitLineVehicle `json:"vehicle"`
}

// TransitAgency contains information about the operator of the line
type TransitAgency struct {
	// Name contains the name of the transit agency
	Name string `json:"name"`
	// URL contains the URL for the transit agency
	URL *url.URL `json:"url"`
	// Phone contains the phone number of the transit agency
	Phone string `json:"phone"`
}

// TransitLineVehicle contains the type of vehicle used on this line
type TransitLineVehicle struct {
	// Name contains the name of the vehicle on this line
	Name string `json:"name"`
	// Type contains the type of vehicle that runs on this line
	Type string `json:"type"`
	// Icon contains the URL for an icon associated with this vehicle type
	Icon *url.URL `json:"icon"`
}
