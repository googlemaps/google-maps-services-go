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
- [Maps Static API]

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

- Go 1.7 or later.
- A Google Maps API key.

### API keys

Each Google Maps Web Service request requires an API key or client ID. API keys
are freely available with a Google Account at
[Google APIs Console][API Console]. The type of API key you need is a **Server key**.

To get an API key:

 1. Visit [Google APIs Console][API Console] and log in with
    a Google Account.
 1. Select one of your existing projects, or create a new project.
 1. Enable the API(s) you want to use. The Go Client for Google Maps Services
    accesses the following APIs:
    - Directions API
    - Distance Matrix API
    - Elevation API
    - Geocoding API
    - Places API
    - Roads API
    - Time Zone API
    - Maps Static API
 1. Create a new **Server key**.
 1. If you'd like to restrict requests to a specific IP address, do so now.

For guided help, follow the instructions for the [Directions API][directions-key].
You only need one API key, but remember to enable all the APIs you need.
For even more information, see the guide to [API keys][apikey].

**Important:** This key should be kept secret on your server.

## Installation

To install the Go Client for Google Maps Services, please execute the following `go get` command.

```bash
    go get googlemaps.github.io/maps
```

## Developer Documentation

View the [reference documentation](https://godoc.org/googlemaps.github.io/maps)

Additional documentation for the included  web services is available at
[developers.google.com/maps][Maps documentation] and
[developers.google.com/places][Places documentation].

- [Directions API]
- [Distance Matrix API]
- [Elevation API]
- [Geocoding API]
- [Places API]
- [Time Zone API]
- [Roads API]
- [Maps Static API]

## Usage

Sample usage of the Directions API with an API key:

```go
package main

import (
	"context"
	"log"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
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
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(route)
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
	"context"
	"log"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

func main() {
	c, err := maps.NewClient(maps.WithClientIDAndSignature("Client ID", "Client Secret"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(route)
}
```

## Features

### Rate limiting

Never sleep between requests again! By default, requests are sent at the expected rate limits for
each web service, typically 50 queries per second for free users. If you want to speed up or slow
down requests, you can do that too, using `maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(qps))`.

### Client IDs

Google Maps APIs Premium Plan customers can use their [client ID and secret][clientid] to authenticate,
instead of an API key.

### Native types

Native objects for each of the API responses.

### Monitoring

It's possible to get metrics for status counts and latency histograms for monitoring.
Use `maps.WithMetricReporter(metrics.OpenCensusReporter{})` to log metrics to OpenCensus,
and `metrics.RegisterViews()` to make the metrics available to be exported.
OpenCensus can export these metrics to a [variety of monitoring services](https://opencensus.io/exporters/).
You can also implement your own metric reporter instead of using the provided one.

[apikey]: https://developers.google.com/maps/faq#keysystem
[clientid]: https://developers.google.com/maps/documentation/business/webservices/auth

[API Console]: https://developers.google.com/console
[Maps documentation]: https://developers.google.com/maps/
[Places documentation]: https://developers.google.com/places/

[Google Maps API Web Services]: https://developers.google.com/maps/apis-by-platform#web_service_apis
[Directions API]: https://developers.google.com/maps/documentation/directions/
[directions-client-id]: https://developers.google.com/maps/documentation/directions/get-api-key#client-id
[directions-key]: https://developers.google.com/maps/documentation/directions/get-api-key#key
[Distance Matrix API]: https://developers.google.com/maps/documentation/distancematrix/
[Elevation API]: https://developers.google.com/maps/documentation/elevation/
[Geocoding API]: https://developers.google.com/maps/documentation/geocoding/
[Places API]: https://developers.google.com/places/web-service/
[Roads API]: https://developers.google.com/maps/documentation/roads/
[Time Zone API]: https://developers.google.com/maps/documentation/timezone/
[Maps Static API]: https://developers.google.com/maps/documentation/maps-static/

[issues]: https://github.com/googlemaps/google-maps-services-go/issues
[contrib]: https://github.com/googlemaps/google-maps-services-go/blob/master/CONTRIB.md
