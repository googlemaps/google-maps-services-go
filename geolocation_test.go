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

// More information about Google Geolocation API is available on
// https://developers.google.com/maps/documentation/geolocation

package maps

import (
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

func TestGeolocation(t *testing.T) {

	// Elevation of Denver, the mile high city
	response := `{
		"location" : {
			"lat" : 39.73915360,
			"lng" : -104.98470340
		},
		"accuracy" : 4.771975994110107
	}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
	r := &GeolocationRequest{}

	resp, err := c.Geolocate(context.Background(), r)
	if err != nil {
		t.Errorf("r.Get returned non nil error, was %+v", err)
	}

	correctResponse := GeolocationResult{
		Location: LatLng{
			Lat: 39.73915360,
			Lng: -104.98470340,
		},
		Accuracy: 4.771975994110107,
	}

	if !reflect.DeepEqual(*resp, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp)
	}
}

func TestGeolocationError(t *testing.T) {

	// Elevation of Denver, the mile high city
	response := `{
		"error": {
			"errors": [
				{
					"domain": "global",
					"reason": "invalid",
					"message": "Invalid value for UnsignedInteger: ",
					"locationType": "other",
					"location": "homeMobileCountryCode"
				}
			],
			"code": 400,
			"message": "Invalid value for UnsignedInteger: "
		}
	}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
	r := &GeolocationRequest{}

	_, err := c.Geolocate(context.Background(), r)
	if err == nil {
		t.Errorf("r.Get returned nil error")
	}

	correctResponse := "Invalid value for UnsignedInteger: "

	if !reflect.DeepEqual(err.Error(), correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, err)
	}
}
