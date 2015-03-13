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

// Package main contains a simple command line tool for Elevation API
// Directions docs: https://developers.google.com/maps/documentation/distancematrix/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	locations = flag.String("locations", "", "defines the location(s) on the earth from which to return elevation data. This parameter takes either a single location as a comma-separated pair or multiple latitude/longitude pairs passed as an array or as an encoded polyline.")
	path      = flag.String("path", "", "defines a path on the earth for which to return elevation data.")
	samples   = flag.Int("samples", 0, "specifies the number of sample points along a path for which to return elevation data.")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	client := &http.Client{}
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	ctx := maps.NewContext(*apiKey, client)
	eReq := &maps.ElevationRequest{}

	if *samples > 0 {
		eReq.Samples = *samples
	}

	if *locations != "" {
		l, err := decodeLocations(*locations)
		if err != nil {
			log.Fatalf("Could not parse locations: %#v", err)
		}
		eReq.Locations = l
	}

	if *path != "" {
		p, err := decodePath(*path)
		if err != nil {
			log.Fatalf("Could not parse path: %#v", err)
		}
		eReq.Path = p
	}

	resp, err := eReq.Get(ctx)
	if err != nil {
		log.Fatalf("Could not request elevations: %v", err)
	}

	pretty.Println(resp)
}

// decodeLocations takes a location argument string and decodes it.
// This argument has three different forms, as per documentation at
// https://developers.google.com/maps/documentation/elevation/#Locations
func decodeLocations(location string) ([]maps.LatLng, error) {
	if strings.HasPrefix(location, "enc:") {
		return maps.DecodePolyline(location[len("enc:"):]), nil
	}
	result := []maps.LatLng{}
	if strings.Contains(location, "|") {
		// | delimited list of locations
		ls := strings.Split(location, "|")
		for _, l := range ls {
			ll := strings.Split(l, ",")
			lat, err := strconv.ParseFloat(ll[0], 64)
			if err != nil {
				return result, err
			}
			lng, err := strconv.ParseFloat(ll[1], 64)
			if err != nil {
				return result, err
			}
			result = append(result, maps.LatLng{Lat: lat, Lng: lng})
		}
		return result, nil
	}

	// single location
	ll := strings.Split(location, ",")
	lat, err := strconv.ParseFloat(ll[0], 64)
	if err != nil {
		return result, err
	}
	lng, err := strconv.ParseFloat(ll[1], 64)
	if err != nil {
		return result, err
	}
	result = append(result, maps.LatLng{Lat: lat, Lng: lng})
	return result, nil
}

// decodePath takes a location argument string and decodes it.
// This argument has two different forms, as per documentation at
// https://developers.google.com/maps/documentation/elevation/#Paths
func decodePath(path string) ([]maps.LatLng, error) {
	if strings.HasPrefix(path, "enc:") {
		return maps.DecodePolyline(path[len("enc:"):]), nil
	}
	result := []maps.LatLng{}
	if strings.Contains(path, "|") {
		// | delimited list of locations
		ls := strings.Split(path, "|")
		for _, l := range ls {
			ll := strings.Split(l, ",")
			lat, err := strconv.ParseFloat(ll[0], 64)
			if err != nil {
				return result, err
			}
			lng, err := strconv.ParseFloat(ll[1], 64)
			if err != nil {
				return result, err
			}
			result = append(result, maps.LatLng{Lat: lat, Lng: lng})
		}
		return result, nil
	}
	return result, fmt.Errorf("Invalid Path argument: '%s'", path)
}
