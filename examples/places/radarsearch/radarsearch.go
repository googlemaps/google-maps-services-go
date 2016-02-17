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

// Package main contains a simple command line tool for Places API Text Search
// Documentation: https://developers.google.com/places/web-service/search#TextSearchRequests
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	clientID  = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature = flag.String("signature", "", "Signature for Maps for Work API access.")
	location  = flag.String("location", "", "The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude.")
	radius    = flag.Uint("radius", 0, "Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters.")
	keyword   = flag.String("keyword", "", "A term to be matched against all content that Google has indexed for this place, including but not limited to name, type, and address, as well as customer reviews and other third-party content.")
	minprice  = flag.String("min_price", "", "Restricts results to only those places within the specified price level.")
	maxprice  = flag.String("max_price", "", "Restricts results to only those places within the specified price level.")
	name      = flag.String("name", "", "One or more terms to be matched against the names of places, separated by a space character. Results will be restricted to those containing the passed name values.")
	opennow   = flag.Bool("open_now", false, "Restricts results to only those places that are open for business at the time the query is sent.")
	placeType = flag.String("type", "", "Restricts the results to places matching the specified type.")
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

	r := &maps.RadarSearchRequest{
		Radius:  *radius,
		Keyword: *keyword,
		Name:    *name,
		OpenNow: *opennow,
	}

	parseLocation(*location, r)
	parsePriceLevels(*minprice, *maxprice, r)
	parsePlaceType(*placeType, r)

	resp, err := client.RadarSearch(context.Background(), r)
	check(err)

	for i, result := range resp.Results {
		r2 := &maps.PlaceDetailsRequest{
			PlaceID: result.PlaceID,
		}
		resp, err := client.PlaceDetails(context.Background(), r2)
		check(err)

		fmt.Printf("%d: %v at %v\n", i, resp.Name, resp.FormattedAddress)
	}
}

func parseLocation(location string, r *maps.RadarSearchRequest) {
	if location != "" {
		l, err := maps.ParseLatLng(location)
		check(err)
		r.Location = &l
	}
}

func parsePriceLevel(priceLevel string) maps.PriceLevel {
	switch priceLevel {
	case "0":
		return maps.PriceLevelFree
	case "1":
		return maps.PriceLevelInexpensive
	case "2":
		return maps.PriceLevelModerate
	case "3":
		return maps.PriceLevelExpensive
	case "4":
		return maps.PriceLevelVeryExpensive
	default:
		usageAndExit(fmt.Sprintf("Unknown price level: '%s'", priceLevel))
	}
	return maps.PriceLevelFree
}

func parsePriceLevels(minprice string, maxprice string, r *maps.RadarSearchRequest) {
	if minprice != "" {
		r.MinPrice = parsePriceLevel(minprice)
	}

	if maxprice != "" {
		r.MaxPrice = parsePriceLevel(minprice)
	}
}

func parsePlaceType(placeType string, r *maps.RadarSearchRequest) {
	if placeType != "" {
		t, err := maps.ParsePlaceType(placeType)
		check(err)
		r.Type = t
	}
}
