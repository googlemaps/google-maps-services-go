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

	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	clientID  = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature = flag.String("signature", "", "Signature for Maps for Work API access.")
	location  = flag.String("location", "", "The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude.")
	radius    = flag.Uint("radius", 0, "Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters.")
	keyword   = flag.String("keyword", "", "Specifies the language in which to return results. Optional.")
	language  = flag.String("language", "", "The language in which to return results.")
	minPrice  = flag.String("minprice", "", "Restricts results to only those places within the specified price level.")
	maxPrice  = flag.String("maxprice", "", "Restricts results to only those places within the specified price level.")
	name      = flag.String("name", "", "One or more terms to be matched against the names of places, separated with a space character.")
	openNow   = flag.Bool("open_now", false, "Restricts results to only those places that are open for business at the time the query is sent.")
	rankBy    = flag.String("rankby", "", "Specifies the order in which results are listed. Valid values are prominence or distance.")
	placeType = flag.String("type", "", "Restricts the results to places matching the specified type.")
	pageToken = flag.String("pagetoken", "", "Set to retrieve the next next page of results.")
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

	r := &maps.NearbySearchRequest{
		Radius:    *radius,
		Keyword:   *keyword,
		Language:  *language,
		Name:      *name,
		OpenNow:   *openNow,
		PageToken: *pageToken,
	}

	parseLocation(*location, r)
	parsePriceLevels(*minPrice, *maxPrice, r)
	parseRankBy(*rankBy, r)
	parsePlaceType(*placeType, r)

	resp, err := client.NearbySearch(context.Background(), r)
	check(err)

	pretty.Println(resp)
}

func parseLocation(location string, r *maps.NearbySearchRequest) {
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

func parsePriceLevels(minPrice string, maxPrice string, r *maps.NearbySearchRequest) {
	if minPrice != "" {
		r.MinPrice = parsePriceLevel(minPrice)
	}

	if maxPrice != "" {
		r.MaxPrice = parsePriceLevel(minPrice)
	}
}

func parseRankBy(rankBy string, r *maps.NearbySearchRequest) {
	switch rankBy {
	case "prominence":
		r.RankBy = maps.RankByProminence
		return
	case "distance":
		r.RankBy = maps.RankByDistance
		return
	case "":
		return
	default:
		usageAndExit(fmt.Sprintf("Unknown rank by: \"%v\"", rankBy))
	}
}

func parsePlaceType(placeType string, r *maps.NearbySearchRequest) {
	if placeType != "" {
		t, err := maps.ParsePlaceType(placeType)
		if err != nil {
			usageAndExit(fmt.Sprintf("Unknown place type \"%v\"", placeType))
		}

		r.Type = t
	}
}
