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

package maps

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
)

var directionsAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/directions/json",
	acceptsClientID: true,
}

// Directions issues the Directions request and retrieves the Response
func (c *Client) Directions(ctx context.Context, r *DirectionsRequest) ([]Route, []GeocodedWaypoint, error) {
	if r.Origin == "" {
		return nil, nil, errors.New("maps: origin missing")
	}
	if r.Destination == "" {
		return nil, nil, errors.New("maps: destination missing")
	}
	if r.Mode != "" && TravelModeDriving != r.Mode && TravelModeWalking != r.Mode && TravelModeBicycling != r.Mode && TravelModeTransit != r.Mode {
		return nil, nil, fmt.Errorf("maps: unknown Mode: '%s'", r.Mode)
	}
	if r.DepartureTime != "" && r.ArrivalTime != "" {
		return nil, nil, errors.New("maps: DepartureTime and ArrivalTime both specified")
	}
	if len(r.TransitMode) != 0 && r.Mode != TravelModeTransit {
		return nil, nil, errors.New("maps: TransitMode specified while Mode != TravelModeTransit")
	}
	if r.TransitRoutingPreference != "" && r.Mode != TravelModeTransit {
		return nil, nil, errors.New("maps: mode of transit '" + string(r.Mode) + "' invalid for TransitRoutingPreference")
	}

	var response struct {
		Routes            []Route            `json:"routes"`
		GeocodedWaypoints []GeocodedWaypoint `json:"geocoded_waypoints"`
		commonResponse
	}

	if err := c.getJSON(ctx, directionsAPI, r, &response); err != nil {
		return nil, nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, nil, err
	}

	return response.Routes, response.GeocodedWaypoints, nil
}

func (r *DirectionsRequest) params() url.Values {
	q := make(url.Values)
	q.Set("origin", r.Origin)
	q.Set("destination", r.Destination)
	if r.Mode != "" {
		q.Set("mode", string(r.Mode))
	}
	if r.DepartureTime != "" {
		q.Set("departure_time", r.DepartureTime)
	}
	if r.ArrivalTime != "" {
		q.Set("arrival_time", r.ArrivalTime)
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
	if r.TransitRoutingPreference != "" {
		q.Set("transit_routing_preference", string(r.TransitRoutingPreference))
	}
	if r.TrafficModel != "" {
		q.Set("traffic_model", string(r.TrafficModel))
	}
	return q
}

// DirectionsRequest is the functional options struct for directions.Get
type DirectionsRequest struct {
	// Origin is the address or textual latitude/longitude value from which you wish to calculate directions. Required.
	Origin string
	// Destination is the address or textual latitude/longitude value from which you wish to calculate directions. Required.
	Destination string
	// Mode specifies the mode of transport to use when calculating directions. Optional.
	Mode Mode
	// DepartureTime specifies the desired time of departure. You can specify the time as an integer in seconds since midnight, January 1, 1970 UTC. Alternatively, you can specify a value of `"now"`. Optional.
	DepartureTime string
	// ArrivalTime specifies the desired time of arrival for transit directions, in seconds since midnight, January 1, 1970 UTC. Optional. You cannot specify both `DepartureTime` and `ArrivalTime`.
	ArrivalTime string
	// Waypoints specifies an array of points to add to a route. Optional.
	Waypoints []string
	// Alternatives specifies if Directions service may provide more than one route alternative in the response. Optional.
	Alternatives bool
	// Avoid indicates that the calculated route(s) should avoid the indicated features. Optional.
	Avoid []Avoid
	// Language specifies the language in which to return results. Optional.
	Language string
	// Units specifies the unit system to use when displaying results. Optional.
	Units Units
	// Region specifies the region code, specified as a ccTLD two-character value. Optional.
	Region string
	// TransitMode specifies one or more preferred modes of transit. This parameter may only be specified for transit directions. Optional.
	TransitMode []TransitMode
	// TransitRoutingPreference specifies preferences for transit routes. Optional.
	TransitRoutingPreference TransitRoutingPreference
	// TrafficModel specifies traffic prediction model when requesting future directions. Optional.
	TrafficModel TrafficModel
}

// GeocodedWaypoint represents the geocoded point for origin, supplied waypoints, or destination for a requested direction request.
type GeocodedWaypoint struct {
	// GeocoderStatus indicates the status code resulting from the geocoding operation. This field may contain the following values.
	GeocoderStatus string `json:"geocoder_status"`
	// PartialMatch indicates that the geocoder did not return an exact match for the original request, though it was able to match part of the requested address.
	PartialMatch bool `json:"partial_match"`
	// PlaceID is a unique identifier that can be used with other Google APIs.
	PlaceID string `json:"place_id"`
	// Types indicates the address type of the geocoding result used for calculating directions.
	Types []string `json:"types"`
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
	Duration time.Duration `json:"duration"`

	// DurationInTraffic indicates the total duration of this leg. This value is an estimate of the time in traffic based on current and historical traffic conditions.
	DurationInTraffic time.Duration `json:"duration_in_traffic"`

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
