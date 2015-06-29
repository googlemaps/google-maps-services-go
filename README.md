Go Client for Google Maps Services
==================================

## Description
Use Go? Want to [geocode][Geocoding API] something? Looking for [directions][Directions API]?
Maybe [matrices of directions][Distance Matrix API]? This library brings the [Google Maps API Web
Services] to your Go application.
![Analytics](https://ga-beacon.appspot.com/UA-12846745-20/google-maps-services-go/readme?pixel)

The Go Client for Google Maps Services is a Go Client library for the following Google Maps
APIs:

 - [Directions API]
 - [Distance Matrix API]
 - [Elevation API]
 - [Geocoding API]
 - [Roads API]
 - [Time Zone API]

Keep in mind that the same [terms and conditions](https://developers.google.com/maps/terms) apply
to usage of the APIs when they're accessed through this library.

## Support

This library is community supported. We're comfortable enough with the stability and features of
the library that we want you to build real production applications on it. We will try to support,
through Stack Overflow, the public and protected surface of the library and maintain backwards
compatibility in the future; however, while the library is in version 0.x, we reserve the right
to make backwards-incompatible changes. If we do remove some functionality (typically because
better functionality exists or if the feature proved infeasible), our intention is to deprecate
and give developers a year to update their code.

If you find a bug, or have a feature suggestion, please [log an issue][issues]. If you'd like to
contribute, please read [How to Contribute][contrib].

## Requirements
 - Go 1.x or later.
 - A Google Maps API key.

### API keys

Each Google Maps Web Service requires an API key or Client ID. API keys are
freely available with a Google Account at https://developers.google.com/console.
To generate a server key for your project:

 1. Visit https://developers.google.com/console and log in with
    a Google Account.
 1. Select an existing project, or create a new project.
 1. Click **Enable an API**.
 1. Browse for the API, and set its status to "On". The Golang Client for Google Maps Services
    accesses the following APIs:
    * Directions API
    * Distance Matrix API
    * Elevation API
    * Geocoding API
    * Roads API
    * Time Zone API
 1. Once you've enabled the APIs, click **Credentials** from the left navigation of the Developer
    Console.
 1. In the "Public API access", click **Create new Key**.
 1. Choose **Server Key**.
 1. If you'd like to restrict requests to a specific IP address, do so now.
 1. Click **Create**.

Your API key should be 40 characters long, and begin with `AIza`.

**Important:** This key should be kept secret on your server.

## Installation

    $ go get github.com/googlemaps/google-maps-services-go

## Developer Documentation

View the [reference documentation](http://googlemaps.github.io/google-maps-services-go/docs)
TODO(brettmorgan): make this documentation.

Additional documentation for the included web services is available at
https://developers.google.com/maps/.

 - [Directions API]
 - [Distance Matrix API]
 - [Elevation API]
 - [Geocoding API]
 - [Time Zone API]
 - [Roads API]

## Usage
Sample usage of the Directions API:

```go
package main

import (
	"log"
	"net/http"

	"github.com/kr/pretty"

	"google.golang.org/maps"
)

func main() {

  apiKey := "Insert-API-Key-Here"
	client := maps.NewClient(&http.Client{}, apiKey)

	r := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}

	resp, err := client.GetDirections(r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(resp)
}
```

## Features
TODO(brettmorgan): Implement features

[apikey]: https://developers.google.com/maps/faq#keysystem
[clientid]: https://developers.google.com/maps/documentation/business/webservices/auth

[Google Maps API Web Services]: https://developers.google.com/maps/documentation/webservices/
[Directions API]: https://developers.google.com/maps/documentation/directions/
[Distance Matrix API]: https://developers.google.com/maps/documentation/distancematrix/
[Elevation API]: https://developers.google.com/maps/documentation/elevation/
[Geocoding API]: https://developers.google.com/maps/documentation/geocoding/
[Roads API]: https://developers.google.com/maps/documentation/roads/
[Time Zone API]: https://developers.google.com/maps/documentation/timezone/

[issues]: https://github.com/googlemaps/google-maps-services-python/issues
[contrib]: https://github.com/googlemaps/google-maps-services-python/blob/master/CONTRIB.md
