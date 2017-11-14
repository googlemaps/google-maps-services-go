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

// Package main contains a simple command line tool for Places API Query Autocomplete
// Documentation: https://developers.google.com/places/web-service/query
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
	apiKey       = flag.String("key", "", "API Key for using Google Maps API.")
	clientID     = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature    = flag.String("signature", "", "Signature for Maps for Work API access.")
	input        = flag.String("input", "", "Text string on which to search.")
	language     = flag.String("language", "", "The language in which to return results.")
	offset       = flag.Uint("offset", 0, "The character position in the input term at which the service uses text for predictions.")
	location     = flag.String("location", "", "The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude.")
	radius       = flag.Uint("radius", 0, "Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters.")
	placeType    = flag.String("types", "", "Restricts the results to places matching the specified type.")
	components   = flag.String("components", "", "A component filter for specifying which country to perform autocomplete inside of")
	strictbounds = flag.Bool("strictbounds", false, "Whether to strictly enforce bounds.")
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

	r := &maps.PlaceAutocompleteRequest{
		Input:        *input,
		Language:     *language,
		Offset:       *offset,
		Radius:       *radius,
		StrictBounds: *strictbounds,
	}

	parseLocation(*location, r)
	parsePlaceType(*placeType, r)

	resp, err := client.PlaceAutocomplete(context.Background(), r)
	check(err)

	pretty.Println(resp)
}

func parseLocation(location string, r *maps.PlaceAutocompleteRequest) {
	if location != "" {
		l, err := maps.ParseLatLng(location)
		check(err)
		r.Location = &l
	}
}

func parsePlaceType(placeType string, r *maps.PlaceAutocompleteRequest) {
	if placeType != "" {
		t, err := maps.ParseAutocompletePlaceType(placeType)
		check(err)
		r.Types = t
	}
}

func parseComponents(components string, r *maps.PlaceAutocompleteRequest) {
	if components != "" {
		c := strings.Split(components, "|")
		for _, cf := range c {
			i := strings.Split(cf, ":")
			switch i[0] {
			case "country":
				r.Components[maps.ComponentCountry] = i[1]
			default:
				log.Fatalf("Unsupported component \"%v\"", i[0])
			}
		}
	}
}
