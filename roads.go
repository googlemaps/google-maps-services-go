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
	"context"
	"errors"
	"net/url"
	"strings"
)

var snapToRoadsAPI = &apiConfig{
	host:             "https://roads.googleapis.com",
	path:             "/v1/snapToRoads",
	acceptsClientID:  false,
	acceptsSignature: false,
}

var nearestRoadsAPI = &apiConfig{
	host:             "https://roads.googleapis.com",
	path:             "/v1/nearestRoads",
	acceptsClientID:  false,
	acceptsSignature: false,
}

var speedLimitsAPI = &apiConfig{
	host:             "https://roads.googleapis.com",
	path:             "/v1/speedLimits",
	acceptsClientID:  false,
	acceptsSignature: false,
}

// SnapToRoad makes a Snap to Road API request
func (c *Client) SnapToRoad(ctx context.Context, r *SnapToRoadRequest) (*SnapToRoadResponse, error) {

	if len(r.Path) == 0 {
		return nil, errors.New("maps: Path empty")
	}

	response := &SnapToRoadResponse{}

	if err := c.getJSON(ctx, snapToRoadsAPI, r, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *SnapToRoadRequest) params() url.Values {
	q := make(url.Values)
	var p []string
	for _, e := range r.Path {
		p = append(p, e.String())
	}

	q.Set("path", strings.Join(p, "|"))
	if r.Interpolate {
		q.Set("interpolate", "true")
	}

	return q
}

// SnapToRoadRequest is the request structure for the Roads Snap to Road API.
type SnapToRoadRequest struct {
	// Path is the path to be snapped.
	Path []LatLng

	// Interpolate is whether to interpolate a path to include all points forming the
	// full road-geometry.
	Interpolate bool
}

// SnapToRoadResponse is an array of snapped points.
type SnapToRoadResponse struct {
	SnappedPoints []SnappedPoint `json:"snappedPoints"`
}

// SnappedPoint is the original path point snapped to a road.
type SnappedPoint struct {
	// Location of the snapped point.
	Location LatLng `json:"location"`

	// OriginalIndex is an integer that indicates the corresponding value in the
	// original request. Not present on interpolated points.
	OriginalIndex *int `json:"originalIndex"`

	// PlaceID is a unique identifier for a place.
	PlaceID string `json:"placeId"`
}

// NearestRoads makes a Nearest Roads API request
func (c *Client) NearestRoads(ctx context.Context, r *NearestRoadsRequest) (*NearestRoadsResponse, error) {

	if len(r.Points) == 0 {
		return nil, errors.New("maps: Points empty")
	}

	response := &NearestRoadsResponse{}

	if err := c.getJSON(ctx, nearestRoadsAPI, r, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *NearestRoadsRequest) params() url.Values {
	q := make(url.Values)
	var p []string
	for _, e := range r.Points {
		p = append(p, e.String())
	}

	q.Set("points", strings.Join(p, "|"))

	return q
}

// NearestRoadsRequest is the request structure for the Nearest Roads API.
type NearestRoadsRequest struct {
	// Points is the list of points to be snapped.
	Points []LatLng
}

// NearestRoadsResponse is an array of snapped points.
type NearestRoadsResponse struct {
	SnappedPoints []SnappedPoint `json:"snappedPoints"`
}

// SpeedLimits makes a Speed Limits API request
func (c *Client) SpeedLimits(ctx context.Context, r *SpeedLimitsRequest) (*SpeedLimitsResponse, error) {

	if len(r.Path) == 0 && len(r.PlaceID) == 0 {
		return nil, errors.New("maps: Path and PlaceID both empty")
	}

	response := &SpeedLimitsResponse{}

	if err := c.getJSON(ctx, speedLimitsAPI, r, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *SpeedLimitsRequest) params() url.Values {
	q := make(url.Values)

	var p []string
	for _, e := range r.Path {
		p = append(p, e.String())
	}

	if len(p) > 0 {
		q.Set("path", strings.Join(p, "|"))
	}
	for _, id := range r.PlaceID {
		q.Add("placeId", id)
	}
	if r.Units != "" {
		q.Set("units", string(r.Units))
	}

	return q
}

type speedLimitUnit string

const (
	// SpeedLimitMPH is for requesting speed limits in Miles Per Hour.
	SpeedLimitMPH = "MPH"
	// SpeedLimitKPH is for requesting speed limits in Kilometers Per Hour.
	SpeedLimitKPH = "KPH"
)

// SpeedLimitsRequest is the request structure for the Roads Speed Limits API.
type SpeedLimitsRequest struct {
	// Path is the path to be snapped and speed limits requested.
	Path []LatLng

	// PlaceID is the PlaceIDs to request speed limits for.
	PlaceID []string

	// Units is whether to return speed limits in `SpeedLimitKPH` or `SpeedLimitMPH`.
	// Optional, default behavior is to return results in KPH.
	Units speedLimitUnit
}

// SpeedLimitsResponse is an array of snapped points and an array of speed limits.
type SpeedLimitsResponse struct {
	SpeedLimits   []SpeedLimit   `json:"speedLimits"`
	SnappedPoints []SnappedPoint `json:"snappedPoints"`
}

// SpeedLimit is the speed limit for a PlaceID
type SpeedLimit struct {
	// PlaceID is a unique identifier for a place.
	PlaceID string `json:"placeId"`
	// SpeedLimit is the speed limit for that road segment.
	SpeedLimit float64 `json:"speedLimit"`
	// Units is either KPH or MPH.
	Units speedLimitUnit `json:"units"`
}
