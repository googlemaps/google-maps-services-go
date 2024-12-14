[![GoDoc](https://godoc.org/googlemaps.github.io/maps?status.svg)][documentation]
[![Go Report Card](https://goreportcard.com/badge/github.com/googlemaps/google-maps-services-go)](https://goreportcard.com/report/github.com/googlemaps/google-maps-services-go)
![GitHub Go version](https://img.shields.io/github/go-mod/go-version/googlemaps/google-maps-services-go)

![Release](https://github.com/googlemaps/google-maps-services-go/workflows/Release/badge.svg)
![Stable](https://img.shields.io/badge/stability-stable-green)
[![Tests/Build Status](https://github.com/googlemaps/google-maps-services-go/actions/workflows/test.yml/badge.svg)](https://github.com/googlemaps/google-maps-services-go/actions/workflows/test.yml)

![GitHub contributors](https://img.shields.io/github/contributors/googlemaps/google-maps-services-go?color=green)
[![GitHub License](https://img.shields.io/github/license/googlemaps/google-maps-services-go?color=blue)][license]
[![Discord](https://img.shields.io/discord/676948200904589322?color=6A7EC2&logo=discord&logoColor=ffffff)][Discord server]
[![StackOverflow](https://img.shields.io/stackexchange/stackoverflow/t/google-maps?color=orange&label=google-maps&logo=stackoverflow)](https://stackoverflow.com/questions/tagged/google-maps)

# Go Client for Google Maps Services

## Description

Use Go? Want to [geocode][Geocoding API] something? Looking for [directions][Directions API]? This client library brings the following [Google Maps Web Services APIs] to your server-side Go applications:

- [Maps Static API]
- [Directions API]
- [Distance Matrix API]
- [Elevation API]
- [Geocoding API]
- [Places API]
- [Roads API]
- [Time Zone API]

> [!TIP]
> See the [Google Maps Platform Cloud Client Library for Go](https://github.com/googleapis/google-cloud-go/tree/main/maps) for our newer APIs
> including Address Validation API, Datasets API, Fleet Engine, new Places API, and Routes API.

## Requirements

* [Sign up with Google Maps Platform]
* A Google Maps Platform [project] with the desired API(s) from the above list enabled
* An [API key] associated with the project above
* Go 1.7+

## API Key Security

This client library is designed for use in server-side applications.

In either case, it is important to add [API key restrictions](https://developers.google.com/maps/api-security-best-practices#restricting-api-keys) to improve its security. Additional security measures, such as hiding your key
from version control, should also be put in place to further improve the security of your API key.

Check out the [API Security Best Practices](https://developers.google.com/maps/api-security-best-practices) guide to learn more.

## Installation

To install the client library, execute the following command:

    go get googlemaps.github.io/maps

## Documentation

You can find the reference [documentation] at GoDoc, and each API also has its own set of documentation:
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
	c, err := maps.NewClient(maps.WithAPIKey("YOUR_API_KEY"))
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

Sample usage of the Geocoding API with an API key to get an [Address Descriptor](https://developers.google.com/maps/documentation/geocoding/address-descriptors/requests-address-descriptors):

```go
package main

import (
	"context"
	"log"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

func main() {
	c, err := maps.NewClient(maps.WithAPIKey("YOUR_API_KEY"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.GeocodingRequest{
		LatLng: &LatLng{Lat: 40.714224, Lng: -73.961452},
		EnableAddressDescriptor: True
	}
	reverseGeocodingResponse, _, err := c.ReverseGeocode(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pretty.Println(reverseGeocodingResponse)
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

## Contributing

Contributions are welcome and encouraged! If you'd like to contribute, send us a [pull request] and refer to our [code of conduct] and [contributing guide].

## Terms of Service

This library uses Google Maps Platform services. Use of Google Maps Platform services through this library is subject to the Google Maps Platform [Terms of Service].

This library is not a Google Maps Platform Core Service. Therefore, the Google Maps Platform Terms of Service (e.g. Technical Support Services, Service Level Agreements, and Deprecation Policy) do not apply to the code in this library.

## Support

This library is offered via an open source [license]. It is not governed by the Google Maps Platform Support [Technical Support Services Guidelines, the SLA, or the [Deprecation Policy]. However, any Google Maps Platform services used by the library remain subject to the Google Maps Platform Terms of Service.

This library adheres to [semantic versioning] to indicate when backwards-incompatible changes are introduced. Accordingly, while the library is in version 0.x, backwards-incompatible changes may be introduced at any time.

If you find a bug, or have a feature request, please [file an issue] on GitHub. If you would like to get answers to technical questions from other Google Maps Platform developers, ask through one of our [developer community channels]. If you'd like to contribute, please check the [contributing guide].

You can also discuss this library on our [Discord server].

[Google Maps Platform Web Services APIs]: https://developers.google.com/maps/apis-by-platform#web_service_apis
[Maps Static API]: https://developers.google.com/maps/documentation/maps-static
[Directions API]: https://developers.google.com/maps/documentation/directions
[Distance Matrix API]: https://developers.google.com/maps/documentation/distancematrix
[Elevation API]: https://developers.google.com/maps/documentation/elevation
[Geocoding API]: https://developers.google.com/maps/documentation/geocoding
[Places API]: https://developers.google.com/places/web-service
[Roads API]: https://developers.google.com/maps/documentation/roads
[Time Zone API]: https://developers.google.com/maps/documentation/timezone

[API key]: https://developers.google.com/maps/documentation/javascript/get-api-key
[documentation]: https://godoc.org/googlemaps.github.io/maps

[code of conduct]: CODE_OF_CONDUCT.md
[contributing guide]: CONTRIB.md
[Deprecation Policy]: https://cloud.google.com/maps-platform/terms
[developer community channels]: https://developers.google.com/maps/developer-community
[Discord server]: https://discord.gg/hYsWbmk
[file an issue]: https://github.com/googlemaps/google-maps-services-go/issues/new/choose
[license]: LICENSE
[pull request]: https://github.com/googlemaps/google-maps-services-go/compare
[project]: https://developers.google.com/maps/documentation/javascript/cloud-setup#enabling-apis
[semantic versioning]: https://semver.org
[Sign up with Google Maps Platform]: https://console.cloud.google.com/google/maps-apis/start
[similar inquiry]: https://github.com/googlemaps/google-maps-services-go/issues
[SLA]: https://cloud.google.com/maps-platform/terms/sla
[Technical Support Services Guidelines]: https://cloud.google.com/maps-platform/terms/tssg
[Terms of Service]: https://cloud.google.com/maps-platform/terms
