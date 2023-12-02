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

package maps

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTextSearchPizzaInNewYork(t *testing.T) {
	response := `{
  "html_attributions" : [],
  "next_page_token" : "CuQB1wAAANI17eHXt1HpqbLjkj7T5Ti69DEAClo02Qampg7Q6W_O_krFbge7hnTtDR7oVF3asexHcGnUtR1ZKjroYd4BTCXxSGPi9LEkjJ0P_zVE7byjEBcHvkdxB6nCHKHAgVNGqe0ZHuwSYKlr3C1-kuellMYwMlg3WSe69bJr1Ck35uToNZkUGvo4yjoYxNFRn1lABEnjPskbMdyHAjUDwvBDxzgGxpd8t0EzA9UOM8Y1jqWnZGJM7u8gacNFcI4prr0Doh9etjY1yHrgGYI4F7lKPbfLQKiks_wYzoHbcAcdbBjkEhAxDHC0XXQ16thDAlwVbEYaGhSaGDw5sHbaZkG9LZIqbcas0IJU8w",
  "results" : [
    {
      "formatted_address" : "60 Greenpoint Ave, Brooklyn, NY 11222, United States",
      "geometry" : {
        "location" : {
          "lat" : 40.729606,
          "lng" : -73.95857599999999
        }
      },
      "icon" : "https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png",
      "name" : "Paulie Gee's",
      "opening_hours" : {
        "open_now" : false,
        "weekday_text" : []
      },
      "photos" : [
        {
          "height" : 427,
          "html_attributions" : [
            "\u003ca href=\"https://maps.google.com/maps/contrib/107146711858841264424\"\u003ePaulie Gee&#39;s\u003c/a\u003e"
          ],
          "photo_reference" : "CmRdAAAAume6Q8oFq9AcGSZOQnqGHfgYHyCsQHO4JK-JbxeZ0rn1s-QeSMmLbFDV3NvWiSX3SOCJBLQnpnmpxCwiviSGdJbb6Ja2aqCKi5usrlMw6_wI_JM4eUe9_wsGhNT5MmPwEhDcY98HKcLeAkBLEvYHMja1GhQpQTCXtzKF8dLeyOhkm2XJmWJ2iA",
          "width" : 640
        }
      ],
      "place_id" : "ChIJuc8AM0BZwokRtpm2S66ltsE",
      "price_level" : 2,
      "rating" : 4.4,
      "types" : [ "restaurant", "food", "point_of_interest", "establishment" ]
    }
  ],
  "status" : "OK"
}`

	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &TextSearchRequest{
		Query: "Pizza in New York",
	}

	resp, err := c.TextSearch(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
	}

	nextPageToken := "CuQB1wAAANI17eHXt1HpqbLjkj7T5Ti69DEAClo02Qampg7Q6W_O_krFbge7hnTtDR7oVF3asexHcGnUtR1ZKjroYd4BTCXxSGPi9LEkjJ0P_zVE7byjEBcHvkdxB6nCHKHAgVNGqe0ZHuwSYKlr3C1-kuellMYwMlg3WSe69bJr1Ck35uToNZkUGvo4yjoYxNFRn1lABEnjPskbMdyHAjUDwvBDxzgGxpd8t0EzA9UOM8Y1jqWnZGJM7u8gacNFcI4prr0Doh9etjY1yHrgGYI4F7lKPbfLQKiks_wYzoHbcAcdbBjkEhAxDHC0XXQ16thDAlwVbEYaGhSaGDw5sHbaZkG9LZIqbcas0IJU8w"
	if resp.NextPageToken != nextPageToken {
		t.Errorf("expected %+v, was %+v", nextPageToken, resp.NextPageToken)
	}

	if len(resp.HTMLAttributions) != 0 {
		t.Errorf("expected %+v, was %+v", 0, len(resp.HTMLAttributions))
	}

	if len(resp.Results) != 1 {
		t.Errorf("expected %+v, was %+v", 1, len(resp.Results))
	}

	result := resp.Results[0]
	name := "Paulie Gee's"
	if name != result.Name {
		t.Errorf("expected %+v, was %+v", name, result.Name)
	}

	formattedAddress := "60 Greenpoint Ave, Brooklyn, NY 11222, United States"
	if formattedAddress != result.FormattedAddress {
		t.Errorf("expected %+v, was %+v", formattedAddress, result.FormattedAddress)
	}

	location := LatLng{Lat: 40.729606, Lng: -73.958576}
	if location != result.Geometry.Location {
		t.Errorf("expected %+v, was %+v", location, result.Geometry.Location)
	}

	icon := "https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png"
	if icon != result.Icon {
		t.Errorf("expected %+v, was %+v", icon, result.Icon)
	}

	placeID := "ChIJuc8AM0BZwokRtpm2S66ltsE"
	if placeID != result.PlaceID {
		t.Errorf("expected %+v, was %+v", placeID, result.PlaceID)
	}

	photo := result.Photos[0]
	photoWidth := 640
	if photoWidth != photo.Width {
		t.Errorf("expected %+v, was %+v", photoWidth, photo.Width)
	}

	photoHeight := 427
	if photoHeight != photo.Height {
		t.Errorf("expected %+v, was %+v", photoHeight, photo.Height)
	}

	photoReference := "CmRdAAAAume6Q8oFq9AcGSZOQnqGHfgYHyCsQHO4JK-JbxeZ0rn1s-QeSMmLbFDV3NvWiSX3SOCJBLQnpnmpxCwiviSGdJbb6Ja2aqCKi5usrlMw6_wI_JM4eUe9_wsGhNT5MmPwEhDcY98HKcLeAkBLEvYHMja1GhQpQTCXtzKF8dLeyOhkm2XJmWJ2iA"
	if photoReference != photo.PhotoReference {
		t.Errorf("expected %+v, was %+v", photoReference, photo.PhotoReference)
	}

	photoAttribution := "<a href=\"https://maps.google.com/maps/contrib/107146711858841264424\">Paulie Gee&#39;s</a>"
	if photoAttribution != photo.HTMLAttributions[0] {
		t.Errorf("expected %+v, was %+v", photoAttribution, photo.HTMLAttributions[0])
	}

	openNow := false
	if openNow != *result.OpeningHours.OpenNow {
		t.Errorf("expected %+v, was %+v", openNow, *result.OpeningHours.OpenNow)
	}

	// Find a way of mapping int -> PriceLevel
	priceLevel := 2
	if priceLevel != result.PriceLevel {
		t.Errorf("expected %+v, was %+v", priceLevel, result.PriceLevel)
	}
}

func TestNearbySearchMinimalRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&location=1%2C2&radius=10000"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &NearbySearchRequest{
		Location: &LatLng{1.0, 2.0},
		Radius:   10000,
	}

	_, err := c.NearbySearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestNearbySearchMaximalRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&keyword=foo&language=es&location=1%2C2&maxprice=3&minprice=0&name=name&opennow=true&pagetoken=NextPageToken&radius=10000&rankby=prominence&type=airport"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &NearbySearchRequest{
		Location:  &LatLng{1.0, 2.0},
		Radius:    10000,
		Keyword:   "foo",
		Language:  "es",
		MinPrice:  PriceLevelFree,
		MaxPrice:  PriceLevelExpensive,
		Name:      "name",
		OpenNow:   true,
		RankBy:    RankByProminence,
		Type:      PlaceTypeAirport,
		PageToken: "NextPageToken",
	}

	_, err := c.NearbySearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestNearbySearchNoLocation(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &NearbySearchRequest{
		Radius: 1000,
	}
	_, err := c.NearbySearch(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: Location and PageToken both missing")
	}
	if err.Error() != "maps: Location and PageToken both missing" {
		t.Errorf("Incorrect error returned \"%v\"", err)
	}
}

func TestNearbySearchNoRadius(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &NearbySearchRequest{
		Location: &LatLng{-33.865, 151.209},
	}
	_, err := c.NearbySearch(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: Radius and PageToken both missing")
	}
	if err.Error() != "maps: Radius and PageToken both missing" {
		t.Errorf("Incorrect error returned \"%v\"", err)
	}
}

func TestNearbySearchRankByDistanceAndRadius(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &NearbySearchRequest{
		Location: &LatLng{1.0, 2.0},
		Radius:   1000,
		RankBy:   RankByDistance,
	}
	_, err := c.NearbySearch(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: Radius specified with RankByDistance")
	}
	if err.Error() != "maps: Radius specified with RankByDistance" {
		t.Errorf("Incorrect error returned \"%v\"", err)
	}
}

func TestNearbySearchRankByDistanceAndNoKeywordNameAndType(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &NearbySearchRequest{
		Location: &LatLng{1.0, 2.0},
		RankBy:   RankByDistance,
	}
	_, err := c.NearbySearch(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: RankBy=distance and Keyword, Name and Type are missing")
	}

	if err.Error() != "maps: RankBy=distance and Keyword, Name and Type are missing" {
		t.Errorf("Incorrect error returned \"%v\"", err)
	}
}

func TestTextSearchMinimalRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&query=Pizza+in+New+York"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &TextSearchRequest{
		Query: "Pizza in New York",
	}

	_, err := c.TextSearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestTextSearchAllTheThingsRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&language=es&location=1%2C2&maxprice=2&minprice=0&opennow=true&pagetoken=NextPageToken&query=Pizza+in+New+York&radius=1000&region=US&type=airport"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &TextSearchRequest{
		Query:     "Pizza in New York",
		Location:  &LatLng{1.0, 2.0},
		Radius:    1000,
		Language:  "es",
		MinPrice:  PriceLevelFree,
		MaxPrice:  PriceLevelModerate,
		OpenNow:   true,
		Type:      PlaceTypeAirport,
		PageToken: "NextPageToken",
		Region:    "US",
	}

	_, err := c.TextSearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestTextSearchMissingQuery(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &TextSearchRequest{}
	_, err := c.TextSearch(context.Background(), r)

	if "maps: Query, PageToken and Type are all missing" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestTextSearchMissingRadius(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &TextSearchRequest{
		Query:    "Foo",
		Location: &LatLng{1, 2},
	}

	_, err := c.TextSearch(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: Radius missing, required with Location")
	}

	if "maps: Radius missing, required with Location" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestQueryAutocompleteMinimalRequestURL(t *testing.T) {
	expectedQuery := "input=quay+resteraunt+sydney&key=AIzaNotReallyAnAPIKey"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &QueryAutocompleteRequest{
		Input: "quay resteraunt sydney",
	}

	_, err := c.QueryAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestQueryAutocompleteMaximalRequestURL(t *testing.T) {
	expectedQuery := "input=quay+resteraunt+sydney&key=AIzaNotReallyAnAPIKey&language=es&location=1%2C2&offset=5&radius=10000"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &QueryAutocompleteRequest{
		Input:    "quay resteraunt sydney",
		Offset:   5,
		Location: &LatLng{1.0, 2.0},
		Radius:   10000,
		Language: "es",
	}

	_, err := c.QueryAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestQueryAutocompleteMissingInput(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &QueryAutocompleteRequest{}

	_, err := c.QueryAutocomplete(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: Input missing")
	}

	if "maps: Input missing" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestPlaceAutocompleteWithStrictbounds(t *testing.T) {
	session := NewPlaceAutocompleteSessionToken()
	expectedQuery := "input=Amoeba&key=AIzaNotReallyAnAPIKey&location=37.76999%2C-122.44696&radius=500&sessiontoken=" + uuid.UUID(session).String() + "&strictbounds=true&types=establishment"
	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &PlaceAutocompleteRequest{
		Input:        "Amoeba",
		Types:        AutocompletePlaceTypeEstablishment,
		Location:     &LatLng{37.76999, -122.44696},
		Radius:       500,
		StrictBounds: true,
		SessionToken: session,
	}

	_, err := c.PlaceAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	} else if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestPlaceAutocompleteMinimalRequestURL(t *testing.T) {
	session := NewPlaceAutocompleteSessionToken()
	expectedQuery := "input=quay+resteraunt+sydney&key=AIzaNotReallyAnAPIKey&sessiontoken=" + uuid.UUID(session).String()
	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	r := &PlaceAutocompleteRequest{
		Input:        "quay resteraunt sydney",
		SessionToken: session,
	}

	_, err := c.PlaceAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	} else if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestPlaceAutocompleteMaximalRequestURL(t *testing.T) {
	expectedQuery := "components=country%3AES%7Ccountry%3AAU&input=quay+resteraunt+sydney&key=AIzaNotReallyAnAPIKey&language=es&location=1%2C2&offset=5&origin=1%2C2&radius=10000&types=geocode"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	placeType, err := ParseAutocompletePlaceType("geocode")
	if err != nil {
		t.Errorf("Unexpected error in parsing place type: %v", err)
	}

	r := &PlaceAutocompleteRequest{
		Input:    "quay resteraunt sydney",
		Offset:   5,
		Location: &LatLng{1.0, 2.0},
		Origin:   &LatLng{1.0, 2.0},
		Radius:   10000,
		Language: "es",
		Types:    placeType,
		Components: map[Component][]string{
			ComponentCountry: {"ES", "AU"},
		},
	}

	_, err = c.PlaceAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	} else if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestPlaceAutocompleteWithMutipleComponentsWillWithoutMutation(t *testing.T) {
	expectedQuery := "components=country%3AES%7Ccountry%3AAU%7Ccountry%3AHK&input=quay+resteraunt+sydney&key=AIzaNotReallyAnAPIKey&language=es&location=1%2C2&offset=5&radius=10000&types=geocode"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))

	placeType, err := ParseAutocompletePlaceType("geocode")
	if err != nil {
		t.Errorf("Unexpected error in parsing place type: %v", err)
	}

	r := &PlaceAutocompleteRequest{
		Input:    "quay resteraunt sydney",
		Offset:   5,
		Location: &LatLng{1.0, 2.0},
		Radius:   10000,
		Language: "es",
		Types:    placeType,
		Components: map[Component][]string{
			ComponentCountry: {"ES", "AU", "HK"},
		},
	}

	var copiedComponents = map[Component][]string{
		ComponentCountry: {"ES", "AU", "HK"},
	}

	// copy the request
	var originalReq = *r
	originalReq.Components = copiedComponents

	_, err = c.PlaceAutocomplete(context.Background(), r)

	if eq := reflect.DeepEqual(*r, originalReq); !eq {
		t.Errorf("Unexpected mutation at PlaceAutocompleteRequest")
	}
	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	} else if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestPlaceAutocompleteMissingInput(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &PlaceAutocompleteRequest{}

	_, err := c.PlaceAutocomplete(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: Input missing")
	}

	if "maps: Input missing" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestPlaceAutocompleteWithStructuredFormatting(t *testing.T) {
	session := NewPlaceAutocompleteSessionToken()
	response := `
{
  "predictions": [
    {
      "description": "Theater de Meervaart, Meer en Vaart, Amsterdam, Netherlands",
      "id": "2d13bd84619a4cc5f9bf0b20f093b841a1403fcd",
      "matched_substrings": [
        {
          "length": 20,
          "offset": 0
        }
      ],
      "place_id": "ChIJVwucBNLjxUcRXzGhUau_gBw",
      "reference": "ClRJAAAAq5qHSaGPrUhUH3LyKrLYmg280v2TYXUCD5h7_m0YGw3Y8Mj1h1bffMyG7CBFlAN17V8kKkzeXwXO94v5513ErtXHVYKnJ9pNg4S7HtGUqEwSECL5WbMbXSbeRs_H2B91qHcaFEbgpLF1aftugYKgJTIupUwYsEbl",
      "structured_formatting": {
        "main_text": "Theater de Meervaart",
        "main_text_matched_substrings": [
          {
            "length": 20,
            "offset": 0
          }
        ],
        "secondary_text": "Meer en Vaart, Amsterdam, Netherlands"
      },
      "terms": [
        {
          "offset": 0,
          "value": "Theater de Meervaart"
        },
        {
          "offset": 22,
          "value": "Meer en Vaart"
        },
        {
          "offset": 37,
          "value": "Amsterdam"
        },
        {
          "offset": 48,
          "value": "Netherlands"
        }
      ],
      "types": [
        "establishment"
      ]
    }
  ],
  "status": "OK"
}`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &PlaceAutocompleteRequest{
		Input:        "Theater de Meervaart",
		Types:        AutocompletePlaceType("establishment"),
		SessionToken: session,
	}

	resp, err := c.PlaceAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
		return
	}

	mainText := "Theater de Meervaart"
	if mainText != resp.Predictions[0].StructuredFormatting.MainText {
		t.Errorf("expected %+v, was %+v", mainText, resp.Predictions[0].StructuredFormatting.MainText)
	}

	secondaryText := "Meer en Vaart, Amsterdam, Netherlands"
	if secondaryText != resp.Predictions[0].StructuredFormatting.SecondaryText {
		t.Errorf("expected %+v, was %+v", mainText, resp.Predictions[0].StructuredFormatting.SecondaryText)
	}

	mainTextSubstringLength := 20
	if mainTextSubstringLength != resp.Predictions[0].StructuredFormatting.MainTextMatchedSubstrings[0].Length {
		t.Errorf("expected %+v, was %+v", mainTextSubstringLength, resp.Predictions[0].StructuredFormatting.MainTextMatchedSubstrings[0].Length)
	}

	mainTextSubstringOffset := 0
	if mainTextSubstringOffset != resp.Predictions[0].StructuredFormatting.MainTextMatchedSubstrings[0].Offset {
		t.Errorf("expected %+v, was %+v", mainTextSubstringLength, resp.Predictions[0].StructuredFormatting.MainTextMatchedSubstrings[0].Offset)
	}
}

func TestPlaceDetails(t *testing.T) {
	response := `
{
    "html_attributions": [],
    "result": {
        "business_status": "OPERATIONAL",
        "dine_in": true,
        "formatted_address": "Upper Level Overseas Passenger Terminal, The Rocks NSW 2000, Australia",
        "geometry": {
            "location": {
                "lat": -33.85756870000001,
                "lng": 151.2100844
            },
            "viewport": {
                "northeast": {
                    "lat": -33.8561714197085,
                    "lng": 151.2113682802915
                },
                "southwest": {
                    "lat": -33.85886938029149,
                    "lng": 151.2086703197085
                }
            }
        },
        "name": "Quay Restaurant",
        "opening_hours": {
            "open_now": false,
            "periods": [
                {
                    "close": {
                        "day": 0,
                        "time": "2045"
                    },
                    "open": {
                        "day": 0,
                        "time": "1200"
                    }
                },
                {
                    "close": {
                        "day": 4,
                        "time": "2045"
                    },
                    "open": {
                        "day": 4,
                        "time": "1800"
                    }
                },
                {
                    "close": {
                        "day": 5,
                        "time": "2045"
                    },
                    "open": {
                        "day": 5,
                        "time": "1800"
                    }
                },
                {
                    "close": {
                        "day": 6,
                        "time": "2045"
                    },
                    "open": {
                        "day": 6,
                        "time": "1200"
                    }
                }
            ],
            "weekday_text": [
                "Monday: Closed",
                "Tuesday: Closed",
                "Wednesday: Closed",
                "Thursday: 6:00 – 8:45 PM",
                "Friday: 6:00 – 8:45 PM",
                "Saturday: 12:00 – 8:45 PM",
                "Sunday: 12:00 – 8:45 PM"
            ]
        },
	"secondary_opening_hours": [
            {
                "open_now": true,
                "periods": [
                    {
                        "close": {
                            "day": 0,
                            "time": "2045"
                        },
                        "open": {
                            "day": 0,
                            "time": "1200"
                        }
                    },
                    {
                        "close": {
                            "day": 4,
                            "time": "2045"
                        },
                        "open": {
                            "day": 4,
                            "time": "1800"
                        }
                    },
                    {
                        "close": {
                            "day": 5,
                            "time": "2045"
                        },
                        "open": {
                            "day": 5,
                            "time": "1800"
                        }
                    },
                    {
                        "close": {
                            "day": 6,
                            "time": "2045"
                        },
                        "open": {
                            "day": 6,
                            "time": "1200"
                        }
                    }
                ],
                "type": "KITCHEN",
                "weekday_text": [
                    "Monday: 11:30 AM – 2:00 PM, 5:00 – 10:00 PM",
                    "Tuesday: 11:30 AM – 2:00 PM, 5:00 – 10:00 PM",
                    "Wednesday: 11:30 AM – 2:00 PM, 5:00 – 10:00 PM",
                    "Thursday: 11:30 AM – 2:00 PM, 5:00 – 10:00 PM",
                    "Friday: 11:30 AM – 2:00 PM, 5:00 – 10:00 PM",
                    "Saturday: 11:30 AM – 2:00 PM, 5:00 – 10:00 PM",
                    "Sunday: 11:30 AM – 2:00 PM, 5:00 – 9:00 PM"
                ]
            }
	    ],
        "place_id": "ChIJ4cQcDV2uEmsRMxTEHBIe9ZQ",
        "serves_dinner": true,
        "utc_offset": 660,
        "wheelchair_accessible_entrance": true
    },
    "status": "OK"
}
`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	placeID := "ChIJ4cQcDV2uEmsRMxTEHBIe9ZQ"
	fields := []PlaceDetailsFieldMask{PlaceDetailsFieldMaskBusinessStatus, PlaceDetailsFieldMaskDineIn, PlaceDetailsFieldMaskFormattedAddress, PlaceDetailsFieldMaskGeometry, PlaceDetailsFieldMaskName, PlaceDetailsFieldMaskCurrentOpeningHours, PlaceDetailsFieldMaskSecondaryOpeningHours, PlaceDetailsFieldMaskPlaceID, PlaceDetailsFieldMaskServesDinner, PlaceDetailsFieldMaskUTCOffset, PlaceDetailsFieldMaskWheelchairAccessibleEntrance}
	r := &PlaceDetailsRequest{
		PlaceID: placeID,
		Fields:  fields,
	}

	resp, err := c.PlaceDetails(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
		return
	}

	businessStatus := "OPERATIONAL"
	if businessStatus != resp.BusinessStatus {
		t.Errorf("expected %+v, was %+v", businessStatus, resp.BusinessStatus)
	}

	if !*&resp.DineIn {
		t.Errorf("Expected DineIn to be true")
	}

	formattedAddress := "Upper Level Overseas Passenger Terminal, The Rocks NSW 2000, Australia"
	if formattedAddress != resp.FormattedAddress {
		t.Errorf("expected %+v, was %+v", formattedAddress, resp.FormattedAddress)
	}

	name := "Quay Restaurant"
	if name != resp.Name {
		t.Errorf("expected %+v, was %+v", name, resp.Name)
	}

	if resp.OpeningHours.Periods[0].Open.Day != time.Sunday || resp.OpeningHours.Periods[0].Close.Day != time.Sunday {
		t.Errorf("OpeningHours.Periods[0].Open.Day or Close.Day incorrect")
	}

	if resp.OpeningHours.Periods[0].Open.Time != "1200" || resp.OpeningHours.Periods[0].Close.Time != "2045" {
		t.Errorf("OpeningHours.Periods[0].Open.Time or Close.Time incorrect")
	}

	if resp.SecondaryOpeningHours[0].Periods[0].Open.Day != time.Sunday || resp.OpeningHours.Periods[0].Close.Day != time.Sunday {
		t.Errorf("SecondaryOpeningHours[0].Periods[0].Open.Day or Close.Day incorrect")
	}

	if resp.SecondaryOpeningHours[0].Periods[0].Open.Time != "1200" || resp.OpeningHours.Periods[0].Close.Time != "2045" {
		t.Errorf("SecondaryOpeningHours[0].Periods[0].Open.Time or Close.Time incorrect")
	}

	weekdayText := "Monday: Closed"
	if weekdayText != resp.OpeningHours.WeekdayText[0] {
		t.Errorf("expected %+v, was %+v", weekdayText, resp.OpeningHours.WeekdayText[0])
	}

	if placeID != resp.PlaceID {
		t.Errorf("expected %+v, was %+v", placeID, resp.PlaceID)
	}
  
	if !*&resp.ServesDinner {
		t.Errorf("Expected ServesDinner to be true")
	}

	if *resp.UTCOffset != 660 {
		t.Errorf("Expected UTCOffset to be 660")
	}

	if !*&resp.WheelchairAccessibleEntrance {
		t.Errorf("Expected WheelchairAccessibleEntrance to be true")
	}

}

func TestPlaceDetailsMissingPlaceID(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &PlaceDetailsRequest{}

	_, err := c.PlaceDetails(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: PlaceID missing")
	}

	if "maps: PlaceID missing" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestPlacePhotoMissingPhotoReference(t *testing.T) {
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &PlacePhotoRequest{}

	_, err := c.PlacePhoto(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: PhotoReference missing")
	}

	if "maps: PhotoReference missing" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestPlacePhotoMissingWidthAndHeight(t *testing.T) {
	photoReference := "ThisIsNotAPhotoReference"
	c, _ := NewClient(WithAPIKey(apiKey))
	r := &PlacePhotoRequest{
		PhotoReference: photoReference,
	}

	_, err := c.PlacePhoto(context.Background(), r)

	if err == nil {
		t.Errorf("Error expected: maps: both MaxHeight & MaxWidth missing")
	}

	if "maps: both MaxHeight & MaxWidth missing" != err.Error() {
		t.Errorf("Wrong error returned \"%v\"", err)
	}
}

func TestTextSearchWithPermanentlyClosed(t *testing.T) {
	response := `
	{
	   "html_attributions" : [],
	   "results" : [
	      {
	         "formatted_address" : "5 Martinez Ave, West End QLD 4810, Australia",
	         "geometry" : {
	            "location" : {
	               "lat" : -19.2690427,
	               "lng" : 146.7832313
	            }
	         },
	         "icon" : "https://maps.gstatic.com/mapfiles/place_api/icons/school-71.png",
	         "id" : "6b19d85f4ac1dd71ba400d8ad7fe540a64beacc7",
	         "name" : "ABC Learning Centre",
			 "permanently_closed" : true,
			 "business_status": "foo",
	         "place_id" : "ChIJLdqTiaj51WsRv4Mkbq2qQEU",
	         "reference" : "CnRmAAAAJJuaK6n6aI7imGz2zcqHpBanTQcafAIyja-5pGX6q67WDRT4DJ8M6HcjfxRCbOM-7RAw10sU9l-lZktErhP4mVmavboCyI_QG8iAHNjBPlqYcfFYjJLUE4gtrYvYhx1VGG88wYBbQXXAH4hcGQc3-xIQyNcdcFc9rmijjlL5g1U4KxoUYxqZLWwPfDWy1hkU0DqTUbAm26k",
	         "types" : [ "school", "point_of_interest", "establishment" ]
	      }
	   ],
	   "status" : "OK"
  }`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &TextSearchRequest{
		Query: "ABC Learning Centres in australia",
	}

	resp, err := c.TextSearch(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
		return
	}

	result := resp.Results[0]

	formattedAddress := "5 Martinez Ave, West End QLD 4810, Australia"
	if formattedAddress != result.FormattedAddress {
		t.Errorf("expected %+v, was %+v", formattedAddress, result.FormattedAddress)
	}

	icon := "https://maps.gstatic.com/mapfiles/place_api/icons/school-71.png"
	if icon != result.Icon {
		t.Errorf("expected %+v, was %+v", icon, result.Icon)
	}

	name := "ABC Learning Centre"
	if name != result.Name {
		t.Errorf("expected %+v, was %+v", name, result.Name)
	}

	permanentlyClosed := true
	if permanentlyClosed != result.PermanentlyClosed {
		t.Errorf("expected %+v, was %+v", permanentlyClosed, result.PermanentlyClosed)
	}

	businessStatus := "foo"
	if businessStatus != result.BusinessStatus {
		t.Errorf("expected %+v, was %+v", businessStatus, result.BusinessStatus)
	}

	placeID := "ChIJLdqTiaj51WsRv4Mkbq2qQEU"
	if placeID != result.PlaceID {
		t.Errorf("expected %+v, was %+v", placeID, result.PlaceID)
	}
}

func TestPlaceAutocompleteJsonMarshalLowerCase(t *testing.T) {
	session := NewPlaceAutocompleteSessionToken()
	response := `
{
  "predictions": [
    {
      "description": "Theater de Meervaart, Meer en Vaart, Amsterdam, Netherlands",
      "id": "2d13bd84619a4cc5f9bf0b20f093b841a1403fcd",
      "matched_substrings": [
        {
          "length": 20,
          "offset": 0
        }
      ],
      "place_id": "ChIJVwucBNLjxUcRXzGhUau_gBw",
      "reference": "ClRJAAAAq5qHSaGPrUhUH3LyKrLYmg280v2TYXUCD5h7_m0YGw3Y8Mj1h1bffMyG7CBFlAN17V8kKkzeXwXO94v5513ErtXHVYKnJ9pNg4S7HtGUqEwSECL5WbMbXSbeRs_H2B91qHcaFEbgpLF1aftugYKgJTIupUwYsEbl",
      "structured_formatting": {
        "main_text": "Theater de Meervaart",
        "main_text_matched_substrings": [
          {
            "length": 20,
            "offset": 0
          }
        ],
        "secondary_text": "Meer en Vaart, Amsterdam, Netherlands"
      },
      "terms": [
        {
          "offset": 0,
          "value": "Theater de Meervaart"
        },
        {
          "offset": 22,
          "value": "Meer en Vaart"
        },
        {
          "offset": 37,
          "value": "Amsterdam"
        },
        {
          "offset": 48,
          "value": "Netherlands"
        }
      ],
      "types": [
        "establishment"
      ]
    }
  ],
  "status": "OK"
}`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &PlaceAutocompleteRequest{
		Input:        "Theater de Meervaart",
		Types:        AutocompletePlaceType("establishment"),
		SessionToken: session,
	}

	resp, err := c.PlaceAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
		return
	}

	json, err := json.Marshal(&resp)

	if err != nil {
		t.Errorf("json.Marshal error: %v", err)
		return
	}

	if strings.Contains(string(json), `"Predictions"`) {
		t.Error("AutocompleteResponse json.Marshal result \"prediction\" key is uppercase")
		return
	}

	if strings.Contains(string(json), `"predictions"`) {
		return
	}

	t.Error("TestPlaceAutocompleteJsonMarshalLowerCase error!")
}

func TestFindPlaceFromText(t *testing.T) {
	expectedQuery := "fields=photos%2Cformatted_address%2Cname%2Copening_hours%2Crating&input=mongolian+grill&inputtype=textquery&key=AIzaNotReallyAnAPIKey&locationbias=circle%3A2000%4047.6918452%2C-122.2226413"
	response := `
{
	"candidates" : [
	   {
		  "formatted_address" : "9736 NE 117th Ln, Kirkland, WA 98034, USA",
		  "name" : "Mongolian Grill Kirkland",
		  "opening_hours" : {
			 "open_now" : false,
			 "weekday_text" : []
		  },
		  "photos" : [
			 {
				"height" : 2891,
				"html_attributions" : [
				   "\u003ca href=\"https://maps.google.com/maps/contrib/111759700246215860219/photos\"\u003eVamsi Kanamaluru\u003c/a\u003e"
				],
				"photo_reference" : "CmRaAAAAwzjnmCwlQAFViioiTzU3jGb1jzTnfUg3CThLhA92w9FeLvCFymiYgL3qlstXd0TngcZ45fF3mwJfPWHWKQ44rllAcC_Izp4A-euYZloBnjFAtEuKOx5gecBG5rR0CnymEhB0LxSGBDoojumIma5k6pudGhQdyUwhjplZjF1StMfaydwbGFE80Q",
				"width" : 3175
			 }
		  ],
		  "rating" : 4.2
	   }
	],
	"debug_log" : {
	   "line" : []
	},
	"status" : "OK"
 }`
	server := mockServerForQuery(expectedQuery, 200, response)
	defer server.s.Close()
	fields := []PlaceSearchFieldMask{PlaceSearchFieldMaskPhotos, PlaceSearchFieldMaskFormattedAddress, PlaceSearchFieldMaskName, PlaceSearchFieldMaskOpeningHours, PlaceSearchFieldMaskRating}

	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))
	r := &FindPlaceFromTextRequest{
		Input:              "mongolian grill",
		InputType:          FindPlaceFromTextInputTypeTextQuery,
		Fields:             fields,
		LocationBias:       FindPlaceFromTextLocationBiasCircular,
		LocationBiasCenter: &LatLng{47.6918452, -122.2226413},
		LocationBiasRadius: 2000,
	}

	resp, err := c.FindPlaceFromText(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
		return
	}

	if 1 != len(resp.Candidates) {
		t.Errorf("expected %+v, was %+v", 1, len(resp.Candidates))
	}

	if "9736 NE 117th Ln, Kirkland, WA 98034, USA" != resp.Candidates[0].FormattedAddress {
		t.Errorf("expected %+v, was %+v", "9736 NE 117th Ln, Kirkland, WA 98034, USA", resp.Candidates[0].FormattedAddress)
	}

	if "Mongolian Grill Kirkland" != resp.Candidates[0].Name {
		t.Errorf("expected %+v, was %+v", "Mongolian Grill Kirkland", resp.Candidates[0].Name)
	}
}

func TestTextSearchZeroResults(t *testing.T) {
	server := mockServer(200, `{"results" : [], "status" : "ZERO_RESULTS"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &TextSearchRequest{
		Query: "Nothing to see here",
	}

	resp, err := c.TextSearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error for ZERO_RESULTS status")
	}

	if len(resp.Results) != 0 {
		t.Errorf("Unexpected results for ZERO_RESULTS status")
	}
}

func TestNearbySearchZeroResults(t *testing.T) {
	server := mockServer(200, `{"results" : [], "status" : "ZERO_RESULTS"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &NearbySearchRequest{
		Location: &LatLng{28.0, 140.0},
		Radius:   100,
	}

	resp, err := c.NearbySearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error for ZERO_RESULTS status")
	}

	if len(resp.Results) != 0 {
		t.Errorf("Unexpected results for ZERO_RESULTS status")
	}
}

func TestFindPlaceFromTextZeroResults(t *testing.T) {
	server := mockServer(200, `{"candidates" : [], "status" : "ZERO_RESULTS"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &FindPlaceFromTextRequest{
		Input:     "+12345506789",
		InputType: FindPlaceFromTextInputTypeTextQuery,
	}

	resp, err := c.FindPlaceFromText(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error for ZERO_RESULTS status")
	}

	if len(resp.Candidates) != 0 {
		t.Errorf("Unexpected candidates for ZERO_RESULTS status")
	}
}

func TestPlaceAutocompleteZeroResults(t *testing.T) {
	server := mockServer(200, `{"predictions" : [], "status" : "ZERO_RESULTS"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &PlaceAutocompleteRequest{
		Input: "gobbledygook",
	}

	resp, err := c.PlaceAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error for ZERO_RESULTS status")
	}

	if len(resp.Predictions) != 0 {
		t.Errorf("Unexpected predictions for ZERO_RESULTS status")
	}
}

func TestQueryAutocompleteZeroResults(t *testing.T) {
	server := mockServer(200, `{"predictions" : [], "status" : "ZERO_RESULTS"}`)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &QueryAutocompleteRequest{
		Input: "gobbledygook",
	}

	resp, err := c.QueryAutocomplete(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error for ZERO_RESULTS status")
	}

	if len(resp.Predictions) != 0 {
		t.Errorf("Unexpected predictions for ZERO_RESULTS status")
	}
}
