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

// Package main contains a simple command line tool for Directions API
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	apiKey                   = flag.String("key", "", "API Key for using Google Maps API.")
	clientID                 = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature                = flag.String("signature", "", "Signature for Maps for Work API access.")
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
	iterations               = flag.Int("iterations", 1, "Number of times to make API request.")
	trafficModel             = flag.String("traffic_model", "", "Specifies traffic prediction model when request future directions. Valid values are optimistic, best_guess, and pessimistic. Optional.")
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
		client, err = maps.NewClient(maps.WithAPIKey(*apiKey), maps.WithRateLimit(2))
	} else if *clientID != "" || *signature != "" {
		client, err = maps.NewClient(maps.WithClientIDAndSignature(*clientID, *signature))
	} else {
		usageAndExit("Please specify an API Key, or Client ID and Signature.")
	}
	check(err)

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
	lookupTrafficModel(*trafficModel, r)

	if *waypoints != "" {
		r.Waypoints = strings.Split(*waypoints, "|")
	}

	if *avoid != "" {
		for _, a := range strings.Split(*avoid, "|") {
			switch a {
			case "tolls":
				r.Avoid = append(r.Avoid, maps.AvoidTolls)
			case "highways":
				r.Avoid = append(r.Avoid, maps.AvoidHighways)
			case "ferries":
				r.Avoid = append(r.Avoid, maps.AvoidFerries)
			default:
				log.Fatalf("Unknown avoid restriction %s", a)
			}
		}
	}
	if *transitMode != "" {
		for _, t := range strings.Split(*transitMode, "|") {
			switch t {
			case "bus":
				r.TransitMode = append(r.TransitMode, maps.TransitModeBus)
			case "subway":
				r.TransitMode = append(r.TransitMode, maps.TransitModeSubway)
			case "train":
				r.TransitMode = append(r.TransitMode, maps.TransitModeTrain)
			case "tram":
				r.TransitMode = append(r.TransitMode, maps.TransitModeTram)
			case "rail":
				r.TransitMode = append(r.TransitMode, maps.TransitModeRail)
			}
		}
	}

	if *iterations == 1 {
		routes, waypoints, err := client.Directions(context.Background(), r)
		check(err)

		pretty.Println(waypoints)
		pretty.Println(routes)
	} else {
		done := make(chan iterationResult)
		for i := 0; i < *iterations; i++ {
			go func(i int) {
				startTime := time.Now()
				_, _, err := client.Directions(context.Background(), r)
				done <- iterationResult{
					fmt.Sprintf("Iteration %2d: round trip %.2f seconds", i, float64(time.Now().Sub(startTime))/1000000000),
					err,
				}
			}(i)
		}

		for i := 0; i < *iterations; i++ {
			result := <-done
			if err != nil {
				fmt.Printf("error: %+v\n", result.err)
			} else {
				fmt.Println(result.result)
			}
		}
	}
}

type iterationResult struct {
	result string
	err    error
}

func lookupMode(mode string, r *maps.DirectionsRequest) {
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
		log.Fatalf("Unknown mode '%s'", mode)
	}
}

func lookupUnits(units string, r *maps.DirectionsRequest) {
	switch units {
	case "metric":
		r.Units = maps.UnitsMetric
	case "imperial":
		r.Units = maps.UnitsImperial
	case "":
		// ignore
	default:
		log.Fatalf("Unknown units '%s'", units)
	}
}

func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DirectionsRequest) {
	switch transitRoutingPreference {
	case "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	case "":
		// ignore
	default:
		log.Fatalf("Unknown transit routing preference %s", transitRoutingPreference)
	}
}

func lookupTrafficModel(trafficModel string, r *maps.DirectionsRequest) {
	switch trafficModel {
	case "optimistic":
		r.TrafficModel = maps.TrafficModelOptimistic
	case "best_guess":
		r.TrafficModel = maps.TrafficModelBestGuess
	case "pessimistic":
		r.TrafficModel = maps.TrafficModelPessimistic
	case "":
		// ignore
	default:
		log.Fatalf("Unknown traffic mode %s", trafficModel)
	}
}
