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
	"strconv"

	"golang.org/x/net/context"
)

var placesTextSearchAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/textsearch/json",
	acceptsClientID: true,
}

// var placePhotosAPI = &apiConfig{
// 	host:            "https://maps.googleapis.com",
// 	path:            "/maps/api/place/photo",
// 	acceptsClientID: true,
// }

// TextSearch issues the Places API Text Search request and retrieves the Response
func (c *Client) TextSearch(ctx context.Context, r *TextSearchRequest) (PlacesSearchResponse, error) {

	if r.Query == "" && r.PageToken == "" {
		return PlacesSearchResponse{}, errors.New("maps: Query and PageToken both missing")
	}

	if r.Location != nil && r.Radius == 0 {
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

	if r.Radius != 0 {
		q.Set("radius", fmt.Sprint(r.Radius))
	}

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	if r.MinPrice != "" {
		q.Set("minprice", string(r.MinPrice))
	}

	if r.MaxPrice != "" {
		q.Set("maxprice", string(r.MaxPrice))
	}

	if r.OpenNow {
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
	Radius uint
	// Language specifies the language in which to return results. Optional.
	Language string
	// minprice restricts results to only those places within the specified price level. Valid values are in the range from 0 (most affordable) to 4 (most expensive), inclusive.
	MinPrice PriceLevel
	// maxprice restricts results to only those places within the specified price level. Valid values are in the range from 0 (most affordable) to 4 (most expensive), inclusive.
	MaxPrice PriceLevel
	// OpenNow returns only those places that are open for business at the time the query is sent. Places that do not specify opening hours in the Google Places database will not be returned if you include this parameter in your query.
	OpenNow bool
	// PageToken returns the next 20 results from a previously run search. Setting a PageToken parameter will execute a search with the same parameters used previously — all parameters other than PageToken will be ignored.
	PageToken string
}

// PlacesSearchResponse is the response to a Places API Search request.
type PlacesSearchResponse struct {
	// Results is the Place results for the search query
	Results []PlacesSearchResult
	// HTMLAttributions contain a set of attributions about this listing which must be displayed to the user.
	HTMLAttributions []string
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
	// AltIDs — An array of zero, one or more alternative place IDs for the place, with a scope related to each alternative ID.
	AltIDs []AltID `json:"alt_ids"`
	// price_level is the price level of the place, on a scale of 0 to 4.
	PriceLevel int `json:"price_level"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `json:"vicinity"`
}

// AltID is the alternative place IDs for a place.
type AltID struct{}

var placeDetailsAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/details/json",
	acceptsClientID: true,
}

// PlaceDetails issues the Places API Place Details request and retrieves the response
func (c *Client) PlaceDetails(ctx context.Context, r *PlaceDetailsRequest) (PlaceDetailsResponse, error) {

	if r.PlaceID == "" {
		return PlaceDetailsResponse{}, errors.New("maps: PlaceID missing")
	}

	var response struct {
		Result           PlaceDetailsResult `json:"result"`
		HTMLAttributions []string           `json:"html_attributions"`
		commonResponse
	}

	if err := c.getJSON(ctx, placeDetailsAPI, r, &response); err != nil {
		return PlaceDetailsResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return PlaceDetailsResponse{}, err
	}

	return PlaceDetailsResponse{response.Result, response.HTMLAttributions}, nil
}

func (r *PlaceDetailsRequest) params() url.Values {
	q := make(url.Values)

	q.Set("placeid", r.PlaceID)

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	return q
}

// PlaceDetailsRequest is the functional options struct for PlaceDetails
type PlaceDetailsRequest struct {
	// PlaceID is a textual identifier that uniquely identifies a place, returned from a Place Search.
	PlaceID string
	// Language is the language code, indicating in which language the results should be returned, if possible.
	Language string
}

// PlaceDetailsResponse is the response to a Places API Place Detail request.
type PlaceDetailsResponse struct {
	// Results is the Place results for the search query
	Result PlaceDetailsResult
	// HTMLAttributions contain a set of attributions about this listing which must be displayed to the user.
	HTMLAttributions []string
}

// PlaceDetailsResult is an individual Places API Place Details result
type PlaceDetailsResult struct {
	// AddressComponents is an array of separate address components used to compose a given address.
	AddressComponents []AddressComponent `json:"address_components"`
	// FormattedAddress is the human-readable address of this place
	FormattedAddress string `json:"formatted_address"`
	// FormattedPhoneNumber contains the place's phone number in its local format. For example, the formatted_phone_number for Google's Sydney, Australia office is (02) 9374 4000.
	FormattedPhoneNumber string `json:"formatted_phone_number"`
	// InternationalPhoneNumber contains the place's phone number in international format. International format includes the country code, and is prefixed with the plus (+) sign. For example, the international_phone_number for Google's Sydney, Australia office is +61 2 9374 4000.
	InternationalPhoneNumber string `json:"international_phone_number"`
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
	// AltIDs — An array of zero, one or more alternative place IDs for the place, with a scope related to each alternative ID.
	AltIDs []AltID `json:"alt_ids"`
	// price_level is the price level of the place, on a scale of 0 to 4.
	PriceLevel int `json:"price_level"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `json:"vicinity"`
	// permanently_closed is a boolean flag indicating whether the place has permanently shut down (value true). If the place is not permanently closed, the flag is absent from the response.
	PermanentlyClosed bool `json:"permanently_closed"`
	// Reviews is an array of up to five reviews. If a language parameter was specified in the Place Details request, the Places Service will bias the results to prefer reviews written in that language.
	Reviews []PlaceReview `json:"reviews"`
	// UTCOffset contains the number of minutes this place’s current timezone is offset from UTC. For example, for places in Sydney, Australia during daylight saving time this would be 660 (+11 hours from UTC), and for places in California outside of daylight saving time this would be -480 (-8 hours from UTC).
	UTCOffset int `json:"utc_offset"`
	// Website lists the authoritative website for this place, such as a business' homepage.
	Website string `json:"website"`
	// url contains the URL of the official Google page for this place. This will be the establishment's Google+ page if the Google+ page exists, otherwise it will be the Google-owned page that contains the best available information about the place. Applications must link to or embed this page on any screen that shows detailed results about the place to the user.
	URL string `json:"url"`
}

// PlaceReview is a review of a Place
type PlaceReview struct {
	// Aspects contains a collection of AspectRatings, each of which provides a rating of a single attribute of the establishment. The first in the collection is considered the primary aspect.
	Aspects []PlaceReviewAspect `json:"aspects"`
	// AuthorName the name of the user who submitted the review. Anonymous reviews are attributed to "A Google user".
	AuthorName string `json:"author_name"`
	// author_url the URL to the users Google+ profile, if available.
	AuthorURL string `json:"author_url"`
	// Language an IETF language code indicating the language used in the user's review. This field contains the main language tag only, and not the secondary tag indicating country or region.
	Language string `json:"language"`
	// Rating the user's overall rating for this place. This is a whole number, ranging from 1 to 5.
	Rating int `json:"rating"`
	// Text is the user's review. When reviewing a location with Google Places, text reviews are considered optional. Therefore, this field may by empty. Note that this field may include simple HTML markup.
	Text string `json:"text"`
	// Time the time that the review was submitted, measured in the number of seconds since since midnight, January 1, 1970 UTC.
	Time int64 `json:"time"`
}

// PlaceReviewAspect provides a rating of a single attribute of the establishment.
type PlaceReviewAspect struct {
	// Rating is the user's rating for this particular aspect, from 0 to 3.
	Rating int `json:"rating"`
	// Type is the name of the aspect that is being rated. The following types are supported: appeal, atmosphere, decor, facilities, food, overall, quality and service.
	Type string `json:"type"`
}

var placesQueryAutocompleteAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/queryautocomplete/json",
	acceptsClientID: true,
}

// QueryAutocomplete issues the Places API Query Autocomplete request and retrieves the response
func (c *Client) QueryAutocomplete(ctx context.Context, r *QueryAutocompleteRequest) (QueryAutocompleteResponse, error) {

	if r.Input == "" {
		return QueryAutocompleteResponse{}, errors.New("maps: Input missing")
	}

	var response struct {
		Predictions []QueryAutocompletePrediction `json:"predictions"`
		commonResponse
	}

	if err := c.getJSON(ctx, placesQueryAutocompleteAPI, r, &response); err != nil {
		return QueryAutocompleteResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return QueryAutocompleteResponse{}, err
	}

	return QueryAutocompleteResponse{response.Predictions}, nil
}

func (r *QueryAutocompleteRequest) params() url.Values {
	q := make(url.Values)

	q.Set("input", r.Input)

	if r.Offset > 0 {
		q.Set("offset", strconv.FormatUint(uint64(r.Offset), 10))
	}

	if r.Location != nil {
		q.Set("location", r.Location.String())
	}

	if r.Radius > 0 {
		q.Set("radius", strconv.FormatUint(uint64(r.Radius), 10))
	}

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	return q
}

// QueryAutocompleteRequest is the functional options struct for Query Autocomplete
type QueryAutocompleteRequest struct {
	// Input is the text string on which to search. The Places service will return candidate matches based on this string and order results based on their perceived relevance.
	Input string
	// Offset is the character position in the input term at which the service uses text for predictions. For example, if the input is 'Googl' and the completion point is 3, the service will match on 'Goo'. The offset should generally be set to the position of the text caret. If no offset is supplied, the service will use the entire term.
	Offset uint
	// Location is the point around which you wish to retrieve place information.
	Location *LatLng
	// Radius is the distance (in meters) within which to return place results. Note that setting a radius biases results to the indicated area, but may not fully restrict results to the specified area.
	Radius uint
	// Language is the language in which to return results.
	Language string
}

// QueryAutocompleteResponse is a response to a Query Autocomplete request.
type QueryAutocompleteResponse struct {
	Predictions []QueryAutocompletePrediction
}

// QueryAutocompletePrediction represents a single Query Autocomplete result returned from the Google Places API Web Service.
type QueryAutocompletePrediction struct {
	// Description of the matched prediction.
	Description string `json:"description"`
	// PlaceID of the Place
	PlaceID string `json:"place_id"`
	// Types is an array indicating the type of the address component.
	Types []string `json:"types"`
	// MatchedSubstring describes the location of the entered term in the prediction result text, so that the term can be highlighted if desired.
	MatchedSubstrings []QueryAutocompleteMatchedSubstring `json:"matched_substrings"`
	// Terms contains an array of terms identifying each section of the returned description (a section of the description is generally terminated with a comma).
	Terms []QueryAutocompleteTermOffset `json:"terms"`
}

// QueryAutocompleteMatchedSubstring describes the location of the entered term in the prediction result text, so that the term can be highlighted if desired.
type QueryAutocompleteMatchedSubstring struct {
	// Length describes the length of the matched substring.
	Length int `json:"length"`
	// Offset defines the start position of the matched substring.
	Offset int `json:"offset"`
}

// QueryAutocompleteTermOffset identifies each section of the returned description (a section of the description is generally terminated with a comma).
type QueryAutocompleteTermOffset struct {
	// Value is the text of the matched term.
	Value string `json:"value"`
	// Offset defines the start position of this term in the description, measured in Unicode characters.
	Offset int `json:"offset"`
}
