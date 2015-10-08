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
	"errors"
	"fmt"
	"net/url"

	"golang.org/x/net/context"
)

var placesTextSearchAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/textsearch/json",
	acceptsClientID: true,
}

// var placeDetailsAPI = &apiConfig{
// 	host:            "https://maps.googleapis.com",
// 	path:            "/maps/api/place/details/json",
// 	acceptsClientID: true,
// }
//
// var placePhotosAPI = &apiConfig{
// 	host:            "https://maps.googleapis.com",
// 	path:            "/maps/api/place/photo",
// 	acceptsClientID: true,
// }
//
// var placesQueryAutocompleteAPI = &apiConfig{
// 	host:            "https://maps.googleapis.com",
// 	path:            "/maps/api/place/queryautocomplete/json",
// 	acceptsClientID: true,
// }

// TextSearch issues the Places API Text Search request and retrieves the Response
func (c *Client) TextSearch(ctx context.Context, r *TextSearchRequest) (PlacesSearchResponse, error) {
	if r.Location != nil && r.Radius == nil {
		return PlacesSearchResponse{}, errors.New("maps: Radius missing, required with Location")
	}

	var response struct {
		Results          []PlacesSearchResult `json:"results"`
		HTMLAttributions []string             `json:"html_attributions"`
		NextPageToken    string               `json:"next_page_token"`
		commonResponse
	}

	if err := c.getJSON(ctx, placesTextSearchAPI, r, &response); err != nil {
		return PlacesSearchResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return PlacesSearchResponse{}, err
	}

	return PlacesSearchResponse{response.Results, response.HTMLAttributions, response.NextPageToken}, nil
}

func (r *TextSearchRequest) params() url.Values {
	q := make(url.Values)

	q.Set("query", r.Query)

	if r.Location != nil {
		q.Set("location", r.Location.String())
	}

	if r.Radius != nil {
		q.Set("radius", fmt.Sprint(r.Radius))
	}

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	if r.MinPrice != nil {
		q.Set("minprice", fmt.Sprint(*r.MinPrice))
	}

	if r.MaxPrice != nil {
		q.Set("maxprice", fmt.Sprint(*r.MaxPrice))
	}

	if r.OpenNow != nil && *r.OpenNow {
		q.Set("opennow", "true")
	}

	return q
}

// TextSearchRequest is the functional options struct for TextSearch
type TextSearchRequest struct {
	// Query is the text string on which to search, for example: "restaurant". The Google Places service will return candidate matches based on this string and order the results based on their perceived relevance.
	Query string
	// Location is the latitude/longitude around which to retrieve place information. If you specify a location parameter, you must also specify a radius parameter.
	Location *LatLng
	// Radius defines the distance (in meters) within which to bias place results. The maximum allowed radius is 50,000 meters. Results inside of this region will be ranked higher than results outside of the search circle; however, prominent results from outside of the search radius may be included.
	Radius *uint
	// Language specifies the language in which to return results. Optional.
	Language string
	// minprice restricts results to only those places within the specified price level. Valid values are in the range from 0 (most affordable) to 4 (most expensive), inclusive.
	MinPrice *PriceLevel
	// maxprice restricts results to only those places within the specified price level. Valid values are in the range from 0 (most affordable) to 4 (most expensive), inclusive.
	MaxPrice *PriceLevel
	// OpenNow returns only those places that are open for business at the time the query is sent. Places that do not specify opening hours in the Google Places database will not be returned if you include this parameter in your query.
	OpenNow *bool
	// PageToken returns the next 20 results from a previously run search. Setting a PageToken parameter will execute a search with the same parameters used previously — all parameters other than PageToken will be ignored.
	PageToken string
}

// PlacesSearchResponse is the response to a Places API Search request.
type PlacesSearchResponse struct {
	// Results is the Place results for the search query
	Results []PlacesSearchResult `json:"results"`
	// HTMLAttributions contain a set of attributions about this listing which must be displayed to the user.
	HTMLAttributions []string `json:"html_attributions"`
	// NextPageToken contains a token that can be used to return up to 20 additional results.
	NextPageToken string
}

// PlacesSearchResult is an individual Places API search result
type PlacesSearchResult struct {
	// FormattedAddress is the human-readable address of this place
	FormattedAddress string `json:"formatted_address"`
	// geometry contains geometry information about the result, generally including the location (geocode) of the place and (optionally) the viewport identifying its general area of coverage.
	Geometry AddressGeometry `json:"geometry"`
	// Name contains the human-readable name for the returned result. For establishment results, this is usually the business name.
	Name string `json:"name"`
	// Icon contains the URL of a recommended icon which may be displayed to the user when indicating this result.
	Icon string `json:"icon"`
	// PlaceID is a textual identifier that uniquely identifies a place.
	PlaceID string `json:"place_id"`
	// Scope indicates the scope of the PlaceID.
	Scope string `json:"scope"`
	// Rating contains the place's rating, from 1.0 to 5.0, based on aggregated user reviews.
	Rating float32
	// Types contains an array of feature types describing the given result.
	Types []string `json:"types"`
	// OpeningHours may contain whether the place is open now or not.
	OpeningHours *OpeningHours `json:"opening_hours"`
	// Photos is an array of photo objects, each containing a reference to an image.
	Photos []Photo `json:"photos"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `json:"vicinity"`
}
