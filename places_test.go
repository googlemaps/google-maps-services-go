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
