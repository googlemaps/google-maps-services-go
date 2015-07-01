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
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestDistanceMatrixSydPyrToPar(t *testing.T) {

	// Distance Matrix from Sydney And Pyrmont to Parramatta with most steps elided.
	response := `{
   "destination_addresses" : [ "Parramatta NSW, Australia" ],
   "origin_addresses" : [ "Sydney NSW, Australia", "Pyrmont NSW, Australia" ],
   "rows" : [
      {
         "elements" : [
            {
               "distance" : {
                  "text" : "23.8 km",
                  "value" : 23846
               },
               "duration" : {
                  "text" : "37 mins",
                  "value" : 2215
               },
               "status" : "OK"
            }
         ]
      },
      {
         "elements" : [
            {
               "distance" : {
                  "text" : "22.2 km",
                  "value" : 22242
               },
               "duration" : {
                  "text" : "34 mins",
                  "value" : 2058
               },
               "status" : "OK"
            }
         ]
      }
   ],
   "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), withBaseURL(server.URL))
	r := &DistanceMatrixRequest{
		Origins:      []string{"Sydney", "Pyrmont"},
		Destinations: []string{"Parramatta"},
	}

	resp, err := c.GetDistanceMatrix(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	correctResponse := DistanceMatrixResponse{
		OriginAddresses:      []string{"Sydney NSW, Australia", "Pyrmont NSW, Australia"},
		DestinationAddresses: []string{"Parramatta NSW, Australia"},
		Rows: []DistanceMatrixElementsRow{
			DistanceMatrixElementsRow{
				Elements: []*DistanceMatrixElement{
					&DistanceMatrixElement{
						Status:   "OK",
						Duration: 2215000000000,
						Distance: Distance{Text: "23.8 km", Value: 23846},
					},
				},
			},
			DistanceMatrixElementsRow{
				Elements: []*DistanceMatrixElement{
					&DistanceMatrixElement{
						Status:   "OK",
						Duration: 2058000000000,
						Distance: Distance{Text: "22.2 km", Value: 22242},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(resp, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp)
	}
}

func TestDistanceMatrixMissingOrigins(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DistanceMatrixRequest{
		Origins:      []string{},
		Destinations: []string{"Parramatta"},
	}

	if _, err := c.GetDistanceMatrix(context.Background(), r); err == nil {
		t.Errorf("Missing Origins should return error")
	}
}

func TestDistanceMatrixMissingDestinations(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DistanceMatrixRequest{
		Origins:      []string{"Sydney", "Pyrmont"},
		Destinations: []string{},
	}

	if _, err := c.GetDistanceMatrix(context.Background(), r); err == nil {
		t.Errorf("Missing Destinations should return error")
	}
}

func TestDistanceMatrixDepartureAndArrivalTime(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &DistanceMatrixRequest{
		Origins:       []string{"Sydney", "Pyrmont"},
		Destinations:  []string{},
		DepartureTime: "now",
		ArrivalTime:   "4pm",
	}

	if _, err := c.GetDistanceMatrix(context.Background(), r); err == nil {
		t.Errorf("Having both Departure time and Arrival time should return error")
	}
}

func TestDistanceMatrixFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), withBaseURL(server.URL))
	r := &DistanceMatrixRequest{
		Origins:      []string{"Sydney", "Pyrmont"},
		Destinations: []string{},
	}

	if _, err := c.GetDistanceMatrix(context.Background(), r); err == nil {
		t.Errorf("Failing server should return error")
	}
}
