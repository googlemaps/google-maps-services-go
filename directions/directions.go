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

// Package directions contains a Google Directions API client.
//
// More information about Google Directions API is available on
// https://developers.google.com/maps/documentation/directions/
package directions // import "google.golang.org/maps/directions"

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

// Response represents a Directions API response.
type Response struct {
	// Routes lists the found routes between origin and destination.
	Routes []Route

	// Status contains the status of the request, and may contain
	// debugging information to help you track down why the Directions
	// service failed.
	// See https://developers.google.com/maps/documentation/directions/#StatusCodes
	Status string
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
	Bounds `json:"bounds"`

	// Copyrights contains the copyrights text to be displayed for this route. You must handle
	// and display this information yourself.
	Copyrights string `json:"copyrights"`

	// Warnings contains an array of warnings to be displayed when showing these directions.
	// You must handle and display these warnings yourself.
	Warnings []string `json:"warnings"`
}

// Bounds represents a bounded area on a map.
type Bounds struct {
	// The north east corner of the bounded area.
	NorthEast LatLng `json:"northeast"`

	// The south west corner of the bounded area.
	SouthWest LatLng `json:"southwest"`
}

// LatLng represents a location.
type LatLng struct {
	// Lat is the latitude of this location.
	Lat float64 `json:"lat"`

	// Lng is the longitude of this location.
	Lng float64 `json:"lng"`
}

// Polyline represents a list of lat,lng points, encoded as a string.
// See: https://developers.google.com/maps/documentation/utilities/polylinealgorithm
type Polyline struct {
	Points string `json:"points"`
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

// Distance represents a distance covered in a step or leg.
type Distance struct {
	// Value indicates the distance in meters
	Value int64 `json:"value"`

	// Text contains a human-readable representation of the distance.
	Text string `json:"text"`
}

// DirectionsRequest is the functional options struct for directions.Get
type DirectionsRequest struct {
	// Origin is the address or textual latitude/longitude value from which you wish to calculate directions. Required.
	Origin string
	// Destination is the address or textual latitude/longitude value from which you wish to calculate directions. Required.
	Destination string
	// Mode specifies the mode of transport to use when calculating directions. Optional.
	Mode string
	// DepartureTime specifies the desired time of departure. You can specify the time as an integer in seconds since midnight, January 1, 1970 UTC. Alternatively, you can specify a value of `"now"`. Optional.
	DepartureTime string
	// ArrivalTime specifies the desired time of arrival for transit directions, in seconds since midnight, January 1, 1970 UTC. Optional. You cannot specify both `DepartureTime` and `ArrivalTime`.
	ArrivalTime string
	// Waypoints specifies an array of points to add to a route. Optional.
	Waypoints []string
	// Alternatives specifies if Directions service may provide more than one route alternative in the response. Optional.
	Alternatives bool
	// Avoid indicates that the calculated route(s) should avoid the indicated features. Optional.
	Avoid []string
	// Language specifies the language in which to return results. Optional.
	Language string
	// Units specifies the unit system to use when displaying results. Optional.
	Units string
	// Region specifies the region code, specified as a ccTLD two-character value. Optional.
	Region string
	// TransitMode specifies one or more preferred modes of transit. This parameter may only be specified for transit directions. Optional.
	TransitMode []string
	// TransitRoutingPreference specifies preferences for transit routes. Optional.
	TransitRoutingPreference string
}

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

// Get issues the Directions request and retrieves the Response
func (dirReq *DirectionsRequest) Get(ctx context.Context) (Response, error) {
	var response Response

	if dirReq.Origin == "" {
		return response, errors.New("directions: Origin required")
	}
	if dirReq.Destination == "" {
		return response, errors.New("directions: Destination required")
	}
	if dirReq.Mode != "" && ModeDriving != dirReq.Mode && ModeWalking != dirReq.Mode && ModeBicycling != dirReq.Mode && ModeTransit != dirReq.Mode {
		return response, fmt.Errorf("directions: unknown Mode: '%s'", dirReq.Mode)
	}
	for _, avoid := range dirReq.Avoid {
		if avoid != AvoidTolls && avoid != AvoidHighways && avoid != AvoidFerries {
			return response, fmt.Errorf("directions: Unknown Avoid restriction '%s'", avoid)
		}
	}
	if dirReq.Units != "" && dirReq.Units != UnitsMetric && dirReq.Units != UnitsImperial {
		return response, fmt.Errorf("directions: Unknown Units '%s'", dirReq.Units)
	}
	for _, transitMode := range dirReq.TransitMode {
		if transitMode != TransitModeBus && transitMode != TransitModeSubway && transitMode != TransitModeTrain && transitMode != TransitModeTram && transitMode != TransitModeRail {
			return response, fmt.Errorf("directions: Unknown TransitMode '%s'", dirReq.TransitMode)
		}
	}
	if dirReq.TransitRoutingPreference != "" && dirReq.TransitRoutingPreference != TransitRoutingPreferenceLessWalking && dirReq.TransitRoutingPreference != TransitRoutingPreferenceFewerTransfers {
		return response, fmt.Errorf("directions: Unknown TransitRoutingPreference '%s'", dirReq.TransitRoutingPreference)
	}
	if dirReq.DepartureTime != "" && dirReq.ArrivalTime != "" {
		return response, errors.New("directions: must not specify both DepartureTime and ArrivalTime")
	}

	if dirReq.DepartureTime != "" && dirReq.ArrivalTime != "" {
		return response, errors.New("directions: must not specify both DepartureTime and ArrivalTime")
	}
	if len(dirReq.TransitMode) != 0 && dirReq.Mode != ModeTransit {
		return response, errors.New("directions: must specify mode of transit when specifying transitMode")
	}
	if dirReq.TransitRoutingPreference != "" && dirReq.Mode != ModeTransit {
		return response, errors.New("directions: must specify mode of transit when specifying transitRoutingPreference")
	}

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/directions/json", nil)
	if err != nil {
		return response, err
	}
	q := req.URL.Query()
	q.Set("origin", dirReq.Origin)
	q.Set("destination", dirReq.Destination)
	q.Set("key", internal.APIKey(ctx))
	if dirReq.Mode != "" {
		q.Set("mode", dirReq.Mode)
	}
	if len(dirReq.Waypoints) != 0 {
		q.Set("waypoints", strings.Join(dirReq.Waypoints, "|"))
	}
	if dirReq.Alternatives {
		q.Set("alternatives", "true")
	}
	if len(dirReq.Avoid) > 0 {
		q.Set("avoid", strings.Join(dirReq.Avoid, "|"))
	}
	if dirReq.Language != "" {
		q.Set("language", dirReq.Language)
	}
	if dirReq.Units != "" {
		q.Set("units", dirReq.Units)
	}
	if dirReq.Region != "" {
		q.Set("region", dirReq.Region)
	}
	if len(dirReq.TransitMode) != 0 {
		q.Set("transit_mode", strings.Join(dirReq.TransitMode, "|"))
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
	// httpDo waits for the closure we provided to return, so it's safe to
	// read response here.
	return response, err
}

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func rawService(ctx context.Context) *http.Client {
	return internal.Service(ctx, "directions", func(hc *http.Client) interface{} {
		// TODO(brettmorgan): Introduce a rate limiting wrapper for hc here.
		return hc
	}).(*http.Client)
}
