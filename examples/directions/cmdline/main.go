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
	"strings"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey                   = flag.String("key", "", "API Key for using Google Maps API.")
	origin                   = flag.String("origin", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	destination              = flag.String("destination", "", "The address or textual latitude/longitude value from which you wish to calculate directions.")
	mode                     = flag.String("mode", "", "The travel mode for this directions request.")
	departureTime            = flag.String("departure_time", "", "The depature time for transit mode directions request.")
	arrivalTime              = flag.String("arrival_time", "", "The arrival time for transit mode directions request.")
	waypoints                = flag.String("waypoints", "", "The waypoints for driving directions request, | separated.")
	alternatives             = flag.Bool("alternatives", false, "Whether the Directions service may provide more than one route alternative in the response.")
	avoid                    = flag.String("avoid", "", "Indicates that the calculated route(s) should avoid the indicated features, | separated.")
	language                 = flag.String("language", "", "Specifies the language in which to return results.")
	units                    = flag.String("units", "", "Specifies the unit system to use when returning results.")
	region                   = flag.String("region", "", "Specifies the region code, specified as a ccTLD (\"top-level domain\") two-character value.")
	transitMode              = flag.String("transit_mode", "", "Specifies one or more preferred modes of transit, | separated. This parameter may only be specified for transit directions.")
	transitRoutingPreference = flag.String("transit_routing_preference", "", "Specifies preferences for transit routes.")
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

	r := &maps.DirectionsRequest{
		Origin:        *origin,
		Destination:   *destination,
		DepartureTime: *departureTime,
		ArrivalTime:   *arrivalTime,
		Alternatives:  *alternatives,
		Language:      *language,
		Region:        *region,
	}

	lookupMode(*mode, r)
	lookupUnits(*units, r)
	lookupTransitRoutingPreference(*transitRoutingPreference, r)

	if *waypoints != "" {
		r.Waypoints = strings.Split(*waypoints, "|")
	}

	if *avoid != "" {
		for _, a := range strings.Split(*avoid, "|") {
			switch {
			case a == "tolls":
				r.Avoid = append(r.Avoid, maps.AvoidTolls)
			case a == "highways":
				r.Avoid = append(r.Avoid, maps.AvoidHighways)
			case a == "ferries":
				r.Avoid = append(r.Avoid, maps.AvoidFerries)
			}
		}

	}
	if *transitMode != "" {
		for _, t := range strings.Split(*transitMode, "|") {
			switch {
			case t == "bus":
				r.TransitMode = append(r.TransitMode, maps.TransitModeBus)
			case t == "subway":
				r.TransitMode = append(r.TransitMode, maps.TransitModeSubway)
			case t == "train":
				r.TransitMode = append(r.TransitMode, maps.TransitModeTrain)
			case t == "tram":
				r.TransitMode = append(r.TransitMode, maps.TransitModeTram)
			case t == "rail":
				r.TransitMode = append(r.TransitMode, maps.TransitModeRail)
			}
		}
	}

	pretty.Println(r)

	ctx := maps.NewContext(*apiKey, client)
	resp, err := r.Get(ctx)
	if err != nil {
		log.Fatalf("fatal error %v", err)
	}

	pretty.Println(resp)
}

func lookupMode(mode string, r *maps.DirectionsRequest) {
	switch {
	case mode == "driving":
		r.Mode = maps.TravelModeDriving
	case mode == "walking":
		r.Mode = maps.TravelModeWalking
	case mode == "bicycling":
		r.Mode = maps.TravelModeBicycling
	case mode == "transit":
		r.Mode = maps.TravelModeTransit
	}
}

func lookupUnits(units string, r *maps.DirectionsRequest) {
	switch {
	case units == "metric":
		r.Units = maps.UnitsMetric
	case units == "imperial":
		r.Units = maps.UnitsImperial
	}
}

func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DirectionsRequest) {
	switch {
	case transitRoutingPreference == "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case transitRoutingPreference == "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	}
}
