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
	"net/http"
	"os"
	"strings"

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
	client := &http.Client{}
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	ctx := maps.NewContext(*apiKey, client)

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

	pretty.Println(r)

	resp, err := r.Get(ctx)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(resp)
}

func lookupMode(mode string, r *maps.DistanceMatrixRequest) {
	if mode != "" {
		switch {
		case mode == "driving":
			r.Mode = maps.TravelModeDriving
		case mode == "walking":
			r.Mode = maps.TravelModeWalking
		case mode == "bicycling":
			r.Mode = maps.TravelModeBicycling
		case mode == "transit":
			r.Mode = maps.TravelModeTransit
		default:
			log.Fatalf("Unknown mode %s", mode)
		}
	}
}

func lookupAvoid(avoid string, r *maps.DistanceMatrixRequest) {
	if avoid != "" {
		switch {
		case avoid == "tolls":
			r.Avoid = maps.AvoidTolls
		case avoid == "highways":
			r.Avoid = maps.AvoidHighways
		case avoid == "ferries":
			r.Avoid = maps.AvoidFerries
		default:
			log.Fatalf("Unknown avoid restriction %s", avoid)
		}
	}
}

func lookupUnits(units string, r *maps.DistanceMatrixRequest) {
	if units != "" {
		switch {
		case units == "metric":
			r.Units = maps.UnitsMetric
		case units == "imperial":
			r.Units = maps.UnitsImperial
		default:
			log.Fatalf("Unknown units %s", units)
		}
	}
}

func lookupTransitMode(transitMode string, r *maps.DistanceMatrixRequest) {
	if transitMode != "" {
		switch {
		case transitMode == "bus":
			r.TransitMode = maps.TransitModeBus
		case transitMode == "subway":
			r.TransitMode = maps.TransitModeSubway
		case transitMode == "train":
			r.TransitMode = maps.TransitModeTrain
		case transitMode == "tram":
			r.TransitMode = maps.TransitModeTram
		case transitMode == "rail":
			r.TransitMode = maps.TransitModeRail
		default:
			log.Fatalf("Unknown transit_mode %s", transitMode)
		}
	}
}

func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DistanceMatrixRequest) {
	if transitRoutingPreference != "" {
		switch {
		case transitRoutingPreference == "fewer_transfers":
			r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
		case transitRoutingPreference == "less_walking":
			r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
		default:
			log.Fatalf("Unknown transit_routing_preference %s", transitRoutingPreference)
		}
	}
}
