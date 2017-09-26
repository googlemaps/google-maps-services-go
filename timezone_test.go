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
	"time"

	"github.com/kr/pretty"
	"golang.org/x/net/context"
)

func TestTimezoneNevada(t *testing.T) {

	response := `{
   "dstOffset" : 0,
   "rawOffset" : -28800,
   "status" : "OK",
   "timeZoneId" : "America/Los_Angeles",
   "timeZoneName" : "Pacific Standard Time"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &TimezoneRequest{
		Location: &LatLng{
			Lat: 39.6034810,
			Lng: -119.6822510,
		},
		Timestamp: time.Unix(1331161200, 0),
	}

	resp, err := c.Timezone(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := &TimezoneResult{
		DstOffset:    0,
		RawOffset:    -28800,
		TimeZoneID:   "America/Los_Angeles",
		TimeZoneName: "Pacific Standard Time",
	}

	if !reflect.DeepEqual(resp, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp)
		pretty.Println(resp)
		pretty.Println(correctResponse)
	}
}

func TestTimezoneLocationMissing(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &TimezoneRequest{
		Timestamp: time.Unix(1331161200, 0),
	}

	if _, err := c.Timezone(context.Background(), r); err == nil {
		t.Errorf("Missing Location should return error")
	}
}

func TestTimezoneWithCancelledContext(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &TimezoneRequest{
		Location: &LatLng{
			Lat: 39.6034810,
			Lng: -119.6822510,
		},
		Timestamp: time.Unix(1331161200, 0),
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Timezone(ctx, r); err == nil {
		t.Errorf("Cancelled context should return error")
	}
}

func TestTimezoneFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &TimezoneRequest{
		Location: &LatLng{Lat: 36.578581, Lng: -118.291994},
	}

	if _, err := c.Timezone(context.Background(), r); err == nil {
		t.Errorf("Failing server should return error")
	}
}

func TestTimezoneRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&language=es&location=1%2C2&timestamp=-62135596800"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &TimezoneRequest{
		Location:  &LatLng{1, 2},
		Timestamp: time.Time{},
		Language:  "es",
	}

	_, err := c.Timezone(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}
