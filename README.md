Go Client for Google Maps Services
==================================

[![Build Status](https://travis-ci.org/googlemaps/google-maps-services-go.svg?branch=master)](https://travis-ci.org/googlemaps/google-maps-services-go)
[![GoDoc](https://godoc.org/googlemaps.github.io/maps?status.svg)](https://godoc.org/googlemaps.github.io/maps)

## Description
Use Go? Want to [geocode][Geocoding API] something? Looking for [directions][Directions API]?
Maybe [matrices of directions][Distance Matrix API]? This library brings the [Google Maps API Web
Services] to your Go application.
![Analytics](https://maps-ga-beacon.appspot.com/UA-12846745-20/google-maps-services-go/readme?pixel)

The Go Client for Google Maps Services is a Go Client library for the following Google Maps
APIs:

 - [Directions API]
 - [Distance Matrix API]
 - [Elevation API]
 - [Geocoding API]
 - [Places API]
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
 - Go 1.5 or later.
 - A Google Maps API key.

### API keys

Each Google Maps Web Service request requires an API key or client ID. API keys
are freely available with a Google Account at
https://developers.google.com/console. The type of API key you need is a 
**Server key**. 

To get an API key:

 1. Visit https://developers.google.com/console and log in with
    a Google Account.
 1. Select one of your existing projects, or create a new project.
 1. Enable the API(s) you want to use. The Go Client for Google Maps Services
    accesses the following APIs:
    * Directions API
    * Distance Matrix API
    * Elevation API
    * Geocoding API
    * Places API
    * Roads API
    * Time Zone API
 1. Create a new **Server key**.
 1. If you'd like to restrict requests to a specific IP address, do so now.
 
For guided help, follow the instructions for the [Directions API][directions-key]. 
You only need one API key, but remember to enable all the APIs you need.
For even more information, see the guide to [API keys][apikey]. 

**Important:** This key should be kept secret on your server.

## Installation

    $ go get googlemaps.github.io/maps

## Developer Documentation

View the [reference documentation](https://godoc.org/googlemaps.github.io/maps)

Additional documentation for the included  web services is available at
https://developers.google.com/maps/ and https://developers.google.com/places/.

 - [Directions API]
 - [Distance Matrix API]
 - [Elevation API]
 - [Geocoding API]
 - [Places API]
 - [Time Zone API]
 - [Roads API]

## Usage
Sample usage of the Directions API with an API key:

```go
package main

import (
	"log"

	"googlemaps.github.io/maps"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
)

func main() {
	c, err := maps.NewClient(maps.WithAPIKey("Insert-API-Key-Here"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}
	resp, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(resp)
}
```

Below is the same example, using client ID and client secret (digital signature)
for authentication. This code assumes you have previously loaded the `clientID`
and `clientSecret` variables with appropriate values.

For a guide on how to generate the `clientSecret` (digital signature), see the
documentation for the API you're using. For example, see the guide for the
[Directions API][directions-client-id].

```go
package main

import (
	"log"

        "googlemaps.github.io/maps"
        "github.com/kr/pretty"
        "golang.org/x/net/context"
)

func main() {
        c, err := maps.NewClient(maps.WithClientIDAndSignature(clientID, clientSecret))
        if err != nil {
            log.Fatalf("fatal error: %s", err)
        }
        r := &maps.DirectionsRequest{
            Origin:      "Sydney",
            Destination: "Perth",
        }
        resp, err := c.Directions(context.Background(), r)
        if err != nil {
            log.Fatalf("fatal error: %s", err)
        }

        pretty.Println(resp)
}
```

## Features

### Rate limiting

Never sleep between requests again! By default, requests are sent at the expected rate limits for
each web service, typically 10 queries per second for free users. If you want to speed up or slow
down requests, you can do that too, using `maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(qps))`.

### Client IDs

Google Maps APIs Premium Plan customers can use their [client ID and secret][clientid] to authenticate,
instead of an API key.

### Native types

Native objects for each of the API responses.

[apikey]: https://developers.google.com/maps/faq#keysystem
[clientid]: https://developers.google.com/maps/documentation/business/webservices/auth

[Google Maps API Web Services]: https://developers.google.com/maps/documentation/webservices/
[Directions API]: https://developers.google.com/maps/documentation/directions/
[directions-client-id]: https://developers.google.com/maps/documentation/directions/get-api-key#client-id
[directions-key]: https://developers.google.com/maps/documentation/directions/get-api-key#key
[Distance Matrix API]: https://developers.google.com/maps/documentation/distancematrix/
[Elevation API]: https://developers.google.com/maps/documentation/elevation/
[Geocoding API]: https://developers.google.com/maps/documentation/geocoding/
[Places API]: https://developers.google.com/places/web-service/
[Roads API]: https://developers.google.com/maps/documentation/roads/
[Time Zone API]: https://developers.google.com/maps/documentation/timezone/

[issues]: https://github.com/googlemaps/google-maps-services-go/issues
[contrib]: https://github.com/googlemaps/google-maps-services-go/blob/master/CONTRIB.md
