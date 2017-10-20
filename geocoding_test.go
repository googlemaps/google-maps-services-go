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
	"net/url"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestGeocodingGoogleHQ(t *testing.T) {
	response := `{
    "results": [
        {
            "address_components": [
                {
                    "long_name": "1600",
                    "short_name": "1600",
                    "types": [
                        "street_number"
                    ]
                },
                {
                    "long_name": "Amphitheatre Pkwy",
                    "short_name": "Amphitheatre Pkwy",
                    "types": [
                        "route"
                    ]
                },
                {
                    "long_name": "Mountain View",
                    "short_name": "Mountain View",
                    "types": [
                        "locality",
                        "political"
                    ]
                },
                {
                    "long_name": "Santa Clara County",
                    "short_name": "Santa Clara County",
                    "types": [
                        "administrative_area_level_2",
                        "political"
                    ]
                },
                {
                    "long_name": "California",
                    "short_name": "CA",
                    "types": [
                        "administrative_area_level_1",
                        "political"
                    ]
                },
                {
                    "long_name": "United States",
                    "short_name": "US",
                    "types": [
                        "country",
                        "political"
                    ]
                },
                {
                    "long_name": "94043",
                    "short_name": "94043",
                    "types": [
                        "postal_code"
                    ]
                }
            ],
            "formatted_address": "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA",
            "geometry": {
                "location": {
                    "lat": 37.4224764,
                    "lng": -122.0842499
                },
                "bounds": {
                    "northeast": {
                        "lat": 37.4238253802915,
                        "lng": -122.0829009197085
                    },
                    "southwest": {
                        "lat": 37.4211274197085,
                        "lng": -122.0855988802915
                    }
                },
                "location_type": "ROOFTOP",
                "viewport": {
                    "northeast": {
                        "lat": 37.4238253802915,
                        "lng": -122.0829009197085
                    },
                    "southwest": {
                        "lat": 37.4211274197085,
                        "lng": -122.0855988802915
                    }
                }
            },
            "place_id": "ChIJ2eUgeAK6j4ARbn5u_wAGqWA",
            "types": [
                "street_address"
            ]
        }
    ],
    "status": "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &GeocodingRequest{
		Address: "1600 Amphitheatre Parkway, Mountain View, CA",
	}

	resp, err := c.Geocode(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := GeocodingResult{
		AddressComponents: []AddressComponent{
			{
				LongName:  "1600",
				ShortName: "1600",
				Types:     []string{"street_number"},
			},
			{
				LongName:  "Amphitheatre Pkwy",
				ShortName: "Amphitheatre Pkwy",
				Types:     []string{"route"},
			},
			{
				LongName:  "Mountain View",
				ShortName: "Mountain View",
				Types:     []string{"locality", "political"},
			},
			{
				LongName:  "Santa Clara County",
				ShortName: "Santa Clara County",
				Types:     []string{"administrative_area_level_2", "political"},
			},
			{
				LongName:  "California",
				ShortName: "CA",
				Types:     []string{"administrative_area_level_1", "political"},
			},
			{
				LongName:  "United States",
				ShortName: "US",
				Types:     []string{"country", "political"},
			},
			{
				LongName:  "94043",
				ShortName: "94043",
				Types:     []string{"postal_code"},
			},
		},
		FormattedAddress: "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA",
		Geometry: AddressGeometry{
			Location: LatLng{Lat: 37.4224764, Lng: -122.0842499},
			Bounds: LatLngBounds{
				NorthEast: LatLng{Lat: 37.4238253802915, Lng: -122.0829009197085},
				SouthWest: LatLng{Lat: 37.4211274197085, Lng: -122.0855988802915},
			},
			LocationType: "ROOFTOP",
			Viewport: LatLngBounds{
				NorthEast: LatLng{Lat: 37.4238253802915, Lng: -122.0829009197085},
				SouthWest: LatLng{Lat: 37.4211274197085, Lng: -122.0855988802915},
			},
			Types: nil,
		},
		PlaceID: "ChIJ2eUgeAK6j4ARbn5u_wAGqWA",
		Types:   []string{"street_address"},
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp[0])
	}
}

func TestGeocodingReverseGeocoding(t *testing.T) {

	response := `{
    "results": [
        {
            "address_components": [
                {
                    "long_name": "277",
                    "short_name": "277",
                    "types": [
                        "street_number"
                    ]
                },
                {
                    "long_name": "Bedford Avenue",
                    "short_name": "Bedford Ave",
                    "types": [
                        "route"
                    ]
                },
                {
                    "long_name": "Williamsburg",
                    "short_name": "Williamsburg",
                    "types": [
                        "neighborhood",
                        "political"
                    ]
                },
                {
                    "long_name": "Brooklyn",
                    "short_name": "Brooklyn",
                    "types": [
                        "sublocality",
                        "political"
                    ]
                },
                {
                    "long_name": "Kings",
                    "short_name": "Kings",
                    "types": [
                        "administrative_area_level_2",
                        "political"
                    ]
                },
                {
                    "long_name": "New York",
                    "short_name": "NY",
                    "types": [
                        "administrative_area_level_1",
                        "political"
                    ]
                },
                {
                    "long_name": "United States",
                    "short_name": "US",
                    "types": [
                        "country",
                        "political"
                    ]
                },
                {
                    "long_name": "11211",
                    "short_name": "11211",
                    "types": [
                        "postal_code"
                    ]
                }
            ],
            "formatted_address": "277 Bedford Avenue, Brooklyn, NY 11211, USA",
            "geometry": {
                "location": {
                    "lat": 40.714232,
                    "lng": -73.9612889
                },
                "bounds": {
                    "northeast": {
                        "lat": 40.7155809802915,
                        "lng": -73.9599399197085
                    },
                    "southwest": {
                        "lat": 40.7128830197085,
                        "lng": -73.96263788029151
                    }
                },
                "location_type": "ROOFTOP",
                "viewport": {
                    "northeast": {
                        "lat": 40.7155809802915,
                        "lng": -73.9599399197085
                    },
                    "southwest": {
                        "lat": 40.7128830197085,
                        "lng": -73.96263788029151
                    }
                }
            },
            "place_id": "ChIJd8BlQ2BZwokRAFUEcm_qrcA",
            "types": [
                "street_address"
            ]
        }
    ],
    "status": "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &GeocodingRequest{
		LatLng: &LatLng{Lat: 40.714224, Lng: -73.961452},
	}

	resp, err := c.ReverseGeocode(context.Background(), r)

	if len(resp) != 1 {
		t.Errorf("expected %+v, was %+v", 1, len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := GeocodingResult{
		AddressComponents: []AddressComponent{
			{
				LongName:  "277",
				ShortName: "277",
				Types:     []string{"street_number"},
			},
			{
				LongName:  "Bedford Avenue",
				ShortName: "Bedford Ave",
				Types:     []string{"route"},
			},
			{
				LongName:  "Williamsburg",
				ShortName: "Williamsburg",
				Types:     []string{"neighborhood", "political"},
			},
			{
				LongName:  "Brooklyn",
				ShortName: "Brooklyn",
				Types:     []string{"sublocality", "political"},
			},
			{
				LongName:  "Kings",
				ShortName: "Kings",
				Types:     []string{"administrative_area_level_2", "political"},
			},
			{
				LongName:  "New York",
				ShortName: "NY",
				Types:     []string{"administrative_area_level_1", "political"},
			},
			{
				LongName:  "United States",
				ShortName: "US",
				Types:     []string{"country", "political"},
			},
			{
				LongName:  "11211",
				ShortName: "11211",
				Types:     []string{"postal_code"},
			},
		},
		FormattedAddress: "277 Bedford Avenue, Brooklyn, NY 11211, USA",
		Geometry: AddressGeometry{
			Location: LatLng{Lat: 40.714232, Lng: -73.9612889},
			Bounds: LatLngBounds{
				NorthEast: LatLng{Lat: 40.7155809802915, Lng: -73.9599399197085},
				SouthWest: LatLng{Lat: 40.7128830197085, Lng: -73.96263788029151},
			},
			LocationType: "ROOFTOP",
			Viewport: LatLngBounds{
				NorthEast: LatLng{Lat: 40.7155809802915, Lng: -73.9599399197085},
				SouthWest: LatLng{Lat: 40.7128830197085, Lng: -73.96263788029151},
			},
			Types: nil,
		},
		PlaceID: "ChIJd8BlQ2BZwokRAFUEcm_qrcA",
		Types:   []string{"street_address"},
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp[0])
	}
}

func TestGeocodingEmptyRequest(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &GeocodingRequest{}

	if _, err := c.Geocode(context.Background(), r); err == nil {
		t.Errorf("Missing Address, Address Components, and LatLng should return error")
	}
}

func TestGeocodingWithCancelledContext(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &GeocodingRequest{
		Address: "Sydney Town Hall",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.Geocode(ctx, r); err == nil {
		t.Errorf("Cancelled context should return non-nil err")
	}
}

func TestGeocodingFailingServer(t *testing.T) {
	server := mockServer(500, `{"status" : "ERROR"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &GeocodingRequest{
		Address: "Sydney Town Hall",
	}

	if _, err := c.Geocode(context.Background(), r); err == nil {
		t.Errorf("Failing server should return error")
	}
}

func TestGeocodingRequestURL(t *testing.T) {
	expectedQuery := "address=Santa+Cruz&bounds=34.236144%2C-118.500938%7C34.172684%2C-118.604794&components=country%3AES&key=AIzaNotReallyAnAPIKey&language=es&location_type=APPROXIMATE&region=es&result_type=country"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &GeocodingRequest{
		Address:      "Santa Cruz",
		Bounds:       &LatLngBounds{LatLng{34.172684, -118.604794}, LatLng{34.236144, -118.500938}},
		Region:       "es",
		ResultType:   []string{"country"},
		LocationType: []GeocodeAccuracy{GeocodeAccuracyApproximate},
		Components:   map[Component]string{ComponentCountry: "ES"},
		Language:     "es",
	}

	_, err := c.Geocode(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestReverseGeocodingPlaceID(t *testing.T) {
	response := `{
    "results": [
        {
            "address_components": [
                {
                    "long_name": "1600",
                    "short_name": "1600",
                    "types": [
                        "street_number"
                    ]
                },
                {
                    "long_name": "Amphitheatre Pkwy",
                    "short_name": "Amphitheatre Pkwy",
                    "types": [
                        "route"
                    ]
                },
                {
                    "long_name": "Mountain View",
                    "short_name": "Mountain View",
                    "types": [
                        "locality",
                        "political"
                    ]
                },
                {
                    "long_name": "Santa Clara County",
                    "short_name": "Santa Clara County",
                    "types": [
                        "administrative_area_level_2",
                        "political"
                    ]
                },
                {
                    "long_name": "California",
                    "short_name": "CA",
                    "types": [
                        "administrative_area_level_1",
                        "political"
                    ]
                },
                {
                    "long_name": "United States",
                    "short_name": "US",
                    "types": [
                        "country",
                        "political"
                    ]
                },
                {
                    "long_name": "94043",
                    "short_name": "94043",
                    "types": [
                        "postal_code"
                    ]
                }
            ],
            "formatted_address": "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA",
            "geometry": {
                "location": {
                    "lat": 37.4224764,
                    "lng": -122.0842499
                },
                "location_type": "ROOFTOP",
                "viewport": {
                    "northeast": {
                        "lat": 37.4238253802915,
                        "lng": -122.0829009197085
                    },
                    "southwest": {
                        "lat": 37.4211274197085,
                        "lng": -122.0855988802915
                    }
                }
            },
            "place_id": "ChIJ2eUgeAK6j4ARbn5u_wAGqWA",
            "types": [
                "street_address"
            ]
        }
    ],
    "status": "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &GeocodingRequest{
		PlaceID: "ChIJ2eUgeAK6j4ARbn5u_wAGqWA",
	}

	resp, err := c.ReverseGeocode(context.Background(), r)
	if len(resp) != 1 {
		t.Errorf("Expected length of response is 1, was %+v", len(resp))
	}
	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	correctResponse := GeocodingResult{
		AddressComponents: []AddressComponent{
			{
				LongName:  "1600",
				ShortName: "1600",
				Types:     []string{"street_number"},
			},
			{
				LongName:  "Amphitheatre Pkwy",
				ShortName: "Amphitheatre Pkwy",
				Types:     []string{"route"},
			},
			{
				LongName:  "Mountain View",
				ShortName: "Mountain View",
				Types:     []string{"locality", "political"},
			},
			{
				LongName:  "Santa Clara County",
				ShortName: "Santa Clara County",
				Types:     []string{"administrative_area_level_2", "political"},
			},
			{
				LongName:  "California",
				ShortName: "CA",
				Types:     []string{"administrative_area_level_1", "political"},
			},
			{
				LongName:  "United States",
				ShortName: "US",
				Types:     []string{"country", "political"},
			},
			{
				LongName:  "94043",
				ShortName: "94043",
				Types:     []string{"postal_code"},
			},
		},
		FormattedAddress: "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA",
		Geometry: AddressGeometry{
			Location:     LatLng{Lat: 37.4224764, Lng: -122.0842499},
			LocationType: "ROOFTOP",
			Viewport: LatLngBounds{
				NorthEast: LatLng{Lat: 37.4238253802915, Lng: -122.0829009197085},
				SouthWest: LatLng{Lat: 37.4211274197085, Lng: -122.0855988802915},
			},
			Types: nil,
		},
		PlaceID: "ChIJ2eUgeAK6j4ARbn5u_wAGqWA",
		Types:   []string{"street_address"},
	}

	if !reflect.DeepEqual(resp[0], correctResponse) {
		t.Errorf("expected %+v, was %+v", correctResponse, resp[0])
	}
}

func TestCustomPassThroughGeocodingURL(t *testing.T) {
	expectedQuery := "address=1600+Amphitheatre+Parkway%2C+Mountain+View%2C+CA&key=AIzaNotReallyAnAPIKey&new_forward_geocoder=true"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))
	custom := make(url.Values)
	custom["new_forward_geocoder"] = []string{"true"}

	r := &GeocodingRequest{
		Address: "1600 Amphitheatre Parkway, Mountain View, CA",
		Custom:  custom,
	}

	_, err := c.Geocode(context.Background(), r)
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}
	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestGeocodingZeroResults(t *testing.T) {
	server := mockServer(200, `{"status" : "ZERO_RESULTS"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &GeocodingRequest{
		Address: "Sydney Town Hall",
	}

	response, err := c.Geocode(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error for ZERO_RESULTS status")
	}

	if response == nil {
		t.Errorf("Unexpected nil response for ZERO_RESULTS status")
	}

	if len(response) != 0 {
		t.Errorf("Unexpected response for ZERO_RESULTS status")
	}
}
