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

package maps

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

// Timezone makes a Timezone API request
func (c *Client) Timezone(ctx context.Context, r *TimezoneRequest) (*TimezoneResult, error) {
	if r.Location == nil {
		return nil, errors.New("timezone: You must specify Location")
	}

	chResult := make(chan timezoneResultWithError)

	go func() {
		result, err := c.doGetTimezone(r)
		chResult <- timezoneResultWithError{result, err}
	}()

	select {
	case result := <-chResult:
		return result.timezone, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) doGetTimezone(r *TimezoneRequest) (*TimezoneResult, error) {
	baseURL := "https://maps.googleapis.com/"
	if c.baseURL != "" {
		baseURL = c.baseURL
	}

	req, err := http.NewRequest("GET", baseURL+"/maps/api/timezone/json", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()

	q.Set("location", r.Location.String())
	q.Set("timestamp", strconv.FormatInt(r.Timestamp.Unix(), 10))
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	query, err := c.generateAuthQuery(req.URL.Path, q, true)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query

	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}

	response := &timezoneResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		err := fmt.Errorf("timezone: %s - %s", response.Status, response.ErrorMessage)
		return nil, err
	}

	var result = &TimezoneResult{
		DstOffset:    response.DstOffset,
		RawOffset:    response.RawOffset,
		TimeZoneID:   response.TimeZoneID,
		TimeZoneName: response.TimeZoneName,
	}

	return result, nil
}

type timezoneResultWithError struct {
	timezone *TimezoneResult
	err      error
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
