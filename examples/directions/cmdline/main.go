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

// Package main contains a simple command line tool for Directions
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/maps"
	"google.golang.org/maps/directions"
)

var (
	apiKey      = flag.String("key", "", "API Key for using Google Maps API.")
	origin      = flag.String("origin", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	destination = flag.String("destination", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
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
	if *origin == "" {
		usageAndExit("Please specify an origin.")
	}
	if *destination == "" {
		usageAndExit("Please specify a destination.")
	}
	ctx := maps.NewContext(*apiKey, client)
	resp, err := directions.Get(ctx, *origin, *destination)
	if err != nil {
		log.Fatalf("Could not request directions: %v", err)
	}

	if len(resp.Routes) == 0 {
		log.Fatalf("No results")
	}
	route := resp.Routes[0]

	fmt.Println("Summary:", route.Summary)
	fmt.Printf("Bounds NorthEast lat/lng: %f,%f\n", route.Bounds.NorthEast.Lat, route.Bounds.NorthEast.Lng)
	fmt.Printf("Bounds SouthWest lat/lng: %f,%f\n", route.Bounds.SouthWest.Lat, route.Bounds.SouthWest.Lng)
	fmt.Println("Copyrights:", route.Copyrights)

	for idx, leg := range route.Legs {
		fmt.Println("Leg", idx, "distance:", leg.Distance)
		fmt.Println("Leg", idx, "duration:", leg.Duration)
		fmt.Println("Leg", idx, "start address:", leg.StartAddress)
		fmt.Println("Leg", idx, "start location:", leg.StartLocation)
		fmt.Println("Leg", idx, "end address:", leg.EndAddress)
		fmt.Println("Leg", idx, "end location:", leg.EndLocation)

		for idx, step := range leg.Steps {
			fmt.Println("Step", idx, "distance:", step.Distance)
			fmt.Println("Step", idx, "duration:", step.Duration)
			fmt.Println("Step", idx, "start location:", step.StartLocation)
			fmt.Println("Step", idx, "end location:", step.EndLocation)
			fmt.Println("Step", idx, "travel mode:", step.TravelMode)
			fmt.Println("Step", idx, "path:", step.Polyline.Decode())
		}
	}

}
