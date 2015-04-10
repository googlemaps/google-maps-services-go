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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	location  = flag.String("location", "", "a comma-separated lat,lng tuple (eg. location=-33.86,151.20), representing the location to look up.")
	timestamp = flag.String("timestamp", "", "specifies the desired time as seconds since midnight, January 1, 1970 UTC.")
	language  = flag.String("language", "", "The language in which to return results.")
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
	t, err := strconv.Atoi(*timestamp)
	if err != nil {
		usageAndExit(fmt.Sprintf("Could not convert timestamp to int: %v", err))
	}

	r := &maps.TimezoneRequest{
		Language:  *language,
		Timestamp: &t,
	}

	parseLocation(*location, r)

	pretty.Println(r)

	resp, err := r.Get(ctx)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	pretty.Println(resp)
}

func parseLocation(location string, r *maps.TimezoneRequest) {
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
	} else {
		usageAndExit("location is required")
	}
}
