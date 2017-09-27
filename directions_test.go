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
	"time"

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

func TestDirectionsTransit(t *testing.T) {
	// Route from Google Sydney to Glebe Pt Rd. Steps and some polylines have
	// been removed.
	response := `{
   "geocoded_waypoints" : [
      {
         "geocoder_status" : "OK",
         "partial_match" : true,
         "place_id" : "ChIJ8UadyjeuEmsRDt5QbiDg720",
         "types" : [ "premise" ]
      },
      {
         "geocoder_status" : "OK",
         "place_id" : "ChIJl7hmVNOvEmsRcW-sYgALB78",
         "types" : [ "route" ]
      }
   ],
   "routes" : [
      {
         "bounds" : {
            "northeast" : {
               "lat" : -33.8668939,
               "lng" : 151.1952284
            },
            "southwest" : {
               "lat" : -33.8785317,
               "lng" : 151.1856793
            }
         },
         "copyrights" : "Map data ©2016 Google",
         "legs" : [
            {
               "arrival_time" : {
                  "text" : "4:09pm",
                  "time_zone" : "Australia/Sydney",
                  "value" : 1455512950
               },
               "departure_time" : {
                  "text" : "4:00pm",
                  "time_zone" : "Australia/Sydney",
                  "value" : 1455512400
               },
               "distance" : {
                  "text" : "2.2 km",
                  "value" : 2241
               },
               "duration" : {
                  "text" : "9 mins",
                  "value" : 550
               },
               "end_address" : "Glebe Point Rd, Glebe NSW 2037, Australia",
               "end_location" : {
                  "lat" : -33.8785317,
                  "lng" : 151.1859855
               },
               "start_address" : "Workplace 6, 48 Pirrama Rd, Pyrmont NSW 2009, Australia",
               "start_location" : {
                  "lat" : -33.8675125,
                  "lng" : 151.1950229
               },
               "steps" : [
               ],
               "via_waypoint" : []
            }
         ],
         "overview_polyline" : {
            "points" : ""
         },
         "summary" : "",
         "warnings" : [
            "Walking directions are in beta.    Use caution – This route may be missing sidewalks or pedestrian paths."
         ],
         "waypoint_order" : []
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &DirectionsRequest{
		Origin:      "Google Sydney",
		Destination: "Glebe Pt Rd, Glebe",
		Mode:        TravelModeTransit,
	}

	resp, _, err := c.Directions(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	tzSydney, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		t.Errorf("coudln't load expected timezone Australia/Sydney: %v", err)
	}
	arrivalTime := time.Date(2016, 2, 15, 16, 9, 10, 0, tzSydney)
	departureTime := time.Date(2016, 2, 15, 16, 0, 0, 0, tzSydney)

	var legs []*Leg
	legs = append(legs, &Leg{
		Steps:         make([]*Step, 0),
		Distance:      Distance{HumanReadable: "2.2 km", Meters: 2241},
		Duration:      time.Duration(550) * time.Second,
		ArrivalTime:   arrivalTime,
		DepartureTime: departureTime,
		StartLocation: LatLng{Lat: -33.8675125, Lng: 151.1950229},
		EndLocation:   LatLng{Lat: -33.8785317, Lng: 151.1859855},
		StartAddress:  "Workplace 6, 48 Pirrama Rd, Pyrmont NSW 2009, Australia",
		EndAddress:    "Glebe Point Rd, Glebe NSW 2037, Australia",
	})

	correctResponse := &Route{
		OverviewPolyline: Polyline{},
		Legs:             legs,
		Bounds: LatLngBounds{
			NorthEast: LatLng{Lat: -33.8668939, Lng: 151.1952284},
			SouthWest: LatLng{Lat: -33.8785317, Lng: 151.1856793},
		},
		Copyrights: "Map data ©2016 Google",
		Warnings: []string{
			"Walking directions are in beta.    Use caution – This route may be missing sidewalks or pedestrian paths.",
		},
		WaypointOrder: make([]int, 0),
	}

	if actualResponse := &resp[0]; !reflect.DeepEqual(actualResponse, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, actualResponse)
	}
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
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	resp, _, err := c.Directions(context.Background(), r)

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

	if _, _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Missing Origin should return error")
	}
}

func TestDirectionsMissingDestination(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DirectionsRequest{
		Origin: "Sydney",
	}

	if _, _, err := c.Directions(context.Background(), r); err == nil {
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

	if _, _, err := c.Directions(context.Background(), r); err == nil {
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

	if _, _, err := c.Directions(context.Background(), r); err == nil {
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

	if _, _, err := c.Directions(context.Background(), r); err == nil {
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

	if _, _, err := c.Directions(context.Background(), r); err == nil {
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
	if _, _, err := c.Directions(ctx, r); err == nil {
		t.Errorf("Cancelled context should return non-nil err")
	}
}

func TestDirectionsFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	if _, _, err := c.Directions(context.Background(), r); err == nil {
		t.Errorf("Failing server should return error")
	}
}

func TestDirectionsRequestURL(t *testing.T) {
	expectedQuery := "alternatives=true&avoid=tolls%7Cferries&destination=Parramatta&key=AIzaNotReallyAnAPIKey&language=es&mode=transit&optimize=true&origin=Sydney&region=es&transit_mode=rail&transit_routing_preference=fewer_transfers&units=imperial&waypoints=Charlestown%2CMA%7Cvia%3ALexington"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &DirectionsRequest{
		Origin:       "Sydney",
		Destination:  "Parramatta",
		Mode:         TravelModeTransit,
		TransitMode:  []TransitMode{TransitModeRail},
		Waypoints:    []string{"Charlestown,MA", "via:Lexington"},
		Alternatives: true,
		Optimize:     true,
		Avoid:        []Avoid{AvoidTolls, AvoidFerries},
		Language:     "es",
		Region:       "es",
		Units:        UnitsImperial,
		TransitRoutingPreference: TransitRoutingPreferenceFewerTransfers,
	}

	if _, _, err := c.Directions(context.Background(), r); err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestTrafficModel(t *testing.T) {
	expectedQuery := "departure_time=now&destination=Parramatta+Town+Hall&key=AIzaNotReallyAnAPIKey&mode=driving&origin=Sydney+Town+Hall&traffic_model=pessimistic"
	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &DirectionsRequest{
		Origin:        "Sydney Town Hall",
		Destination:   "Parramatta Town Hall",
		Mode:          TravelModeDriving,
		DepartureTime: "now",
		TrafficModel:  TrafficModelPessimistic,
	}

	if _, _, err := c.Directions(context.Background(), r); err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestFare(t *testing.T) {
	// Directions response, sans steps.
	response := `{
   "geocoded_waypoints" : [
      {
         "geocoder_status" : "OK",
         "place_id" : "ChIJt4ABNBBfrIkRXFPZnsnIicw",
         "types" : [ "street_address" ]
      },
      {
         "geocoder_status" : "OK",
         "place_id" : "EjExMDAxLTExOTkgSGlsbHNib3JvdWdoIFN0LCBSYWxlaWdoLCBOQyAyNzYwMywgVVNB",
         "types" : [ "street_address" ]
      }
   ],
   "routes" : [
      {
         "bounds" : {
            "northeast" : {
               "lat" : 35.782453,
               "lng" : -78.625349
            },
            "southwest" : {
               "lat" : 35.7766535,
               "lng" : -78.6548689
            }
         },
         "copyrights" : "Map data ©2016 Google",
         "fare" : {
            "currency" : "USD",
            "text" : "$1.25",
            "value" : 1.25
         },
         "legs" : [
            {
               "arrival_time" : {
                  "text" : "8:23am",
                  "time_zone" : "America/New_York",
                  "value" : 1476620634
               },
               "departure_time" : {
                  "text" : "8:01am",
                  "time_zone" : "America/New_York",
                  "value" : 1476619311
               },
               "distance" : {
                  "text" : "2.1 mi",
                  "value" : 3408
               },
               "duration" : {
                  "text" : "22 mins",
                  "value" : 1323
               },
               "end_address" : "1001-1199 Hillsborough St, Raleigh, NC 27603, USA",
               "end_location" : {
                  "lat" : 35.782453,
                  "lng" : -78.6548689
               },
               "start_address" : "822 E Hargett St, Raleigh, NC 27601, USA",
               "start_location" : {
                  "lat" : 35.7777912,
                  "lng" : -78.625349
               },
               "traffic_speed_entry" : [],
               "via_waypoint" : []
            }
         ],
         "summary" : "",
         "warnings" : [
            "Walking directions are in beta.    Use caution – This route may be missing sidewalks or pedestrian paths."
         ],
         "waypoint_order" : []
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &DirectionsRequest{
		Origin:      "35.7777111, -78.625354",
		Destination: "35.7826242, -78.6547025",
		Mode:        TravelModeTransit,
	}

	resp, _, err := c.Directions(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	fare := &Fare{
		Currency: "USD",
		Text:     "$1.25",
		Value:    1.25,
	}

	if !reflect.DeepEqual(resp[0].Fare, fare) {
		t.Errorf("expected %+v, was %+v", fare, resp[0].Fare)
	}
}

func TestNoFare(t *testing.T) {
	// Directions response, sans steps.
	response := `{
   "geocoded_waypoints" : [
      {
         "geocoder_status" : "OK",
         "place_id" : "ChIJt4ABNBBfrIkRXFPZnsnIicw",
         "types" : [ "street_address" ]
      },
      {
         "geocoder_status" : "OK",
         "place_id" : "EjExMDAxLTExOTkgSGlsbHNib3JvdWdoIFN0LCBSYWxlaWdoLCBOQyAyNzYwMywgVVNB",
         "types" : [ "street_address" ]
      }
   ],
   "routes" : [
      {
         "bounds" : {
            "northeast" : {
               "lat" : 35.782453,
               "lng" : -78.625349
            },
            "southwest" : {
               "lat" : 35.7766535,
               "lng" : -78.6548689
            }
         },
         "copyrights" : "Map data ©2016 Google",
         "legs" : [
            {
               "arrival_time" : {
                  "text" : "8:23am",
                  "time_zone" : "America/New_York",
                  "value" : 1476620634
               },
               "departure_time" : {
                  "text" : "8:01am",
                  "time_zone" : "America/New_York",
                  "value" : 1476619311
               },
               "distance" : {
                  "text" : "2.1 mi",
                  "value" : 3408
               },
               "duration" : {
                  "text" : "22 mins",
                  "value" : 1323
               },
               "end_address" : "1001-1199 Hillsborough St, Raleigh, NC 27603, USA",
               "end_location" : {
                  "lat" : 35.782453,
                  "lng" : -78.6548689
               },
               "start_address" : "822 E Hargett St, Raleigh, NC 27601, USA",
               "start_location" : {
                  "lat" : 35.7777912,
                  "lng" : -78.625349
               },
               "traffic_speed_entry" : [],
               "via_waypoint" : []
            }
         ],
         "summary" : "",
         "warnings" : [
            "Walking directions are in beta.    Use caution – This route may be missing sidewalks or pedestrian paths."
         ],
         "waypoint_order" : []
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &DirectionsRequest{
		Origin:      "35.7777111, -78.625354",
		Destination: "35.7826242, -78.6547025",
		Mode:        TravelModeTransit,
	}

	resp, _, err := c.Directions(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	if resp[0].Fare != nil {
		t.Errorf("expected %+v, was %+v", nil, resp[0].Fare)
	}
}
