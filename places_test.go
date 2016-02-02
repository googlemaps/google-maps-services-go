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
	"testing"
	"time"

	"golang.org/x/net/context"
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
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
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

func TestTextSearchMinimalRequestURL(t *testing.T) {
	expectedQuery := "key=AIzaNotReallyAnAPIKey&query=Pizza+in+New+York"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.s.URL

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
	expectedQuery := "key=AIzaNotReallyAnAPIKey&language=es&location=1%2C2&maxprice=2&minprice=0&opennow=true&pagetoken=NextPageToken&query=Pizza+in+New+York&radius=1000"

	server := mockServerForQuery(expectedQuery, 200, `{"status":"OK"}"`)
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.s.URL

	r := &TextSearchRequest{
		Query:     "Pizza in New York",
		Location:  &LatLng{1.0, 2.0},
		Radius:    1000,
		Language:  "es",
		MinPrice:  PriceLevelFree,
		MaxPrice:  PriceLevelModerate,
		OpenNow:   true,
		PageToken: "NextPageToken",
	}

	_, err := c.TextSearch(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
	}
}

func TestPlaceDetails(t *testing.T) {
	response := `
{
   "html_attributions" : [],
   "result" : {
      "address_components" : [],
      "formatted_address" : "3, Overseas Passenger Terminal, George St & Argyle Street, The Rocks NSW 2000, Australia",
      "formatted_phone_number" : "(02) 9251 5600",
      "geometry" : {
         "location" : {
            "lat" : -33.858018,
            "lng" : 151.210091
         }
      },
      "icon" : "https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png",
      "international_phone_number" : "+61 2 9251 5600",
      "name" : "Quay",
      "opening_hours" : {
         "open_now" : true,
         "periods" : [
            {
               "close" : {
                  "day" : 1,
                  "time" : "1700"
               },
               "open" : {
                  "day" : 1,
                  "time" : "1330"
               }
            }
         ],
         "weekday_text" : [
            "Monday: 1:30 – 5:00 pm"
         ]
      },
      "photos" : [
         {
            "height" : 612,
            "html_attributions" : [
               "\u003ca href=\"https://maps.google.com/maps/contrib/107255044321733286691\"\u003eFrom a Google User\u003c/a\u003e"
            ],
            "photo_reference" : "CmRdAAAAm1qTaarpM_sUatFI7JxjwxVTgKCGSjz62q_vHpNMoZDP3PpBHGW-rAHQEEprl_c1MyvXFhvZb2mXj8yhKvnEMsSveb-cMuDaDgS7LS8sPPrMrt5s_Mx0G0ereom3j6KxEhAkaQH1_nWxpl4W2mFZ1CKoGhQV_Jx9MIn0skBS3tRAuIFzgHARww",
            "width" : 816
         }
      ],
      "place_id" : "ChIJ02qnq0KuEmsRHUJF4zo1x4I",
      "price_level" : 4,
      "rating" : 4.1,
      "reviews" : [
         {
            "aspects" : [
               {
                  "rating" : 1,
                  "type" : "overall"
               }
            ],
            "author_name" : "Rachel Lewis",
            "author_url" : "https://plus.google.com/114299517944848975298",
            "language" : "en",
            "rating" : 3,
            "text" : "Overall disappointing. This is the second time i've been there and my experience was... Nothing to nibble on for 45 mins and then the bread came. My first entree was the marron which I thought was tasteless - perhaps others would say delicate? but there you go. The XO sea was fantastic. I chose the  vegetarian main dish which was all about the texture which was great but nothing at all outstanding about the dish. My husband and daughter chose the duck for their main course it was the smallest main course i've ever seen - their faces were priceless when it arrived!. Snow egg was beautiful but the granita on the bottom had some solid chunks of hard ice. The service was quite good...",
            "time" : 1441848853
         }
      ],
      "scope" : "GOOGLE",
      "types" : [ "restaurant", "food", "point_of_interest", "establishment" ],
      "url" : "https://plus.google.com/105746337161979416551/about?hl=en-US",
      "user_ratings_total" : 275,
      "utc_offset" : 660,
      "vicinity" : "3 Overseas Passenger Terminal, George Street, The Rocks",
      "website" : "http://www.quay.com.au/"
   },
   "status" : "OK"
}
`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
	placeID := "ChIJ02qnq0KuEmsRHUJF4zo1x4I"
	r := &PlaceDetailsRequest{
		PlaceID: placeID,
	}

	resp, err := c.PlaceDetails(context.Background(), r)

	if err != nil {
		t.Errorf("r.Get returned non nil error: %v", err)
		return
	}

	formattedAddress := "3, Overseas Passenger Terminal, George St & Argyle Street, The Rocks NSW 2000, Australia"
	if formattedAddress != resp.FormattedAddress {
		t.Errorf("expected %+v, was %+v", formattedAddress, resp.FormattedAddress)
	}

	formattedPhoneNumber := "(02) 9251 5600"
	if formattedPhoneNumber != resp.FormattedPhoneNumber {
		t.Errorf("expected %+v, was %+v", formattedPhoneNumber, resp.FormattedPhoneNumber)
	}

	icon := "https://maps.gstatic.com/mapfiles/place_api/icons/restaurant-71.png"
	if icon != resp.Icon {
		t.Errorf("expected %+v, was %+v", icon, resp.Icon)
	}

	internationalPhoneNumber := "+61 2 9251 5600"
	if internationalPhoneNumber != resp.InternationalPhoneNumber {
		t.Errorf("expected %+v, was %+v", internationalPhoneNumber, resp.InternationalPhoneNumber)
	}

	name := "Quay"
	if name != resp.Name {
		t.Errorf("expected %+v, was %+v", name, resp.Name)
	}

	if !*resp.OpeningHours.OpenNow {
		t.Errorf("Expected OpenNow to be true")
	}

	if resp.OpeningHours.Periods[0].Open.Day != time.Monday || resp.OpeningHours.Periods[0].Close.Day != time.Monday {
		t.Errorf("OpeningHours.Periods[0].Open.Day or Close.Day incorrect")
	}

	if resp.OpeningHours.Periods[0].Open.Time != "1330" || resp.OpeningHours.Periods[0].Close.Time != "1700" {
		t.Errorf("OpeningHours.Periods[0].Open.Time or Close.Time incorrect")
	}

	weekdayText := "Monday: 1:30 – 5:00 pm"
	if weekdayText != resp.OpeningHours.WeekdayText[0] {
		t.Errorf("expected %+v, was %+v", weekdayText, resp.OpeningHours.WeekdayText[0])
	}

	if placeID != resp.PlaceID {
		t.Errorf("expected %+v, was %+v", placeID, resp.PlaceID)
	}

	authorName := "Rachel Lewis"
	if authorName != resp.Reviews[0].AuthorName {
		t.Errorf("expected %+v, was %+v", authorName, resp.Reviews[0].AuthorName)
	}

	authorURL := "https://plus.google.com/114299517944848975298"
	if authorURL != resp.Reviews[0].AuthorURL {
		t.Errorf("expected %+v, was %+v", authorURL, resp.Reviews[0].AuthorURL)
	}

	language := "en"
	if language != resp.Reviews[0].Language {
		t.Errorf("expected %+v, was %+v", language, resp.Reviews[0].Language)
	}

	rating := 3
	if rating != resp.Reviews[0].Rating {
		t.Errorf("expected %+v, was %+v", rating, resp.Reviews[0].Rating)
	}

	text := "Overall disappointing. This is the second time i've been there and my experience was... Nothing to nibble on for 45 mins and then the bread came. My first entree was the marron which I thought was tasteless - perhaps others would say delicate? but there you go. The XO sea was fantastic. I chose the  vegetarian main dish which was all about the texture which was great but nothing at all outstanding about the dish. My husband and daughter chose the duck for their main course it was the smallest main course i've ever seen - their faces were priceless when it arrived!. Snow egg was beautiful but the granita on the bottom had some solid chunks of hard ice. The service was quite good..."
	if text != resp.Reviews[0].Text {
		t.Errorf("expected %+v, was %+v", text, resp.Reviews[0].Text)
	}

	time := 1441848853
	if time != resp.Reviews[0].Time {
		t.Errorf("expected %+v, was %+v", time, resp.Reviews[0].Time)
	}

}

