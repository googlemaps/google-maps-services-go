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

// Package main contains a simple command line tool for Places Photos API
// Documentation: https://developers.google.com/places/web-service/photos
package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var (
	apiKey         = flag.String("key", "", "API Key for using Google Maps API.")
	clientID       = flag.String("client_id", "", "ClientID for Maps for Work API access.")
	signature      = flag.String("signature", "", "Signature for Maps for Work API access.")
	photoreference = flag.String("photoreference", "", "Textual identifier that uniquely identifies a place photo.")
	maxheight      = flag.Int("maxheight", 0, "Specifies the maximum desired height, in pixels, of the image returned by the Place Photos service. One of maxheight and maxwidth is required.")
	maxwidth       = flag.Int("maxwidth", 0, "Specifies the maximum desired width, in pixels, of the image returned by the Place Photos service. One of maxheight and maxwidth is required.")
	basename       = flag.String("basename", "", "Base name of file to write image to. If not specified, no file will be written.")
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

	r := &maps.PlacePhotoRequest{
		PhotoReference: *photoreference,
		MaxHeight:      uint(*maxheight),
		MaxWidth:       uint(*maxwidth),
	}

	resp, err := client.PlacePhoto(context.Background(), r)
	check(err)

	log.Printf("Content-Type: %v\n", resp.ContentType)
	img, err := resp.Image()
	check(err)
	log.Printf("Image bounds: %v", img.Bounds())

	if *basename != "" {
		filename := fmt.Sprintf("%s.%s", *basename, "jpg")
		f, err := os.Create(filename)
		check(err)
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
		check(err)

		log.Printf("Wrote image to %s\n", filename)
	}
}
