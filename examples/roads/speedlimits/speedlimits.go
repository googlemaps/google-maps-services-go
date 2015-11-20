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

// Package main contains a simple command line tool for Timezone API
// Directions docs: https://developers.google.com/maps/documentation/timezone/
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
	apiKey   = flag.String("key", "", "API Key for using Google Maps API.")
	path     = flag.String("path", "", "The path to be snapped. The path parameter accepts a list of latitude/longitude pairs. Latitude and longitude values should be separated by commas. Coordinates should be separated by the pipe character.")
	placeIDs = flag.String("place_ids", "", "The place ID of the road segment. Place IDs are returned by the snapToRoads method. You can pass up to 100 Place IDs with each request. Place IDs should be separated by a comma.")
	units    = flag.String("units", "", "Whether to return speed limits in kilometers or miles per hour. This can be set to either KPH or MPH. Defaults to KPH.")
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
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	client, err := maps.NewClient(maps.WithAPIKey(*apiKey))
	if err != nil {
		log.Fatalf("error %v", err)
	}
	r := &maps.SpeedLimitsRequest{}

	if *units == "KPH" {
		r.Units = maps.SpeedLimitKPH
	}
	if *units == "MPH" {
		r.Units = maps.SpeedLimitMPH
	}

	if *path == "" && *placeIDs == "" {
		usageAndExit("Please specify either a path to be snapped, or a list of Place IDs.")
	}
	parsePath(*path, r)
	parsePlaceIDs(*placeIDs, r)

	resp, err := client.SpeedLimits(context.Background(), r)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	pretty.Println(resp)
}

// parsePath takes a location argument string and decodes it.
func parsePath(path string, r *maps.SpeedLimitsRequest) {
	if path != "" {
		ls := strings.Split(path, "|")
		for _, l := range ls {
			ll, err := maps.ParseLatLng(l)
			check(err)
			r.Path = append(r.Path, ll)
		}
	}
}

// parsePlacesIds takes a placesIds argument string and decodes it.
func parsePlaceIDs(placeIDs string, r *maps.SpeedLimitsRequest) {
	if placeIDs != "" {
		ids := strings.Split(placeIDs, ",")
		for _, id := range ids {
			r.PlaceID = append(r.PlaceID, id)
		}
	}
}
