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
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestElevationDenver(t *testing.T) {

	// Elevation of Denver, the mile high city
	response := `{
   "results" : [
      {
         "elevation" : 1608.637939453125,
         "location" : {
            "lat" : 39.73915360,
            "lng" : -104.98470340
         },
         "resolution" : 4.771975994110107
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &ElevationRequest{
		Locations: []LatLng{
			{
				Lat: 39.73915360,
				Lng: -104.9847034,
			},
		},
	}

	resp, err := c.Elevation(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	correctResponse := ElevationResult{
		Location: &LatLng{
			Lat: 39.73915360,
			Lng: -104.98470340,
		},
		Elevation:  1608.637939453125,
		Resolution: 4.771975994110107,
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp[0])
	}
}

func TestElevationSampledPath(t *testing.T) {

	response := `{
  "results" : [
        {
           "elevation" : 4411.941894531250,
           "location" : {
              "lat" : 36.5785810,
              "lng" : -118.2919940
           },
           "resolution" : 19.08790397644043
        },
        {
           "elevation" : 1381.861694335938,
           "location" : {
              "lat" : 36.41150289067028,
              "lng" : -117.5602607523847
           },
           "resolution" : 19.08790397644043
        },
        {
           "elevation" : -84.61699676513672,
           "location" : {
              "lat" : 36.239980,
              "lng" : -116.831710
           },
           "resolution" : 19.08790397644043
        }
     ],
     "status" : "OK"
  }`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &ElevationRequest{
		Path: []LatLng{
			{Lat: 36.578581, Lng: -118.291994},
			{Lat: 36.23998, Lng: -116.83171},
		},
		Samples: 3,
	}

	resp, err := c.Elevation(context.Background(), r)

	if len(resp) != 3 {
		t.Errorf("Expected length of response is 3, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	correctResponse := ElevationResult{
		Location: &LatLng{
			Lat: 36.5785810,
			Lng: -118.2919940,
		},
		Elevation:  4411.941894531250,
		Resolution: 19.08790397644043,
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp[0])
	}
}

func TestElevationNoPathOrLocations(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &ElevationRequest{}

	if _, err := c.Elevation(context.Background(), r); err == nil {
		t.Errorf("Missing both Path and Locations should return error")
	}
}

func TestElevationPathWithNoSamples(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &ElevationRequest{
		Path: []LatLng{
			{Lat: 36.578581, Lng: -118.291994},
			{Lat: 36.23998, Lng: -116.83171},
		},
	}

	if _, err := c.Elevation(context.Background(), r); err == nil {
		t.Errorf("Missing both Path and Locations should return error")
	}
}

func TestElevationFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &ElevationRequest{
		Path: []LatLng{
			{Lat: 36.578581, Lng: -118.291994},
			{Lat: 36.23998, Lng: -116.83171},
		},
		Samples: 3,
	}

	if _, err := c.Elevation(context.Background(), r); err == nil {
		t.Errorf("Failing server should return error")
	}
}

func TestElevationCancelledContext(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &ElevationRequest{
		Path: []LatLng{
			{Lat: 36.578581, Lng: -118.291994},
			{Lat: 36.23998, Lng: -116.83171},
		},
		Samples: 3,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Elevation(ctx, r); err == nil {
		t.Errorf("Cancelled context should return non-nil err")
	}
}

func TestElevationRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&locations=enc%3A_ibE_seK_seK_seK&path=enc%3A_qo%5D_%7Brc%40_seK_seK&samples=10"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &ElevationRequest{
		Locations: []LatLng{{1, 2}, {3, 4}},
		Path:      []LatLng{{5, 6}, {7, 8}},
		Samples:   10,
	}

	_, err := c.Elevation(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}
