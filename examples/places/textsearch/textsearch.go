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
	"strconv"
	"strings"

	"developers.google.com/maps"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	clientID  = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature = flag.String("signature", "", "Signature for Maps for Work API access.")
	query     = flag.String("query", "", "Text Search query to execute.")
	language  = flag.String("language", "", "The language in which to return results.")
	location  = flag.String("location", "", "The latitude/longitude around which to retrieve place information. This must be specified as latitude,longitude.")
	radius    = flag.Uint("radius", 0, "Defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters.")
	minprice  = flag.String("min_price", "", "Restricts results to only those places within the specified price level.")
	maxprice  = flag.String("max_price", "", "Restricts results to only those places within the specified price level.")
	opennow   = flag.Bool("open_now", false, "Restricts results to only those places that are open for business at the time the query is sent.")
)

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
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
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r := &maps.TextSearchRequest{
		Query:    *query,
		Language: *language,
		Radius:   *radius,
		OpenNow:  *opennow,
	}

	parseLocation(*location, r)
	parsePriceLevels(*minprice, *maxprice, r)

	resp, err := client.TextSearch(context.Background(), r)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	pretty.Println(resp)
}

func parseLocation(location string, r *maps.TextSearchRequest) {
	if location != "" {
		l := strings.Split(location, ",")
		lat, err := strconv.ParseFloat(l[0], 64)
		if err != nil {
			log.Fatalf("Couldn't parse latlng: %#v", err)
		}
		lng, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			log.Fatalf("Couldn't parse latlng: %#v", err)
		}
		r.Location = &maps.LatLng{
			Lat: lat,
			Lng: lng,
		}
	}
}

func parsePriceLevels(minprice string, maxprice string, r *maps.TextSearchRequest) {
	if minprice != "" {
		switch minprice {
		case "0":
			r.MinPrice = maps.PriceLevelFree
		case "1":
			r.MinPrice = maps.PriceLevelInexpensive
		case "2":
			r.MinPrice = maps.PriceLevelModerate
		case "3":
			r.MinPrice = maps.PriceLevelExpensive
		case "4":
			r.MinPrice = maps.PriceLevelVeryExpensive
		default:
			usageAndExit(fmt.Sprintf("Unknown min_price level: '%s'", minprice))
		}
	}

	if maxprice != "" {
		switch maxprice {
		case "0":
			r.MaxPrice = maps.PriceLevelFree
		case "1":
			r.MaxPrice = maps.PriceLevelInexpensive
		case "2":
			r.MaxPrice = maps.PriceLevelModerate
		case "3":
			r.MaxPrice = maps.PriceLevelExpensive
		case "4":
			r.MaxPrice = maps.PriceLevelVeryExpensive
		default:
			usageAndExit(fmt.Sprintf("Unknown max_price level: '%s'", maxprice))
		}
	}
}
