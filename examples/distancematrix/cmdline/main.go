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

// Package main contains a simple command line tool for DistanceMatrix
// Directions docs: https://developers.google.com/maps/documentation/distancematrix/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey                   = flag.String("key", "", "API Key for using Google Maps API.")
	origins                  = flag.String("origins", "", "One or more addresses and/or textual latitude/longitude values, separated with the pipe (|) character, from which to calculate distance and time.")
	destinations             = flag.String("destinations", "", "One or more addresses and/or textual latitude/longitude values, separated with the pipe (|) character, to which to calculate distance and time.")
	mode                     = flag.String("mode", "", "Specifies the mode of transport to use when calculating distance.")
	language                 = flag.String("language", "", "The language in which to return results.")
	avoid                    = flag.String("avoid", "", "Introduces restrictions to the route.")
	units                    = flag.String("units", "", "Specifies the unit system to use when expressing distance as text.")
	departureTime            = flag.String("departure_time", "", "The desired time of departure.")
	arrivalTime              = flag.String("arrival_time", "", "Specifies the desired time of arrival.")
	transitMode              = flag.String("transit_mode", "", "Specifies one or more preferred modes of transit.")
	transitRoutingPreference = flag.String("transit_routing_preference", "", "Specifies preferences for transit requests.")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	client, err := maps.NewClient(maps.WithAPIKey(*apiKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.DistanceMatrixRequest{
		Language:      *language,
		DepartureTime: *departureTime,
		ArrivalTime:   *arrivalTime,
	}

	if *origins != "" {
		r.Origins = strings.Split(*origins, "|")
	}
	if *destinations != "" {
		r.Destinations = strings.Split(*destinations, "|")
	}

	lookupMode(*mode, r)
	lookupAvoid(*avoid, r)
	lookupUnits(*units, r)
	lookupTransitMode(*transitMode, r)
	lookupTransitRoutingPreference(*transitRoutingPreference, r)

	resp, err := client.GetDistanceMatrix(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(resp)
}

func lookupMode(mode string, r *maps.DistanceMatrixRequest) {
	switch mode {
	case "driving":
		r.Mode = maps.TravelModeDriving
	case "walking":
		r.Mode = maps.TravelModeWalking
	case "bicycling":
		r.Mode = maps.TravelModeBicycling
	case "transit":
		r.Mode = maps.TravelModeTransit
	case "":
		// ignore
	default:
		log.Fatalf("Unknown mode %s", mode)
	}
}

func lookupAvoid(avoid string, r *maps.DistanceMatrixRequest) {
	switch avoid {
	case "tolls":
		r.Avoid = maps.AvoidTolls
	case "highways":
		r.Avoid = maps.AvoidHighways
	case "ferries":
		r.Avoid = maps.AvoidFerries
	case "":
		//ignore
	default:
		log.Fatalf("Unknown avoid restriction %s", avoid)
	}
}

func lookupUnits(units string, r *maps.DistanceMatrixRequest) {
	switch units {
	case "metric":
		r.Units = maps.UnitsMetric
	case "imperial":
		r.Units = maps.UnitsImperial
	case "":
		// ignore
	default:
		log.Fatalf("Unknown units %s", units)
	}
}

func lookupTransitMode(transitMode string, r *maps.DistanceMatrixRequest) {
	switch transitMode {
	case "bus":
		r.TransitMode = maps.TransitModeBus
	case "subway":
		r.TransitMode = maps.TransitModeSubway
	case "train":
		r.TransitMode = maps.TransitModeTrain
	case "tram":
		r.TransitMode = maps.TransitModeTram
	case "rail":
		r.TransitMode = maps.TransitModeRail
	case "":
		// ignore
	default:
		log.Fatalf("Unknown transit_mode %s", transitMode)
	}
}

func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DistanceMatrixRequest) {
	switch transitRoutingPreference {
	case "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	case "":
		// ignore
	default:
		log.Fatalf("Unknown transit_routing_preference %s", transitRoutingPreference)
	}
}
