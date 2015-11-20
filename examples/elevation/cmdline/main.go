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
	"os"
	"strings"

	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	clientID  = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature = flag.String("signature", "", "Signature for Maps for Work API access.")
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

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func main() {
	flag.Parse()

	var client *maps.Client
	var err error
	if *apiKey != "" {
		client, err = maps.NewClient(maps.WithAPIKey(*apiKey))
	} else if *clientID != "" || *signature != "" {
		client, err = maps.NewClient(maps.WithClientIDAndSignature(*clientID, *signature))
	} else {
		usageAndExit("Please specify an API Key, or Client ID and Signature.")
	}
	check(err)

	r := &maps.ElevationRequest{}

	if *samples > 0 {
		r.Samples = *samples
	}

	if *locations != "" {
		l, err := decodeLocations(*locations)
		check(err)
		r.Locations = l
	}

	if *path != "" {
		p, err := decodePath(*path)
		check(err)
		r.Path = p
	}

	resp, err := client.Elevation(context.Background(), r)
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

	if strings.Contains(location, "|") {
		return maps.ParseLatLngList(location)
	}

	// single location
	ll, err := maps.ParseLatLng(location)
	check(err)
	return []maps.LatLng{ll}, nil
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
			ll, err := maps.ParseLatLng(l)
			check(err)
			result = append(result, ll)
		}
		return result, nil
	}
	return result, fmt.Errorf("Invalid Path argument: '%s'", path)
}