func TestPlacePhoto(t *testing.T) {
	photoReference := "ThisIsNotAPhotoReference"
	expectedQuery := "key=AIzaNotReallyAnAPIKey&maxheight=400&photoreference=ThisIsNotAPhotoReference"

	server := mockServerForQuery(expectedQuery, 200, "An Image?")
	defer server.s.Close()

	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.s.URL

	r := &PlacePhotoRequest{
		PhotoReference: photoReference,
		MaxHeight:      400,
	}

	_, err := c.PlacePhoto(context.Background(), r)

	if err != nil {
		t.Errorf("Unexpected error in constructing request URL: %+v", err)
	}

	if server.successful != 1 {
		t.Errorf("Got URL(s) %v, want %s", server.failed, expectedQuery)
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
	         "place_id" : "ChIJLdqTiaj51WsRv4Mkbq2qQEU",
	         "reference" : "CnRmAAAAJJuaK6n6aI7imGz2zcqHpBanTQcafAIyja-5pGX6q67WDRT4DJ8M6HcjfxRCbOM-7RAw10sU9l-lZktErhP4mVmavboCyI_QG8iAHNjBPlqYcfFYjJLUE4gtrYvYhx1VGG88wYBbQXXAH4hcGQc3-xIQyNcdcFc9rmijjlL5g1U4KxoUYxqZLWwPfDWy1hkU0DqTUbAm26k",
	         "types" : [ "school", "point_of_interest", "establishment" ]
	      }
	   ],
	   "status" : "OK"
  }`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey))
	c.baseURL = server.URL
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

	placeID := "ChIJLdqTiaj51WsRv4Mkbq2qQEU"
	if placeID != result.PlaceID {
		t.Errorf("expected %+v, was %+v", placeID, result.PlaceID)
	}
}
