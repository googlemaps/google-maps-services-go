Go Client for Google Maps Services
==================================

[![GoDoc](https://godoc.org/googlemaps.github.io/maps?status.svg)](https://godoc.org/googlemaps.github.io/maps)
[![Go Report Card](https://goreportcard.com/badge/github.com/googlemaps/google-maps-services-go)](https://goreportcard.com/report/github.com/googlemaps/google-maps-services-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/googlemaps/google-maps-services-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Description

Use Go? This library brings many [Google Maps Platform Web Services APIs] to your Go application.

The Go Client for Google Maps Services is a Go Client library for the following Google Maps Platform
APIs:

- [Directions API]
- [Distance Matrix API]
- [Elevation API]
- [Geocoding API]
- [Places API]
- [Roads API]
- [Time Zone API]
- [Maps Static API]

> [!TIP]
> See the [Google Maps Platform Cloud Client Library for Go](https://github.com/googleapis/google-cloud-go/tree/main/maps) for our newer APIs
> including Address Validation API, Datasets API, Fleet Engine, new Places API, and Routes API.

## Requirements

- Go 1.7 or later.
- A Google Maps Platform [API key] from a project with the APIs below enabled.

> [!IMPORTANT]  
> This key should be kept secret on your server.

## Installation

To install the Go Client for Google Maps Services, please execute the following `go get` command.

```bash
go get googlemaps.github.io/maps
```

## Documentation

View the [reference documentation](https://godoc.org/googlemaps.github.io/maps).

Additional documentation about the APIs is available at:

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

## Features

### Rate limiting

Never sleep between requests again! By default, requests are sent at the expected rate limits for
each web service, typically 50 queries per second for free users. If you want to speed up or slow
down requests, you can do that too, using `maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(qps))`.

### Native types

Native objects for each of the API responses.

### Monitoring

It's possible to get metrics for status counts and latency histograms for monitoring.
Use `maps.WithMetricReporter(metrics.OpenCensusReporter{})` to log metrics to OpenCensus,
and `metrics.RegisterViews()` to make the metrics available to be exported.
OpenCensus can export these metrics to a [variety of monitoring services](https://opencensus.io/exporters/).
You can also implement your own metric reporter instead of using the provided one.

## Terms of Service

This library uses Google Maps Platform services, and any use of Google Maps Platform is subject to the [Terms of Service](https://cloud.google.com/maps-platform/terms).

For clarity, this library, and each underlying component, is not a Google Maps Platform Core Service.

## Support

This library is offered via an open source license. It is not governed by the Google Maps Platform Support [Technical Support Services Guidelines](https://cloud.google.com/maps-platform/terms/tssg), the [SLA](https://cloud.google.com/maps-platform/terms/sla), or the [Deprecation Policy](https://cloud.google.com/maps-platform/terms) (however, any Google Maps Platform services used by the library remain subject to the Google Maps Platform Terms of Service).

This library adheres to [semantic versioning](https://semver.org/) to indicate when backwards-incompatible changes are introduced.

If you find a bug, or have a feature request, please [file an issue][issues] on GitHub. If you would like to get answers to technical questions from other Google Maps Platform developers, ask through one of our [developer community channels](https://developers.google.com/maps/developer-community). If you'd like to contribute, please check the [Contributing guide][contrib].

You can also discuss this library on our [Discord server](https://discord.gg/hYsWbmk).

[API key]: https://developers.google.com/maps/documentation/places/web-service/get-api-key

[Google Maps Platform Web Services APIs]: https://developers.google.com/maps/apis-by-platform#web_service_apis
[Directions API]: https://developers.google.com/maps/documentation/directions/
[Distance Matrix API]: https://developers.google.com/maps/documentation/distancematrix/
[Elevation API]: https://developers.google.com/maps/documentation/elevation/
[Geocoding API]: https://developers.google.com/maps/documentation/geocoding/
[Places API]: https://developers.google.com/places/web-service/
[Roads API]: https://developers.google.com/maps/documentation/roads/
[Time Zone API]: https://developers.google.com/maps/documentation/timezone/
[Maps Static API]: https://developers.google.com/maps/documentation/maps-static/

[issues]: https://github.com/googlemaps/google-maps-services-go/issues
[contrib]: https://github.com/googlemaps/google-maps-services-go/blob/master/CONTRIB.md
