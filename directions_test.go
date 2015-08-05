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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

const apiKey = "AIzaNotReallyAnAPIKey"

type countingServer struct {
	s          *httptest.Server
	successful int
	failed     []string
}

// mockServerForQuery returns a mock server that only responds to a particular query string.
func mockServerForQuery(query string, code int, body string) *countingServer {
	server := &countingServer{}

	server.s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if query != "" && r.URL.RawQuery != query {
			server.failed = append(server.failed, r.URL.RawQuery)
			http.Error(w, "fail", 999)
			return
		}
		server.successful++

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))

	return server
}

// Create a mock HTTP Server that will return a response with HTTP code and body.
func mockServer(code int, body string) *httptest.Server {
	serv := mockServerForQuery("", code, body)
	return serv.s
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
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	resp, err := c.Directions(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	var steps []*Step
	steps = append(steps, &Step{
		HTMLInstructions: "Head <b>south</b> on <b>George St</b> toward <b>Barrack St</b>",
		Distance:         Distance{HumanReadable: "0.4 km", Meters: 366},
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
		Distance:      Distance{HumanReadable: "23.8 km", Meters: 23846},
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
		t.Errorf("expected %+v, was %+v", correctResponse, &resp[0])
	}
}

func TestDirectionsMissingOrigin(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Destination: "Parramatta",
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Missing Origin should return error")
	}
}

func TestDirectionsMissingDestination(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Origin: "Sydney",
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Missing Destination should return error")
	}
}

func TestDirectionsBadMode(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
		Mode:        "Not a Mode",
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Bad Mode should return error")
	}
}

func TestDirectionsDeclaringBothDepartureAndArrivalTime(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Origin:        "Sydney",
		Destination:   "Parramatta",
		DepartureTime: "Now",
		ArrivalTime:   "4pm",
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Declaring both DepartureTime and ArrivalTime should return error")
	}
}

func TestDirectionsTravelModeTransit(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	var transitModes []TransitMode
	transitModes = append(transitModes, TransitModeBus)
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
		TransitMode: transitModes,
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Declaring TransitMode without Mode=Transit should return error")
	}
}

func TestDirectionsTransitRoutingPreference(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Origin:                   "Sydney",
		Destination:              "Parramatta",
		TransitRoutingPreference: TransitRoutingPreferenceFewerTransfers,
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Declaring TransitRoutingPreference without Mode=TravelModeTransit should return error")
	}
}

func TestDirectionsWithCancelledContext(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Directions(ctx, r); err == nil {
		t.Errorf("Cancelled context should return non-nil err")
	}
}

func TestDirectionsFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	if _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Failing server should return error")
	}
}

func TestDirectionsRequestURL(t *testing.T) {
	expectedQuery := "alternatives=true&avoid=tolls%7Cferries&destination=Parramatta&key=AIzaNotReallyAnAPIKey&language=es&mode=transit&origin=Sydney&region=es&transit_mode=rail&transit_routing_preference=fewer_transfers&units=imperial&waypoints=Charlestown%2CMA%7Cvia%3ALexington"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.s.URL

	r := &DirectionsRequest{
		Origin:       "Sydney",
		Destination:  "Parramatta",
		Mode:         TravelModeTransit,
		TransitMode:  []TransitMode{TransitModeRail},
		Waypoints:    []string{"Charlestown,MA", "via:Lexington"},
		Alternatives: true,
		Avoid:        []Avoid{AvoidTolls, AvoidFerries},
		Language:     "es",
		Region:       "es",
		Units:        UnitsImperial,
		TransitRoutingPreference: TransitRoutingPreferenceFewerTransfers,
	}

	_, err := c.Directions(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}
