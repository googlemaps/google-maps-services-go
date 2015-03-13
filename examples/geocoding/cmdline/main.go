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

// Package main contains a simple command line tool for Geocoding API
// Directions docs: https://developers.google.com/maps/documentation/geocoding/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kr/pretty"
	"google.golang.org/maps"
)

var (
	apiKey  = flag.String("key", "", "API Key for using Google Maps API.")
	address = flag.String("address", "", "The street address that you want to geocode, in the format used by the national postal service of the country concerned.")
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
	r := &maps.GeocodingRequest{
		Address: *address,
	}

	resp, err := r.Get(ctx)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	pretty.Println(resp)
}
