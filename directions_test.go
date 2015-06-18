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
	"net/url"
	"reflect"
	"testing"
)

// Test that two values are equal, log if not equal.
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// Test two values are not equal, log if they are equal.
func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// Create a mock HTTP Server that will return a response with HTTP code and body.
func mockServer(code int, body string) (*httptest.Server, *http.Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: tr}

	return server, httpClient
}

func TestSydneyToParramatta(t *testing.T) {

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
         "overview_polyline" : {
            "points" : ""
         },
         "summary" : "A4 and M4",
         "warnings" : [],
         "waypoint_order" : []
      }
   ],
   "status" : "OK"
}`

	apiKey := "AIzaNotReallyAnAPIKey"
	server, client := mockServer(200, response)
	defer server.Close()

	ctx := NewContextWithBaseURL(apiKey, client, server.URL)
	r := &DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Parramatta",
	}

	resp, err := r.Get(ctx)

	expect(t, len(resp), 1)
	expect(t, err, nil)

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

	// Attempting to directly compare &resp[0] and correctResponse failed, yet this works. Help...
	expect(t, resp[0].Summary, correctResponse.Summary)
	expect(t, resp[0].Legs[0].Steps[0].HTMLInstructions, correctResponse.Legs[0].Steps[0].HTMLInstructions)
	expect(t, resp[0].Legs[0].Steps[0].Distance, correctResponse.Legs[0].Steps[0].Distance)
	expect(t, resp[0].Legs[0].Steps[0].Duration, correctResponse.Legs[0].Steps[0].Duration)
	expect(t, resp[0].Legs[0].Steps[0].StartLocation, correctResponse.Legs[0].Steps[0].StartLocation)
	expect(t, resp[0].Legs[0].Steps[0].EndLocation, correctResponse.Legs[0].Steps[0].EndLocation)
	expect(t, resp[0].Legs[0].Steps[0].Polyline, correctResponse.Legs[0].Steps[0].Polyline)
	expect(t, resp[0].Legs[0].Steps[0].TransitDetails, correctResponse.Legs[0].Steps[0].TransitDetails)
	expect(t, resp[0].Legs[0].Steps[0].TravelMode, correctResponse.Legs[0].Steps[0].TravelMode)
	expect(t, resp[0].Legs[0].Distance, correctResponse.Legs[0].Distance)
	expect(t, resp[0].Legs[0].Duration, correctResponse.Legs[0].Duration)
	expect(t, resp[0].Legs[0].StartLocation, correctResponse.Legs[0].StartLocation)
	expect(t, resp[0].Legs[0].EndLocation, correctResponse.Legs[0].EndLocation)
	expect(t, resp[0].Legs[0].StartAddress, correctResponse.Legs[0].StartAddress)
	expect(t, resp[0].Legs[0].EndAddress, correctResponse.Legs[0].EndAddress)
	expect(t, resp[0].OverviewPolyline, correctResponse.OverviewPolyline)
	expect(t, resp[0].Bounds, correctResponse.Bounds)
	expect(t, resp[0].Copyrights, correctResponse.Copyrights)

}
