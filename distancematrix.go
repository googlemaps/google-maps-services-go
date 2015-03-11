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

// More information about Google Distance Matrix API is available on
// https://developers.google.com/maps/documentation/distancematrix/
package maps // import "google.golang.org/maps"

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/maps/internal"
)

// Get makes a Distance Matrix API request
func (dmReq *DistanceMatrixRequest) Get(ctx context.Context) (DistanceMatrixResponse, error) {
	var response DistanceMatrixResponse
	if len(dmReq.Origins) == 0 {
		return response, errors.New("distancematrix: Origins must contain at least one start address")
	}
	if len(dmReq.Destinations) == 0 {
		return response, errors.New("distancematrix: Destinations must contain at least one end address")
	}
	if dmReq.Mode != "" && ModeDriving != dmReq.Mode && ModeWalking != dmReq.Mode && ModeBicycling != dmReq.Mode && ModeTransit != dmReq.Mode {
		return response, fmt.Errorf("distancematrix: unknown Mode: '%s'", dmReq.Mode)
	}
	if dmReq.Avoid != "" && dmReq.Avoid != AvoidTolls && dmReq.Avoid != AvoidHighways && dmReq.Avoid != AvoidFerries {
		return response, fmt.Errorf("distancematrix: Unknown Avoid restriction '%s'", dmReq.Avoid)
	}
	if dmReq.Units != "" && dmReq.Units != UnitsMetric && dmReq.Units != UnitsImperial {
		return response, fmt.Errorf("distancematrix: Unknown Units '%s'", dmReq.Units)
	}
	if dmReq.TransitMode != "" && dmReq.TransitMode != TransitModeBus && dmReq.TransitMode != TransitModeSubway && dmReq.TransitMode != TransitModeTrain && dmReq.TransitMode != TransitModeTram && dmReq.TransitMode != TransitModeRail {
		return response, fmt.Errorf("distancematrix: Unknown TransitMode '%s'", dmReq.TransitMode)
	}
	if dmReq.TransitRoutingPreference != "" && dmReq.TransitRoutingPreference != TransitRoutingPreferenceLessWalking && dmReq.TransitRoutingPreference != TransitRoutingPreferenceFewerTransfers {
		return response, fmt.Errorf("distancematrix: Unknown TransitRoutingPreference '%s'", dmReq.TransitRoutingPreference)
	}
	if dmReq.DepartureTime != "" && dmReq.ArrivalTime != "" {
		return response, errors.New("distancematrix: must not specify both DepartureTime and ArrivalTime")
	}

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/distancematrix/json", nil)
	if err != nil {
		return response, err
	}
	q := req.URL.Query()
	q.Set("origins", strings.Join(dmReq.Origins, "|"))
	q.Set("destinations", strings.Join(dmReq.Destinations, "|"))
	q.Set("key", internal.APIKey(ctx))
	if dmReq.Mode != "" {
		q.Set("mode", dmReq.Mode)
	}
	if dmReq.Language != "" {
		q.Set("language", dmReq.Language)
	}
	if dmReq.Avoid != "" {
		q.Set("avoid", dmReq.Avoid)
	}
	if dmReq.Units != "" {
		q.Set("units", dmReq.Units)
	}
	if dmReq.DepartureTime != "" {
		q.Set("departure_time", dmReq.DepartureTime)
	}
	if dmReq.ArrivalTime != "" {
		q.Set("arrival_time", dmReq.ArrivalTime)
	}
	if dmReq.TransitMode != "" {
		q.Set("transit_mode", dmReq.TransitMode)
	}
	if dmReq.TransitRoutingPreference != "" {
		q.Set("transit_routing_preference", dmReq.TransitRoutingPreference)
	}

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
