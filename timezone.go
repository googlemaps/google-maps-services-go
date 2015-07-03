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
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

// GetTimezone makes a Timezone API request
func (c *Client) GetTimezone(ctx context.Context, r *TimezoneRequest) (TimezoneResult, error) {
	baseURL := "https://maps.googleapis.com/"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}

	req, err := http.NewRequest("GET", baseURL+"/maps/api/timezone/json", nil)
	if err != nil {
		return TimezoneResult{}, err
	}
	q := req.URL.Query()
	q.Set("key", c.apiKey)

	if r.Location == nil {
		return TimezoneResult{}, errors.New("timezone: You must specify Location")
	}

	q.Set("location", r.Location.String())
	q.Set("timestamp", strconv.FormatInt(r.Timestamp.Unix(), 10))
	if r.Language != "" {
		q.Set("language", r.Language)
	}

	req.URL.RawQuery = q.Encode()
	chResult := make(chan timezoneResponse)

	go func() {
		resp, err := c.httpDo(req)
		if err != nil {
			chResult <- timezoneResponse{err: err}
			return
		}

		var response timezoneResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			chResult <- timezoneResponse{err: err}
			return
		}
		if response.Status != "OK" {
			chResult <- timezoneResponse{err: fmt.Errorf("timezone: %s - %s", response.Status, response.ErrorMessage)}
			return
		}

		chResult <- response
	}()

	select {
	case resp := <-chResult:
		return TimezoneResult{
			DstOffset:    resp.DstOffset,
			RawOffset:    resp.RawOffset,
			TimeZoneID:   resp.TimeZoneID,
			TimeZoneName: resp.TimeZoneName,
		}, resp.err
	case <-ctx.Done():
		return TimezoneResult{}, ctx.Err()
	}

}

// TimezoneRequest is the request structure for Timezone API.
type TimezoneRequest struct {
	// Location represents the location to look up.
	Location *LatLng
	// Timestamp specifies the desired time. Time Zone API uses the timestamp to determine whether or not Daylight Savings should be applied.
	Timestamp time.Time
	// Language in which to return results.
	Language string
}

type timezoneResponse struct {
	// DstOffset is the offset for daylight-savings time in seconds.
	DstOffset int `json:"dstOffset"`
	// RawOffset is the offset from UTC for the given location.
	RawOffset int `json:"rawOffset"`
	// TimeZoneID is a string containing the "tz" ID of the time zone.
	TimeZoneID string `json:"timeZoneId"`
	// TimeZoneName is a string containing the long form name of the time zone.
	TimeZoneName string `json:"timeZoneName"`

	// Status indicating if this request was successful
	Status string `json:"status"`
	// ErrorMessage is the explanatory field added when Status is an error.
	ErrorMessage string `json:"error_message"`

	err error
}

// TimezoneResult is a single geocoded address
type TimezoneResult struct {
	// DstOffset is the offset for daylight-savings time in seconds.
	DstOffset int
	// RawOffset is the offset from UTC for the given location.
	RawOffset int
	// TimeZoneID is a string containing the "tz" ID of the time zone.
	TimeZoneID string
	// TimeZoneName is a string containing the long form name of the time zone.
	TimeZoneName string
}
