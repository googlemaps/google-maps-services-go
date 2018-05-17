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
	apiKey    = flag.String("key", "", "API Key for using Google Maps API.")
	input     = flag.String("input", "", "The text input specifying which place to search for (for example, a name, address, or phone number).")
	inputType = flag.String("inputtype", "", "The type of input. This can be one of either textquery or phonenumber.")
	fields    = flag.String("fields", "", "Comma seperated list of Fields")
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
	} else {
		usageAndExit("Please specify an API Key.")
	}
	check(err)

	r := &maps.FindPlaceFromTextRequest{}

	r.Input = *input

	parseInputType(*inputType, r)

	if *fields != "" {
		f, err := parseFields(*fields)
		check(err)
		r.Fields = f
	}

	resp, err := client.FindPlaceFromText(context.Background(), r)
	check(err)

	pretty.Println(resp)
}

func parseFields(fields string) ([]maps.PlaceDetailsFieldMask, error) {
	var res []maps.PlaceDetailsFieldMask
	for _, s := range strings.Split(fields, ",") {
		f, err := maps.ParsePlaceDetailsFieldMask(s)
		if err != nil {
			return nil, err
		}
		res = append(res, f)
	}
	return res, nil
}

func parseInputType(inputType string, r *maps.FindPlaceFromTextRequest) {
	switch inputType {
	case "textquery":
		r.InputType = maps.FindPlaceFromTextInputTypeTextQuery
	case "phonenumber":
		r.InputType = maps.FindPlaceFromTextInputTypePhoneNumber
	default:
		usageAndExit(fmt.Sprintf("Unknown input type: '%s'", inputType))
	}
}
