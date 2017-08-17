// Copyright 2017 Google Inc. All Rights Reserved.
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

// Package main contains a simple command line tool for Static Maps API
// Documentation: https://developers.google.com/maps/documentation/static-maps/
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
	center    = flag.String("center", "", "Center the center of the map, equidistant from all edges of the map.")
	zoom      = flag.Int("zoom", -1, "Zoom the zoom level of the map, which determines the magnification level of the map.")
	size      = flag.String("size", "", "Size defines the rectangular dimensions of the map image.")
	scale     = flag.Int("scale", -1, "Scale affects the number of pixels that are returned.")
	format    = flag.String("format", "", "Format defines the format of the resulting image.")
	maptype   = flag.String("maptype", "", "Maptype defines the type of map to construct.")
	language  = flag.String("langauge", "", "Language defines the language to use for display of labels on map tiles.")
	region    = flag.String("region", "", "Region the appropriate borders to display, based on geo-political sensitivities.")
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

	r := &maps.StaticMapRequest{
		Center:   *center,
		Zoom:     *zoom,
		Size:     *size,
		Scale:    *scale,
		Format:   maps.Format(*format),
		Language: *language,
		Region:   *region,
		MapType:  maps.MapType(*maptype),
	}

	resp, err := client.StaticMap(context.Background(), r)
	check(err)

	pretty.Println(resp)
}
