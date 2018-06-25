// Copyright 2018 Google Inc. All Rights Reserved.
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

// Package main contains a simple command line tool for Find Place From Text API
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
	input        = flag.String("input", "", "The text input specifying which place to search for (for example, a name, address, or phone number).")
	inputType    = flag.String("inputtype", "", "The type of input. This can be one of either textquery or phonenumber.")
	fields       = flag.String("fields", "", "Comma seperated list of Fields")
	locationbias = flag.String("locationbias", "", "Location bias for this request. Optional. One of ipbias, point, circle, or rectangle.")
	point        = flag.String("point", "", "The latitude/longitude for location bias point. This must be specified as latitude,longitude.")
	center       = flag.String("center", "", "The center latitude/longitude for location bias circle. This must be specified as latitude,longitude.")
	radius       = flag.Int("radius", 0, "The radius for location bias circle.")
	southwest    = flag.String("southwest", "", "The South West latitude/longitude for location bias rectangle. This must be specified as latitude,longitude.")
	northeast    = flag.String("northeast", "", "The North East latitude/longitude for location bias rectangle. This must be specified as latitude,longitude.")
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
	if *apiKey == "" {
		usageAndExit("Please specify an API Key.")
	}
	client, err = maps.NewClient(maps.WithAPIKey(*apiKey))
	check(err)

	r := &maps.FindPlaceFromTextRequest{
		Input:     *input,
		InputType: parseInputType(*inputType),
	}

	if *locationbias != "" {
		lb, err := maps.ParseFindPlaceFromTextLocationBiasType(*locationbias)
		check(err)
		r.LocationBias = lb
		switch lb {
		case maps.FindPlaceFromTextLocationBiasPoint:
			l, err := maps.ParseLatLng(*point)
			check(err)
			r.LocationBiasPoint = &l
		case maps.FindPlaceFromTextLocationBiasCircular:
			l, err := maps.ParseLatLng(*center)
			check(err)
			r.LocationBiasCenter = &l
			r.LocationBiasRadius = *radius
		case maps.FindPlaceFromTextLocationBiasRectangular:
			sw, err := maps.ParseLatLng(*southwest)
			check(err)
			r.LocationBiasSouthWest = &sw
			ne, err := maps.ParseLatLng(*northeast)
			check(err)
			r.LocationBiasNorthEast = &ne
		}
	}

	if *fields != "" {
		f, err := parseFields(*fields)
		check(err)
		r.Fields = f
	}

	resp, err := client.FindPlaceFromText(context.Background(), r)
	check(err)

	pretty.Println(resp)
}

func parseFields(fields string) ([]maps.PlaceSearchFieldMask, error) {
	var res []maps.PlaceSearchFieldMask
	for _, s := range strings.Split(fields, ",") {
		f, err := maps.ParsePlaceSearchFieldMask(s)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func parseInputType(inputType string) maps.FindPlaceFromTextInputType {
	var it maps.FindPlaceFromTextInputType
	switch inputType {
	case "textquery":
		it = maps.FindPlaceFromTextInputTypeTextQuery
	case "phonenumber":
		it = maps.FindPlaceFromTextInputTypePhoneNumber
	default:
		usageAndExit(fmt.Sprintf("Unknown input type: '%s'", inputType))
	}
	return it
}
