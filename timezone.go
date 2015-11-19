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
	"errors"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"
)

var timezoneAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/timezone/json",
	acceptsClientID: true,
}

// Timezone makes a Timezone API request
func (c *Client) Timezone(ctx context.Context, r *TimezoneRequest) (*TimezoneResult, error) {
	if r.Location == nil {
		return nil, errors.New("maps: Location missing")
	}

	var response struct {
		TimezoneResult
		commonResponse
	}

	if err := c.getJSON(ctx, timezoneAPI, r, &response); err != nil {
		return nil, err
	}

	if err := response.StatusError(); err != nil {
		return nil, err
	}

	return &response.TimezoneResult, nil
}

func (r *TimezoneRequest) params() url.Values {
	q := make(url.Values)
	q.Set("location", r.Location.String())
	q.Set("timestamp", strconv.FormatInt(r.Timestamp.Unix(), 10))
	if r.Language != "" {
		q.Set("language", r.Language)
	}
	return q
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

// TimezoneResult is a single timezone result.
type TimezoneResult struct {
	// DstOffset is the offset for daylight-savings time in seconds.
	DstOffset int `json:"dstOffset"`
	// RawOffset is the offset from UTC for the given location.
	RawOffset int `json:"rawOffset"`
	// TimeZoneID is a string containing the "tz" ID of the time zone.
	TimeZoneID string `json:"timeZoneId"`
	// TimeZoneName is a string containing the long form name of the time zone.
	TimeZoneName string `json:"timeZoneName"`
}
