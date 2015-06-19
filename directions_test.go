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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const apiKey = "AIzaNotReallyAnAPIKey"

// Create a mock HTTP Server that will return a response with HTTP code and body.
func mockServer(code int, body string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))

	return server
}

func TestDirectionsSydneyToParramatta(t *testing.T) {

	// Route from Sydney to Parramatta with most steps elided.
	response := `{
   "routes" : [
      {
         "bounds" : {
            "northeast" : {
               "lat" : -33.8150985,
               "lng" : 151.2070825
            },
            "southwest" : {
               "lat" : -33.8770049,
               "lng" : 151.0031658
            }
         },
         "copyrights" : "Map data ©2015 Google",
         "legs" : [
            {
               "distance" : {
                  "text" : "23.8 km",
                  "value" : 23846
               },
               "duration" : {
                  "text" : "37 mins",
                  "value" : 2214
               },
               "end_address" : "Parramatta NSW, Australia",
               "end_location" : {
                  "lat" : -33.8150985,
                  "lng" : 151.0031658
               },
               "start_address" : "Sydney NSW, Australia",
               "start_location" : {
                  "lat" : -33.8674944,
                  "lng" : 151.2070825
               },
               "steps" : [
                  {
                     "distance" : {
                        "text" : "0.4 km",
                        "value" : 366
                     },
                     "duration" : {
                        "text" : "2 mins",
                        "value" : 103
                     },
                     "end_location" : {
                        "lat" : -33.8707786,
                        "lng" : 151.206934
                     },
                     "html_instructions" : "Head \u003cb\u003esouth\u003c/b\u003e on \u003cb\u003eGeorge St\u003c/b\u003e toward \u003cb\u003eBarrack St\u003c/b\u003e",
                     "polyline" : {
                        "points" : "xvumEgs{y[V@|AH|@DdABbC@@?^@N?zD@\\?F@"
                     },
                     "start_location" : {
                        "lat" : -33.8674944,
                        "lng" : 151.2070825
                     },
                     "travel_mode" : "DRIVING"
                  }
               ],
               "via_waypoint" : []
            }
         ],
         "summary" : "A4 and M4"
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	resp, err := r.Get(ctx)

	if len(resp) != 1 {
		t.Errorf("Expected length of response was 1, actually %g", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	var steps []*Step
	steps = append(steps, &Step{
		HTMLInstructions: "Head <b>south</b> on <b>George St</b> toward <b>Barrack St</b>",
		Distance:         Distance{Text: "0.4 km", Value: 366},
		Duration:         103000000000,
		StartLocation:    LatLng{Lat: -33.8674944, Lng: 151.2070825},
		EndLocation:      LatLng{Lat: -33.8707786, Lng: 151.206934},
		Polyline:         Polyline{Points: "xvumEgs{y[V@|AH|@DdABbC@@?^@N?zD@\\?F@"},
		Steps:            nil,
		TransitDetails:   (*TransitDetails)(nil),
		TravelMode:       "DRIVING",
	})

	var legs []*Leg
	legs = append(legs, &Leg{
		Steps:         steps,
		Distance:      Distance{Text: "23.8 km", Value: 23846},
		Duration:      2214000000000,
		StartLocation: LatLng{Lat: -33.8674944, Lng: 151.2070825},
		EndLocation:   LatLng{Lat: -33.8150985, Lng: 151.0031658},
		StartAddress:  "Sydney NSW, Australia",
		EndAddress:    "Parramatta NSW, Australia",
	})

	correctResponse := &Route{
		Summary:          "A4 and M4",
		Legs:             legs,
		OverviewPolyline: Polyline{},
		Bounds: LatLngBounds{
			NorthEast: LatLng{Lat: -33.8150985, Lng: 151.2070825},
			SouthWest: LatLng{Lat: -33.8770049, Lng: 151.0031658},
		},
		Copyrights: "Map data ©2015 Google",
	}

	if !reflect.DeepEqual(&resp[0], correctResponse) {
		t.Errorf("Actual response != expected")
	}
}

func TestDirectionsMissingOrigin(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &DirectionsRequest{
		Destination: "Parramatta",
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Missing Origin should return error")
	}
}

func TestDirectionsMissingDestination(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &DirectionsRequest{
		Origin: "Sydney",
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Missing Destination should return error")
	}
}

func TestDirectionsBadMode(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
		Mode:        "Not a Mode",
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Bad Mode should return error")
	}
}

func TestDirectionsDeclaringBothDepartureAndArrivalTime(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &DirectionsRequest{
		Origin:        "Sydney",
		Destination:   "Parramatta",
		DepartureTime: "Now",
		ArrivalTime:   "4pm",
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Declaring both DepartureTime and ArrivalTime should return error")
	}
}

func TestDirectionsTravelModeTransit(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	var transitModes []transitMode
	transitModes = append(transitModes, TransitModeBus)
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
		TransitMode: transitModes,
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Declaring TransitMode without Mode=Transit should return error")
	}
}

func TestDirectionsTransitRoutingPreference(t *testing.T) {
	client := &http.Client{}
	ctx := NewContext(apiKey, client)
	r := &DirectionsRequest{
		Origin:                   "Sydney",
		Destination:              "Parramatta",
		TransitRoutingPreference: TransitRoutingPreferenceFewerTransfers,
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Declaring TransitRoutingPreference without Mode=TravelModeTransit should return error")
	}
}

func TestDirectionsFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	client := &http.Client{}
	ctx := newContextWithBaseURL(apiKey, client, server.URL)
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	if _, err := r.Get(ctx); err == nil {
		t.Errorf("Failing server should return error")
	}
}
