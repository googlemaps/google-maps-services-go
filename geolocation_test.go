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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"golang.org/x/net/context"
)

func TestGeolocation(t *testing.T) {

	// Denver, the mile high city
	response := `{
		"location" : {
			"lat" : 39.73915360,
			"lng" : -104.98470340
		},
		"accuracy" : 4.771975994110107
	}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
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

	// An error response
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
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
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

func TestCellTowerAndWiFiRequest(t *testing.T) {
	// Denver, the mile high city
	response := `{
		"location" : {
			"lat" : 39.73915360,
			"lng" : -104.98470340
		},
		"accuracy" : 4.771975994110107
	}`

	server := &countingServer{}

	failResponse := func(reason string, w http.ResponseWriter, r *http.Request) {
		server.failed = append(server.failed, r.URL.RawQuery)
		s := fmt.Sprintf(`{"status":"fail", "message": "%s"}`, reason)
		http.Error(w, s, 999)
	}

	server.s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			failResponse("failed to read body", w, r)
			return
		}
		body := string(b)
		expected := `{"homeMobileCountryCode":310,` +
			`"homeMobileNetworkCode":410,` +
			`"radioType":"gsm",` +
			`"carrier":"Vodafone",` +
			`"considerIp":true,` +
			`"cellTowers":[{"cellId":42,` +
			`"locationAreaCode":415,` +
			`"mobileCountryCode":310,` +
			`"mobileNetworkCode":410,` +
			`"signalStrength":-60,` +
			`"timingAdvance":15}],` +
			`"wifiAccessPoints":[{"macAddress":"00:25:9c:cf:1c:ac",` +
			`"signalStrength":-43,` +
			`"channel":11}]}`
		if body != expected {
			pretty.Errorf("Body is incorrect: %v", body)
			failResponse("failed to parse body", w, r)
			return
		}

		server.successful++

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, response)
	}))

	defer server.s.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))
	r := &GeolocationRequest{
		HomeMobileCountryCode: 310,
		HomeMobileNetworkCode: 410,
		RadioType:             RadioTypeGSM,
		Carrier:               "Vodafone",
		ConsiderIP:            true,
		CellTowers: []CellTower{{
			CellID:            42,
			LocationAreaCode:  415,
			MobileCountryCode: 310,
			MobileNetworkCode: 410,
			Age:               0,
			SignalStrength:    -60,
			TimingAdvance:     15,
		}},
		WiFiAccessPoints: []WiFiAccessPoint{{
			MACAddress:         "00:25:9c:cf:1c:ac",
			SignalStrength:     -43,
			Age:                0,
			Channel:            11,
			SignalToNoiseRatio: 0,
		}},
	}

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
