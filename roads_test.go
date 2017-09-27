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

func TestSnapToRoad(t *testing.T) {

	response := `{
  "snappedPoints": [
    {
      "location": {
        "latitude": -35.2784167,
        "longitude": 149.1294692
      },
      "originalIndex": 0,
      "placeId": "ChIJoR7CemhNFmsRQB9QbW7qABM"
    },
    {
      "location": {
        "latitude": -35.280321693840129,
        "longitude": 149.12908274880189
      },
      "originalIndex": 1,
      "placeId": "ChIJiy6YT2hNFmsRkHZAbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2803415,
        "longitude": 149.1290788
      },
      "placeId": "ChIJiy6YT2hNFmsRkHZAbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2803415,
        "longitude": 149.1290788
      },
      "placeId": "ChIJI2FUTGhNFmsRcHpAbW7qABM"
    },
    {
      "location": {
        "latitude": -35.280451499999991,
        "longitude": 149.1290784
      },
      "placeId": "ChIJI2FUTGhNFmsRcHpAbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2805167,
        "longitude": 149.1290879
      },
      "placeId": "ChIJI2FUTGhNFmsRcHpAbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2805901,
        "longitude": 149.1291104
      },
      "placeId": "ChIJI2FUTGhNFmsRcHpAbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2805901,
        "longitude": 149.1291104
      },
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM"
    },
    {
      "location": {
        "latitude": -35.280734599999995,
        "longitude": 149.1291517
      },
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2807852,
        "longitude": 149.1291716
      },
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM"
    },
    {
      "location": {
        "latitude": -35.2808499,
        "longitude": 149.1292099
      },
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM"
    },
    {
      "location": {
        "latitude": -35.280923099999995,
        "longitude": 149.129278
      },
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM"
    },
    {
      "location": {
        "latitude": -35.280960897210818,
        "longitude": 149.1293250692261
      },
      "originalIndex": 2,
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM"
    },
    {
      "location": {
        "latitude": -35.284728724835304,
        "longitude": 149.12835061713685
      },
      "originalIndex": 7,
      "placeId": "ChIJW5JAZmpNFmsRegG0-Jc80sM"
    }
  ]
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &SnapToRoadRequest{
		Path: []LatLng{
			{Lat: -35.27801, Lng: 149.12958},
			{Lat: -35.28032, Lng: 149.12907},
			{Lat: -35.28099, Lng: 149.12929},
			{Lat: -35.28144, Lng: 149.12984},
			{Lat: -35.28194, Lng: 149.13003},
			{Lat: -35.28282, Lng: 149.12956},
			{Lat: -35.28302, Lng: 149.12881},
			{Lat: -35.28473, Lng: 149.12836},
		},
	}

	resp, err := c.SnapToRoad(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	// Required because we can't do &2 in the data structure below.
	index0 := 0
	index1 := 1
	index2 := 2
	index7 := 7

	correctResponse := &SnapToRoadResponse{
		SnappedPoints: []SnappedPoint{
			{
				Location:      LatLng{Lat: -35.2784167, Lng: 149.1294692},
				OriginalIndex: &index0,
				PlaceID:       "ChIJoR7CemhNFmsRQB9QbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.28032169384013, Lng: 149.1290827488019},
				OriginalIndex: &index1,
				PlaceID:       "ChIJiy6YT2hNFmsRkHZAbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2803415, Lng: 149.1290788},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJiy6YT2hNFmsRkHZAbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2803415, Lng: 149.1290788},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJI2FUTGhNFmsRcHpAbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.28045149999999, Lng: 149.1290784},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJI2FUTGhNFmsRcHpAbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2805167, Lng: 149.1290879},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJI2FUTGhNFmsRcHpAbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2805901, Lng: 149.1291104},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJI2FUTGhNFmsRcHpAbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2805901, Lng: 149.1291104},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJW9R7smlNFmsRMH1AbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.280734599999995, Lng: 149.1291517},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJW9R7smlNFmsRMH1AbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2807852, Lng: 149.1291716},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJW9R7smlNFmsRMH1AbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.2808499, Lng: 149.1292099},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJW9R7smlNFmsRMH1AbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.280923099999995, Lng: 149.129278},
				OriginalIndex: (*int)(nil),
				PlaceID:       "ChIJW9R7smlNFmsRMH1AbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.28096089721082, Lng: 149.1293250692261},
				OriginalIndex: &index2,
				PlaceID:       "ChIJW9R7smlNFmsRMH1AbW7qABM",
			},
			{
				Location:      LatLng{Lat: -35.284728724835304, Lng: 149.12835061713685},
				OriginalIndex: &index7,
				PlaceID:       "ChIJW5JAZmpNFmsRegG0-Jc80sM",
			},
		},
	}

	if !reflect.DeepEqual(resp, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp)
	}
}

func TestSnapToRoadNoPath(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &SnapToRoadRequest{}

	if _, err := c.SnapToRoad(context.Background(), r); err == nil {
		t.Errorf("Empty path should return error")
	}
}

func TestSnapToRoadWithCancelledContext(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &SnapToRoadRequest{
		Path: []LatLng{
			{Lat: -35.27801, Lng: 149.12958},
			{Lat: -35.28032, Lng: 149.12907},
			{Lat: -35.28099, Lng: 149.12929},
			{Lat: -35.28144, Lng: 149.12984},
			{Lat: -35.28194, Lng: 149.13003},
			{Lat: -35.28282, Lng: 149.12956},
			{Lat: -35.28302, Lng: 149.12881},
			{Lat: -35.28473, Lng: 149.12836},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.SnapToRoad(ctx, r); err == nil {
		t.Errorf("Cancelled context should return non-nil err")
	}
}

func TestSpeedLimit(t *testing.T) {
	response := `{
  "speedLimits": [
    {
      "placeId": "ChIJ1Wi6I2pNFmsRQL9GbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJ58xCoGlNFmsRUEZUbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJ9RhaiGlNFmsR0IxAbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJabjuhGlNFmsREIxAbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJcSAlFWpNFmsRMHlUbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJI2FUTGhNFmsRcHpAbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJiy6YT2hNFmsRkHZAbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJoR7CemhNFmsRQB9QbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJP2m_FWpNFmsRIHlUbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJtV7La2pNFmsRAGpHbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJW5JAZmpNFmsRegG0-Jc80sM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJW9R7smlNFmsRMH1AbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    },
    {
      "placeId": "ChIJy8c0r2lNFmsRQEZUbW7qABM",
      "speedLimit": 60,
      "units": "KPH"
    }
  ]
}`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &SpeedLimitsRequest{
		PlaceID: []string{
			"ChIJ1Wi6I2pNFmsRQL9GbW7qABM",
			"ChIJ58xCoGlNFmsRUEZUbW7qABM",
			"ChIJ9RhaiGlNFmsR0IxAbW7qABM",
			"ChIJabjuhGlNFmsREIxAbW7qABM",
			"ChIJcSAlFWpNFmsRMHlUbW7qABM",
			"ChIJI2FUTGhNFmsRcHpAbW7qABM",
			"ChIJiy6YT2hNFmsRkHZAbW7qABM",
			"ChIJoR7CemhNFmsRQB9QbW7qABM",
			"ChIJP2m_FWpNFmsRIHlUbW7qABM",
			"ChIJtV7La2pNFmsRAGpHbW7qABM",
			"ChIJW5JAZmpNFmsRegG0-Jc80sM",
			"ChIJW9R7smlNFmsRMH1AbW7qABM",
			"ChIJy8c0r2lNFmsRQEZUbW7qABM",
		},
	}

	resp, err := c.SpeedLimits(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := &SpeedLimitsResponse{
		SpeedLimits: []SpeedLimit{
			{PlaceID: "ChIJ1Wi6I2pNFmsRQL9GbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJ58xCoGlNFmsRUEZUbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJ9RhaiGlNFmsR0IxAbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJabjuhGlNFmsREIxAbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJcSAlFWpNFmsRMHlUbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJI2FUTGhNFmsRcHpAbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJiy6YT2hNFmsRkHZAbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJoR7CemhNFmsRQB9QbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJP2m_FWpNFmsRIHlUbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJtV7La2pNFmsRAGpHbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJW5JAZmpNFmsRegG0-Jc80sM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJW9R7smlNFmsRMH1AbW7qABM", SpeedLimit: 60, Units: "KPH"},
			{PlaceID: "ChIJy8c0r2lNFmsRQEZUbW7qABM", SpeedLimit: 60, Units: "KPH"},
		},
		SnappedPoints: nil,
	}

	if !reflect.DeepEqual(resp, correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp)
	}
}

func TestSpeedLimitsNoPlaceIDs(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &SpeedLimitsRequest{}

	if _, err := c.SpeedLimits(context.Background(), r); err == nil {
		t.Errorf("Empty PlaceIDs should return error")
	}
}

func TestSpeedLimitsWithCancelledContext(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &SpeedLimitsRequest{
		PlaceID: []string{
			"ChIJ1Wi6I2pNFmsRQL9GbW7qABM",
			"ChIJ58xCoGlNFmsRUEZUbW7qABM",
			"ChIJ9RhaiGlNFmsR0IxAbW7qABM",
			"ChIJabjuhGlNFmsREIxAbW7qABM",
			"ChIJcSAlFWpNFmsRMHlUbW7qABM",
			"ChIJI2FUTGhNFmsRcHpAbW7qABM",
			"ChIJiy6YT2hNFmsRkHZAbW7qABM",
			"ChIJoR7CemhNFmsRQB9QbW7qABM",
			"ChIJP2m_FWpNFmsRIHlUbW7qABM",
			"ChIJtV7La2pNFmsRAGpHbW7qABM",
			"ChIJW5JAZmpNFmsRegG0-Jc80sM",
			"ChIJW9R7smlNFmsRMH1AbW7qABM",
			"ChIJy8c0r2lNFmsRQEZUbW7qABM",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.SpeedLimits(ctx, r); err == nil {
		t.Errorf("Cancelled context should return non-nil err")
	}
}

func TestSnapToRoadRequestURL(t *testing.T) {
	expectedQuery := "interpolate=true&key=AIzaNotReallyAnAPIKey&path=1%2C2%7C3%2C4"

	server := mockServerForQuery(expectedQuery, 200, `{}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &SnapToRoadRequest{
		Path:        []LatLng{{1, 2}, {3, 4}},
		Interpolate: true,
	}

	_, err := c.SnapToRoad(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestSpeedLimitsRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&path=-35.27801%2C149.12958%7C-35.28032%2C149.12907&placeId=ChIJ1Wi6I2pNFmsRQL9GbW7qABM&placeId=ChIJ58xCoGlNFmsRUEZUbW7qABM&units=MPH"

	server := mockServerForQuery(expectedQuery, 200, `{}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &SpeedLimitsRequest{
		Path: []LatLng{
			{Lat: -35.27801, Lng: 149.12958},
			{Lat: -35.28032, Lng: 149.12907},
		},
		PlaceID: []string{
			"ChIJ1Wi6I2pNFmsRQL9GbW7qABM",
			"ChIJ58xCoGlNFmsRUEZUbW7qABM",
		},
		Units: SpeedLimitMPH,
	}

	_, err := c.SpeedLimits(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}
