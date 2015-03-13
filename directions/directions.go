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
	"log"
	"net/http"
	"time"

	"google.golang.org/maps"
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
	OverviewPolyline maps.Polyline `json:"overview_polyline"`

	// Bounds contains the viewport bounding box of the overview polyline.
	maps.Bounds `json:"bounds"`

	// Copyrights contains the copyrights text to be displayed for this route. You must handle
	// and display this information yourself.
	Copyrights string `json:"copyrights"`

	// Warnings contains an array of warnings to be displayed when showing these directions.
	// You must handle and display these warnings yourself.
	Warnings []string `json:"warnings"`
}

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
	StartLocation maps.LatLng `json:"start_location"`

	// EndLocation contains the latitude/longitude coordinates of the destination of this leg.
	EndLocation maps.LatLng `json:"end_location"`

	// StartAddress contains the human-readable address (typically a street address)
	// reflecting the start location of this leg.
	StartAddress string `json:"start_address"`

	// EndAddress contains the human-readable address (typically a street address)
	// reflecting the end location of this leg.
	EndAddress string `json:"end_address"`
}

// Step represents a single step of a leg.
type Step struct {
	// html_instructions contains formatted instructions for this step, presented as an HTML text string.
	HTMLInstructions string `json:"html_instructions"`

	// Distance contains the distance covered by this step until the next step.
	Distance `json:"distance"`

	// Duration contains the typical time required to perform the step, until the next step.
	time.Duration `json:"duration"`

	// StartLocation contains the location of the starting point of this step, as a single set of lat
	// and lng fields.
	StartLocation maps.LatLng `json:"start_location"`

	// end_location contains the location of the last point of this step, as a single set of lat and
	// lng fields.
	EndLocation maps.LatLng `json:"end_location"`

	// polyline contains a single points object that holds an encoded polyline representation of the
	// step. This polyline is an approximate (smoothed) path of the step.
	maps.Polyline `json:"polyline"`

	// Steps contains detailed directions for walking or driving steps in transit directions. Substeps
	// are only available when travel_mode is set to "transit". The inner steps array is of the same
	// type as steps.
	Steps []*Step `json:"steps"`

	// TransitDetails contains transit specific information. This field is only returned with travel
	// mode is set to "transit".
	TransitDetails `json:"transit_details"`

	// TravelMode indicates the travel mode of this step.
	TravelMode string `json:"travel_mode"`
}

type TransitDetails struct {
	// TODO(brettmorgan): fill this in
}

// Distance represents a distance covered in a step or leg.
type Distance struct {
	// Value indicates the distance in meters
	Value int64 `json:"value"`

	// Text contains a human-readable representation of the distance.
	Text string `json:"text"`
}

// Get retrieves directions between the specified origin and destination.
func Get(ctx context.Context, origin, destination string) (Response, error) {
	var response Response
	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/directions/json", nil)
	if err != nil {
		return response, err
	}
	q := req.URL.Query()
	q.Set("origin", origin)
	q.Set("destination", destination)
	q.Set("key", internal.APIKey(ctx))
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
