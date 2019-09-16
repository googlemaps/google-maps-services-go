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
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"

	// Included for image/jpeg's decoder
	_ "image/jpeg"
)

var placesNearbySearchAPI = &apiConfig{
	host:             "https://maps.googleapis.com",
	path:             "/maps/api/place/nearbysearch/json",
	acceptsClientID:  true,
	acceptsSignature: false,
}

// NearbySearch lets you search for places within a specified area. You can refine
// your search request by supplying keywords or specifying the type of place you are
// searching for.
func (c *Client) NearbySearch(ctx context.Context, r *NearbySearchRequest) (PlacesSearchResponse, error) {

	if r.PageToken == "" {
		if r.Location == nil {
			return PlacesSearchResponse{}, errors.New("maps: Location and PageToken both missing")
		}

		// Radius is required, unless rank by distance, in which case it isn't allowed.

		if r.Radius == 0 && r.RankBy != RankByDistance {
			return PlacesSearchResponse{}, errors.New("maps: Radius and PageToken both missing")
		}

		if r.Radius > 0 && r.RankBy == RankByDistance {
			return PlacesSearchResponse{}, errors.New("maps: Radius specified with RankByDistance")
		}

		if r.RankBy == RankByDistance && r.Keyword == "" && r.Name == "" && r.Type == "" {
			return PlacesSearchResponse{}, errors.New("maps: RankBy=distance and Keyword, Name and Type are missing")
		}
	}

	var response struct {
		Results          []PlacesSearchResult `json:"results,omitempty"`
		HTMLAttributions []string             `json:"html_attributions,omitempty"`
		NextPageToken    string               `json:"next_page_token,omitempty"`
		commonResponse
	}

	if err := c.getJSON(ctx, placesNearbySearchAPI, r, &response); err != nil {
		return PlacesSearchResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return PlacesSearchResponse{}, err
	}

	return PlacesSearchResponse{response.Results, response.HTMLAttributions, response.NextPageToken}, nil

}

func (r *NearbySearchRequest) params() url.Values {
	q := make(url.Values)

	if r.Location != nil {
		q.Set("location", r.Location.String())
	}

	if r.Radius != 0 {
		q.Set("radius", fmt.Sprint(r.Radius))
	}

	if r.Keyword != "" {
		q.Set("keyword", r.Keyword)
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

	if r.Name != "" {
		q.Set("name", r.Name)
	}

	if r.OpenNow {
		q.Set("opennow", "true")
	}

	if r.RankBy != "" {
		q.Set("rankby", string(r.RankBy))
	}

	if r.Type != "" {
		q.Set("type", string(r.Type))
	}

	if r.PageToken != "" {
		q.Set("pagetoken", r.PageToken)
	}

	return q
}

// NearbySearchRequest is the functional options struct for NearbySearch
type NearbySearchRequest struct {
	// Location is the latitude/longitude around which to retrieve place information.
	// If you specify a location parameter, you must also specify a radius parameter.
	Location *LatLng
	// Radius defines the distance (in meters) within which to bias place results. The
	//maximum allowed radius is 50,000 meters. Results inside of this region will be
	//ranked higher than results outside of the search circle; however, prominent
	// results from outside of the search radius may be included.
	Radius uint
	// Keyword is a term to be matched against all content that Google has indexed for
	// this place, including but not limited to name, type, and address, as well as
	// customer reviews and other third-party content.
	Keyword string
	// Language specifies the language in which to return results. Optional.
	Language string
	// MinPrice restricts results to only those places within the specified price level.
	// Valid values are in the range from 0 (most affordable) to 4 (most expensive),
	// inclusive.
	MinPrice PriceLevel
	// MaxPrice restricts results to only those places within the specified price level.
	// Valid values are in the range from 0 (most affordable) to 4 (most expensive),
	// inclusive.
	MaxPrice PriceLevel
	// Name is one or more terms to be matched against the names of places, separated
	// with a space character.
	Name string
	// OpenNow returns only those places that are open for business at the time the
	// query is sent. Places that do not specify opening hours in the Google Places
	// database will not be returned if you include this parameter in your query.
	OpenNow bool
	// RankBy specifies the order in which results are listed.
	RankBy
	// Type restricts the results to places matching the specified type.
	Type PlaceType
	// PageToken returns the next 20 results from a previously run search. Setting a
	// PageToken parameter will execute a search with the same parameters used
	// previously — all parameters other than PageToken will be ignored.
	PageToken string
}

var placesTextSearchAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/textsearch/json",
	acceptsClientID: true,
}

// TextSearch issues the Places API Text Search request and retrieves the Response
func (c *Client) TextSearch(ctx context.Context, r *TextSearchRequest) (PlacesSearchResponse, error) {

	if r.Query == "" && r.PageToken == "" && r.Type == "" {
		return PlacesSearchResponse{}, errors.New("maps: Query, PageToken and Type are all missing")
	}

	if r.Location != nil && r.Radius == 0 {
		return PlacesSearchResponse{}, errors.New("maps: Radius missing, required with Location")
	}

	var response struct {
		Results          []PlacesSearchResult `json:"results,omitempty"`
		HTMLAttributions []string             `json:"html_attributions,omitempty"`
		NextPageToken    string               `json:"next_page_token,omitempty"`
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

	if r.Type != "" {
		q.Set("type", string(r.Type))
	}

	if r.PageToken != "" {
		q.Set("pagetoken", r.PageToken)
	}

	if r.Region != "" {
		q.Set("region", r.Region)
	}

	return q
}

// TextSearchRequest is the functional options struct for TextSearch
type TextSearchRequest struct {
	// Query is the text string on which to search, for example: "restaurant". The
	// Google Places service will return candidate matches based on this string and
	// order the results based on their perceived relevance.
	Query string
	// Location is the latitude/longitude around which to retrieve place information. If
	// you specify a location parameter, you must also specify a radius parameter.
	Location *LatLng
	// Radius defines the distance (in meters) within which to bias place results. The
	// maximum allowed radius is 50,000 meters. Results inside of this region will be
	// ranked higher than results outside of the search circle; however, prominent
	// results from outside of the search radius may be included.
	Radius uint
	// Language specifies the language in which to return results. Optional.
	Language string
	// MinPrice restricts results to only those places within the specified price level.
	// Valid values are in the range from 0 (most affordable) to 4 (most expensive),
	// inclusive.
	MinPrice PriceLevel
	// MaxPrice restricts results to only those places within the specified price level.
	// Valid values are in the range from 0 (most affordable) to 4 (most expensive),
	// inclusive.
	MaxPrice PriceLevel
	// OpenNow returns only those places that are open for business at the time the
	// query is sent. Places that do not specify opening hours in the Google Places
	// database will not be returned if you include this parameter in your query.
	OpenNow bool
	// Type restricts the results to places matching the specified type.
	Type PlaceType
	// PageToken returns the next 20 results from a previously run search. Setting a
	// PageToken parameter will execute a search with the same parameters used
	// previously — all parameters other than PageToken will be ignored.
	PageToken string
	// The region code, specified as a ccTLD (country code top-level domain) two-character
	// value. Most ccTLD codes are identical to ISO 3166-1 codes, with some exceptions.
	// This parameter will only influence, not fully restrict, search results. If more
	// relevant results exist outside of the specified region, they may be included. When
	// this parameter is used, the country name is omitted from the resulting formatted_address
	// for results in the specified region.
	Region string
}

// PlacesSearchResponse is the response to a Places API Search request.
type PlacesSearchResponse struct {
	// Results is the Place results for the search query
	Results []PlacesSearchResult
	// HTMLAttributions contain a set of attributions about this listing which must be
	// displayed to the user.
	HTMLAttributions []string
	// NextPageToken contains a token that can be used to return up to 20 additional
	// results.
	NextPageToken string
}

// PlacesSearchResult is an individual Places API search result
type PlacesSearchResult struct {
	// FormattedAddress is the human-readable address of this place
	FormattedAddress string `json:"formatted_address,omitempty"`
	// Geometry contains geometry information about the result, generally including the
	// location (geocode) of the place and (optionally) the viewport identifying its
	// general area of coverage.
	Geometry AddressGeometry `json:"geometry,omitempty"`
	// Name contains the human-readable name for the returned result. For establishment
	// results, this is usually the business name.
	Name string `json:"name,omitempty"`
	// Icon contains the URL of a recommended icon which may be displayed to the user
	// when indicating this result.
	Icon string `json:"icon,omitempty"`
	// PlaceID is a textual identifier that uniquely identifies a place.
	PlaceID string `json:"place_id,omitempty"`
	// Scope indicates the scope of the PlaceID.
	Scope string `json:"scope,omitempty"`
	// Rating contains the place's rating, from 1.0 to 5.0, based on aggregated user
	// reviews.
	Rating float32 `json:"rating,omitempty"`
	// UserRatingsTotal contains total number of the place's ratings
	UserRatingsTotal int `json:"user_ratings_total,omitempty"`
	// Types contains an array of feature types describing the given result.
	Types []string `json:"types,omitempty"`
	// OpeningHours may contain whether the place is open now or not.
	OpeningHours *OpeningHours `json:"opening_hours,omitempty"`
	// Photos is an array of photo objects, each containing a reference to an image.
	Photos []Photo `json:"photos,omitempty"`
	// AltIDs — An array of zero, one or more alternative place IDs for the place, with
	// a scope related to each alternative ID.
	AltIDs []AltID `json:"alt_ids,omitempty"`
	// PriceLevel is the price level of the place, on a scale of 0 to 4.
	PriceLevel int `json:"price_level,omitempty"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `json:"vicinity,omitempty"`
	// PermanentlyClosed is a boolean flag indicating whether the place has permanently
	// shut down.
	PermanentlyClosed bool `json:"permanently_closed,omitempty"`
	// BusinessStatus is a string indicating the operational status of the
	// place, if it is a business.
	BusinessStatus string `json:"business_status,omitempty"`
	// ID is an identifier.
	ID string `json:"id,omitempty"`
}

// AltID is the alternative place IDs for a place.
type AltID struct {
	// PlaceID is the APP scoped Place ID that you received when you initially created
	// this Place, before it was given a Google wide Place ID.
	PlaceID string `json:"place_id,omitempty"`
	// Scope is the scope of this alternative place ID. It will always be APP,
	// indicating that the alternative place ID is recognised by your application only.
	Scope string `json:"scope,omitempty"`
}

var placeDetailsAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/details/json",
	acceptsClientID: true,
}

// PlaceDetails issues the Places API Place Details request and retrieves the response
func (c *Client) PlaceDetails(ctx context.Context, r *PlaceDetailsRequest) (PlaceDetailsResult, error) {

	if r.PlaceID == "" {
		return PlaceDetailsResult{}, errors.New("maps: PlaceID missing")
	}

	var response struct {
		Result           PlaceDetailsResult `json:"result,omitempty"`
		HTMLAttributions []string           `json:"html_attributions,omitempty"`
		commonResponse
	}

	if err := c.getJSON(ctx, placeDetailsAPI, r, &response); err != nil {
		return PlaceDetailsResult{}, err
	}

	if err := response.StatusError(); err != nil {
		return PlaceDetailsResult{}, err
	}

	response.Result.HTMLAttributions = response.HTMLAttributions
	return response.Result, nil
}

func (r *PlaceDetailsRequest) params() url.Values {
	q := make(url.Values)

	q.Set("placeid", r.PlaceID)

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	if len(r.Fields) > 0 {
		q.Set("fields", strings.Join(placeDetailsFieldMasksAsStringArray(r.Fields), ","))
	}

	if st := uuid.UUID(r.SessionToken).String(); st != "00000000-0000-0000-0000-000000000000" {
		q.Set("sessiontoken", st)
	}

	if r.Region != "" {
		q.Set("region", r.Region)
	}

	return q
}

// PlaceDetailsRequest is the functional options struct for PlaceDetails
type PlaceDetailsRequest struct {
	// PlaceID is a textual identifier that uniquely identifies a place, returned from a
	// Place Search.
	PlaceID string
	// Language is the language code, indicating in which language the results should be
	// returned, if possible.
	Language string
	// Fields allows you to select which parts of the returned details structure
	// should be filled in. For more detail, please see the following URL:
	// https://cloud.google.com/maps-platform/user-guide/product-changes/#places
	Fields []PlaceDetailsFieldMask
	// SessionToken is a token that marks this request as part of a Place Autocomplete
	// Session. Optional.
	SessionToken PlaceAutocompleteSessionToken
	// Region is the region code, specified as a ccTLD (country code top-level domain)
	// two-character value. Most ccTLD codes are identical to ISO 3166-1 codes, with
	// some exceptions. This parameter will only influence, not fully restrict, results.
	Region string
}

// PlaceDetailsResult is an individual Places API Place Details result
type PlaceDetailsResult struct {
	// AddressComponents is an array of separate address components used to compose a
	// given address.
	AddressComponents []AddressComponent `json:"address_components,omitempty"`
	// FormattedAddress is the human-readable address of this place.
	FormattedAddress string `json:"formatted_address,omitempty"`
	// AdrAddress is the address in the "adr" microformat.
	AdrAddress string `json:"adr_address,omitempty"`
	// FormattedPhoneNumber contains the place's phone number in its local format. For
	// example, the formatted_phone_number for Google's Sydney, Australia office is
	// (02) 9374 4000.
	FormattedPhoneNumber string `json:"formatted_phone_number,omitempty"`
	// InternationalPhoneNumber contains the place's phone number in international
	// format. International format includes the country code, and is prefixed with the
	// plus (+) sign. For example, the international_phone_number for Google's Sydney,
	// Australia office is +61 2 9374 4000.
	InternationalPhoneNumber string `json:"international_phone_number,omitempty"`
	// Geometry contains geometry information about the result, generally including the
	// location (geocode) of the place and (optionally) the viewport identifying its
	// general area of coverage.
	Geometry AddressGeometry `json:"geometry,omitempty"`
	// Name contains the human-readable name for the returned result. For establishment
	// results, this is usually the business name.
	Name string `json:"name,omitempty"`
	// Icon contains the URL of a recommended icon which may be displayed to the user
	// when indicating this result.
	Icon string `json:"icon,omitempty"`
	// PlaceID is a textual identifier that uniquely identifies a place.
	PlaceID string `json:"place_id,omitempty"`
	// Scope indicates the scope of the PlaceID.
	Scope string `json:"scope,omitempty"`
	// Rating contains the place's rating, from 1.0 to 5.0, based on aggregated user
	// reviews.
	Rating float32 `json:"rating,omitempty"`
	// UserRatingsTotal contains total number of the place's ratings
	UserRatingsTotal int `json:"user_ratings_total,omitempty"`
	// Types contains an array of feature types describing the given result.
	Types []string `json:"types,omitempty"`
	// OpeningHours may contain whether the place is open now or not.
	OpeningHours *OpeningHours `json:"opening_hours,omitempty"`
	// Photos is an array of photo objects, each containing a reference to an image.
	Photos []Photo `json:"photos,omitempty"`
	// AltIDs — An array of zero, one or more alternative place IDs for the place, with
	// a scope related to each alternative ID.
	AltIDs []AltID `json:"alt_ids,omitempty"`
	// PriceLevel is the price level of the place, on a scale of 0 to 4.
	PriceLevel int `json:"price_level,omitempty"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `json:"vicinity,omitempty"`
	// PermanentlyClosed is a boolean flag indicating whether the place has permanently
	// shut down (value true). If the place is not permanently closed, the flag is
	// absent from the response.
	PermanentlyClosed bool `json:"permanently_closed,omitempty"`
	// BusinessStatus is a string indicating the operational status of the
	// place, if it is a business.
	BusinessStatus string `json:"business_status,omitempty"`
	// Reviews is an array of up to five reviews. If a language parameter was specified
	// in the Place Details request, the Places Service will bias the results to prefer
	// reviews written in that language.
	Reviews []PlaceReview `json:"reviews,omitempty"`
	// UTCOffset contains the number of minutes this place’s current timezone is offset
	// from UTC. For example, for places in Sydney, Australia during daylight saving
	// time this would be 660 (+11 hours from UTC), and for places in California outside
	// of daylight saving time this would be -480 (-8 hours from UTC).
	UTCOffset *int `json:"utc_offset,omitempty"`
	// Website lists the authoritative website for this place, such as a business'
	// homepage.
	Website string `json:"website,omitempty"`
	// URL contains the URL of the official Google page for this place. This will be the
	// establishment's Google+ page if the Google+ page exists, otherwise it will be the
	// Google-owned page that contains the best available information about the place.
	// Applications must link to or embed this page on any screen that shows detailed
	// results about the place to the user.
	URL string `json:"url,omitempty"`
	// HTMLAttributions contain a set of attributions about this listing which must be
	// displayed to the user.
	HTMLAttributions []string `json:"html_attributions,omitempty"`
}

// PlaceReview is a review of a Place
type PlaceReview struct {
	// Aspects contains a collection of AspectRatings, each of which provides a rating
	// of a single attribute of the establishment. The first in the collection is
	// considered the primary aspect.
	Aspects []PlaceReviewAspect `json:"aspects,omitempty"`
	// AuthorName the name of the user who submitted the review. Anonymous reviews are
	// attributed to "A Google user".
	AuthorName string `json:"author_name,omitempty"`
	// AuthorURL the URL to the user's Google+ profile, if available.
	AuthorURL string `json:"author_url,omitempty"`
	// AuthorPhoto the Google+ profile photo url of the user who submitted the review, if available.
	AuthorProfilePhoto string `json:"profile_photo_url"`
	// Language an IETF language code indicating the language used in the user's review.
	// This field contains the main language tag only, and not the secondary tag
	// indicating country or region.
	Language string `json:"language,omitempty"`
	// Rating the user's overall rating for this place. This is a whole number, ranging
	// from 1 to 5.
	Rating int `json:"rating,omitempty"`
	// Text is the user's review. When reviewing a location with Google Places, text
	// reviews are considered optional. Therefore, this field may by empty. Note that
	// this field may include simple HTML markup.
	Text string `json:"text,omitempty"`
	// Time the time that the review was submitted, measured in the number of seconds
	// since since midnight, January 1, 1970 UTC.
	Time int `json:"time,omitempty"` // TODO(samthor): convert this to a real time.Time
}

// PlaceReviewAspect provides a rating of a single attribute of the establishment.
type PlaceReviewAspect struct {
	// Rating is the user's rating for this particular aspect, from 0 to 3.
	Rating int `json:"rating"`
	// Type is the name of the aspect that is being rated. The following types are
	// supported: appeal, atmosphere, decor, facilities, food, overall, quality and
	// service.
	Type string `json:"type,omitempty"`
}

var placesQueryAutocompleteAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/queryautocomplete/json",
	acceptsClientID: true,
}

// QueryAutocomplete issues the Places API Query Autocomplete request and retrieves
// the response
func (c *Client) QueryAutocomplete(ctx context.Context, r *QueryAutocompleteRequest) (AutocompleteResponse, error) {

	if r.Input == "" {
		return AutocompleteResponse{}, errors.New("maps: Input missing")
	}

	var response struct {
		Predictions []AutocompletePrediction `json:"predictions,omitempty"`
		commonResponse
	}

	if err := c.getJSON(ctx, placesQueryAutocompleteAPI, r, &response); err != nil {
		return AutocompleteResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return AutocompleteResponse{}, err
	}

	return AutocompleteResponse{response.Predictions}, nil
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
	// Input is the text string on which to search. The Places service will return
	// candidate matches based on this string and order results based on their perceived
	// relevance.
	Input string
	// Offset is the character position in the input term at which the service uses text
	// for predictions. For example, if the input is 'Googl' and the completion point is
	// 3, the service will match on 'Goo'. The offset should generally be set to the
	// position of the text caret. If no offset is supplied, the service will use the
	// entire term.
	Offset uint
	// Location is the point around which you wish to retrieve place information.
	Location *LatLng
	// Radius is the distance (in meters) within which to return place results. Note
	// that setting a radius biases results to the indicated area, but may not fully
	// restrict results to the specified area.
	Radius uint
	// Language is the language in which to return results.
	Language string
}

// AutocompleteResponse is a response to a Query Autocomplete request.
type AutocompleteResponse struct {
	Predictions []AutocompletePrediction `json:"predictions"`
}

// AutocompletePrediction represents a single Query Autocomplete result returned from
// the Google Places API Web Service.
type AutocompletePrediction struct {
	// Description of the matched prediction.
	Description string `json:"description,omitempty"`
	// DistanceMeters is the straight-line distance from the prediction to the
	// Origin if Origin was passed in the Query
	DistanceMeters int `json:"distance_meters,omitempty"`
	// PlaceID is the ID of the Place
	PlaceID string `json:"place_id,omitempty"`
	// Types is an array indicating the type of the address component.
	Types []string `json:"types,omitempty"`
	// MatchedSubstring describes the location of the entered term in the prediction
	// result text, so that the term can be highlighted if desired.
	MatchedSubstrings []AutocompleteMatchedSubstring `json:"matched_substrings,omitempty"`
	// Terms contains an array of terms identifying each section of the returned
	// description (a section of the description is generally terminated with a comma).
	Terms []AutocompleteTermOffset `json:"terms,omitempty"`
	// StructuredFormatting contains the main and secondary text of a prediction
	StructuredFormatting AutocompleteStructuredFormatting `json:"structured_formatting,omitempty"`
}

// AutocompleteMatchedSubstring describes the location of the entered term in the
// prediction result text, so that the term can be highlighted if desired.
type AutocompleteMatchedSubstring struct {
	// Length describes the length of the matched substring.
	Length int `json:"length"`
	// Offset defines the start position of the matched substring.
	Offset int `json:"offset"`
}

// AutocompleteTermOffset identifies each section of the returned description (a
// section of the description is generally terminated with a comma).
type AutocompleteTermOffset struct {
	// Value is the text of the matched term.
	Value string `json:"value,omitempty"`
	// Offset defines the start position of this term in the description, measured in
	// Unicode characters.
	Offset int `json:"offset"`
}

// AutocompleteStructuredFormatting contains the main and secondary text of an
// autocomplete prediction
type AutocompleteStructuredFormatting struct {
	MainText                  string                         `json:"main_text,omitempty"`
	MainTextMatchedSubstrings []AutocompleteMatchedSubstring `json:"main_text_matched_substrings,omitempty"`
	SecondaryText             string                         `json:"secondary_text,omitempty"`
}

var placesPlaceAutocompleteAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/autocomplete/json",
	acceptsClientID: true,
}

// PlaceAutocomplete issues the Places API Place Autocomplete request and retrieves
// the response
func (c *Client) PlaceAutocomplete(ctx context.Context, r *PlaceAutocompleteRequest) (AutocompleteResponse, error) {

	if r.Input == "" {
		return AutocompleteResponse{}, errors.New("maps: Input missing")
	}

	var response struct {
		Predictions []AutocompletePrediction `json:"predictions,omitempty"`
		commonResponse
	}

	if err := c.getJSON(ctx, placesPlaceAutocompleteAPI, r, &response); err != nil {
		return AutocompleteResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return AutocompleteResponse{}, err
	}

	return AutocompleteResponse{response.Predictions}, nil
}

func (r *PlaceAutocompleteRequest) params() url.Values {
	q := make(url.Values)

	q.Set("input", r.Input)

	if st := uuid.UUID(r.SessionToken).String(); st != "00000000-0000-0000-0000-000000000000" {
		q.Set("sessiontoken", st)
	}

	if r.Offset > 0 {
		q.Set("offset", strconv.FormatUint(uint64(r.Offset), 10))
	}

	if r.Location != nil {
		q.Set("location", r.Location.String())
	}

	if r.Origin != nil {
		q.Set("origin", r.Origin.String())
	}

	if r.Radius > 0 {
		q.Set("radius", strconv.FormatUint(uint64(r.Radius), 10))
	}

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	if r.Types != "" {
		q.Set("types", string(r.Types))
	}

	if r.StrictBounds {
		q.Set("strictbounds", "true")
	}

	var cf []string
	for c, f := range r.Components {
		fc := make([]string, len(f))
		for i, v := range f {
			fc[i] = string(c) + ":" + v
		}
		cf = append(cf, strings.Join(fc, "|"))
	}
	if len(cf) > 0 {
		q.Set("components", strings.Join(cf, "|"))
	}

	return q
}

// PlaceAutocompleteSessionToken is a session token for Place Autocomplete.
type PlaceAutocompleteSessionToken uuid.UUID

// NewPlaceAutocompleteSessionToken constructs a new Place Autocomplete session token.
func NewPlaceAutocompleteSessionToken() PlaceAutocompleteSessionToken {
	return PlaceAutocompleteSessionToken(uuid.New())
}

// PlaceAutocompleteRequest is the functional options struct for Place Autocomplete
type PlaceAutocompleteRequest struct {
	// Input is the text string on which to search. The Places service will return
	// candidate matches based on this string and order results based on their perceived
	// relevance.
	Input string
	// Offset is the character position in the input term at which the service uses text
	// for predictions. For example, if the input is 'Googl' and the completion point is
	// 3, the service will match on 'Goo'. The offset should generally be set to the
	// position of the text caret. If no offset is supplied, the service will use the
	// entire term.
	Offset uint
	// Location is the point around which you wish to retrieve place information.
	Location *LatLng
	// Origin is the point from which to calculate the straight-line distance to the
	// destination (returned as distance_meters).
	Origin *LatLng
	// Radius is the distance (in meters) within which to return place results. Note
	// that setting a radius biases results to the indicated area, but may not fully
	// restrict results to the specified area.
	Radius uint
	// Language is the language in which to return results.
	Language string
	// Type restricts the results to places matching the specified type.
	Types AutocompletePlaceType
	// Components is a grouping of places to which you would like to restrict your
	// results. Currently, you can use components to filter by country.
	Components map[Component][]string
	// StrictBounds return only those places that are strictly within the region defined
	// by location and radius.
	StrictBounds bool
	// SessionToken is a token that means you will get charged by autocomplete session
	// instead of by character for Autocomplete
	SessionToken PlaceAutocompleteSessionToken
}

var placesPhotoAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/photo",
	acceptsClientID: true,
}

// PlacePhoto issues the Places API Photo request and retrieves the response
func (c *Client) PlacePhoto(ctx context.Context, r *PlacePhotoRequest) (PlacePhotoResponse, error) {

	if r.PhotoReference == "" {
		return PlacePhotoResponse{}, errors.New("maps: PhotoReference missing")
	}

	if r.MaxHeight == 0 && r.MaxWidth == 0 {
		return PlacePhotoResponse{}, errors.New("maps: both MaxHeight & MaxWidth missing")
	}

	resp, err := c.getBinary(ctx, placesPhotoAPI, r)
	if err != nil {
		return PlacePhotoResponse{}, err
	}

	if resp.statusCode == http.StatusForbidden {
		return PlacePhotoResponse{}, errors.New("maps: request exceeds your available quota")
	}

	return PlacePhotoResponse{resp.contentType, resp.data}, nil
}

func (r *PlacePhotoRequest) params() url.Values {
	q := make(url.Values)

	q.Set("photoreference", r.PhotoReference)

	if r.MaxHeight > 0 {
		q.Set("maxheight", strconv.FormatUint(uint64(r.MaxHeight), 10))
	}

	if r.MaxWidth > 0 {
		q.Set("maxwidth", strconv.FormatUint(uint64(r.MaxWidth), 10))
	}

	return q
}

// PlacePhotoRequest is the functional options struct for Places Photo API
type PlacePhotoRequest struct {
	// PhotoReference is a string used to identify the photo when you perform a Photo
	// request.
	PhotoReference string
	// MaxHeight is the maximum height of the image. One of MaxHeight and MaxWidth is
	// required.
	MaxHeight uint
	// MaxWidth is the maximum width of the image. One of MaxHeight and MaxWidth is
	// required.
	MaxWidth uint
}

// PlacePhotoResponse is a response to the Place Photo request
type PlacePhotoResponse struct {
	// ContentType is the server reported type of the Image.
	ContentType string
	// Data is the server returned image data. You must close this after you are
	// finished.
	Data io.ReadCloser
}

// Image will read and close  response.Data and return it as an image.
func (resp *PlacePhotoResponse) Image() (image.Image, error) {
	defer resp.Data.Close()
	if resp.ContentType != "image/jpeg" {
		return nil, errors.New("Image of unknown format: " + resp.ContentType)
	}
	img, _, err := image.Decode(resp.Data)
	return img, err
}

// FindPlaceFromTextInputType is the different types of inputs.
type FindPlaceFromTextInputType string

// The types of FindPlaceFromText Input Types.
const (
	FindPlaceFromTextInputTypeTextQuery   = FindPlaceFromTextInputType("textquery")
	FindPlaceFromTextInputTypePhoneNumber = FindPlaceFromTextInputType("phonenumber")
)

// FindPlaceFromTextLocationBiasType is the type of location bias for this request
type FindPlaceFromTextLocationBiasType string

// The types of FindPlaceFromTextLocationBiasType
const (
	FindPlaceFromTextLocationBiasIP          = FindPlaceFromTextLocationBiasType("ipbias")
	FindPlaceFromTextLocationBiasPoint       = FindPlaceFromTextLocationBiasType("point")
	FindPlaceFromTextLocationBiasCircular    = FindPlaceFromTextLocationBiasType("circle")
	FindPlaceFromTextLocationBiasRectangular = FindPlaceFromTextLocationBiasType("rectangle")
)

// ParseFindPlaceFromTextLocationBiasType will parse a string to a FindPlaceFromTextLocationBiasType
func ParseFindPlaceFromTextLocationBiasType(locationBias string) (FindPlaceFromTextLocationBiasType, error) {
	t := strings.ToLower(locationBias)
	switch t {
	case "ipbias":
		return FindPlaceFromTextLocationBiasIP, nil
	case "point":
		return FindPlaceFromTextLocationBiasPoint, nil
	case "circle":
		return FindPlaceFromTextLocationBiasCircular, nil
	case "rectangle":
		return FindPlaceFromTextLocationBiasRectangular, nil
	}
	return FindPlaceFromTextLocationBiasType(""), fmt.Errorf("Unknown FindPlaceFromTextLocationBiasType \"%v\"", locationBias)
}

// FindPlaceFromTextRequest is the options struct for Find Place From Text API
type FindPlaceFromTextRequest struct {
	// The text input specifying which place to search for (for example, a name,
	// address, or phone number). Required.
	Input string

	// The type of input. Required.
	InputType FindPlaceFromTextInputType

	// Fields allows you to select which parts of the returned details structure
	// should be filled in.
	Fields []PlaceSearchFieldMask

	// LocationBias is the type of location bias to apply to this request
	LocationBias FindPlaceFromTextLocationBiasType

	// LocationBiasPoint is the point for LocationBias type Point
	LocationBiasPoint *LatLng

	// LocationBiasCenter is the center for LocationBias type Circle
	LocationBiasCenter *LatLng

	// LocationBiasRadius is the radius for LocationBias type Circle
	LocationBiasRadius int

	// LocationBiasSouthWest is the South West boundary for LocationBias type Rectangle
	LocationBiasSouthWest *LatLng

	// LocationBiasSouthWest is the North East boundary for LocationBias type Rectangle
	LocationBiasNorthEast *LatLng
}

func (r *FindPlaceFromTextRequest) params() url.Values {
	q := make(url.Values)

	q.Set("input", r.Input)

	q.Set("inputtype", string(r.InputType))

	if len(r.Fields) > 0 {
		q.Set("fields", strings.Join(placeSearchFieldMasksAsStringArray(r.Fields), ","))
	}

	if r.LocationBias != "" {
		switch r.LocationBias {
		case FindPlaceFromTextLocationBiasIP:
			q.Set("locationbias", "ipbias")
		case FindPlaceFromTextLocationBiasPoint:
			q.Set("locationbias", fmt.Sprintf("point:%s", r.LocationBiasPoint.String()))
		case FindPlaceFromTextLocationBiasCircular:
			q.Set("locationbias", fmt.Sprintf("circle:%d@%s", r.LocationBiasRadius, r.LocationBiasCenter.String()))
		case FindPlaceFromTextLocationBiasRectangular:
			q.Set("locationbias", fmt.Sprintf("rectangle:%s|%s", r.LocationBiasSouthWest.String(), r.LocationBiasNorthEast.String()))
		}
	}

	return q
}

// FindPlaceFromTextResponse is a response to the Find Place From Text request
type FindPlaceFromTextResponse struct {
	Candidates       []PlacesSearchResult
	HTMLAttributions []string
}

var findPlaceFromTextAPI = &apiConfig{
	host:            "https://maps.googleapis.com",
	path:            "/maps/api/place/findplacefromtext/json",
	acceptsClientID: false,
}

// FindPlaceFromText takes a text input, and returns a place. The text input
// can be any kind of Places data, for example, a name, address, or phone number.
func (c *Client) FindPlaceFromText(ctx context.Context, r *FindPlaceFromTextRequest) (FindPlaceFromTextResponse, error) {

	if r.Input == "" {
		return FindPlaceFromTextResponse{}, errors.New("maps: Input required")
	}

	if r.InputType == "" {
		return FindPlaceFromTextResponse{}, errors.New("maps: InputType required")
	}

	if r.LocationBias != "" {
		switch r.LocationBias {
		case FindPlaceFromTextLocationBiasPoint:
			if r.LocationBiasPoint == nil {
				return FindPlaceFromTextResponse{}, errors.New("maps: LocationBiasPoint required when LocationBias set to FindPlaceFromTextLocationBiasPoint")
			}
		case FindPlaceFromTextLocationBiasCircular:
			if r.LocationBiasCenter == nil || r.LocationBiasRadius == 0 {
				return FindPlaceFromTextResponse{}, errors.New("maps: LocationBiasCenter and LocationBiasRadius required when LocationBias set to FindPlaceFromTextLocationBiasCircle")
			}
		case FindPlaceFromTextLocationBiasRectangular:
			if r.LocationBiasSouthWest == nil || r.LocationBiasNorthEast == nil {
				return FindPlaceFromTextResponse{}, errors.New("maps: LocationBiasSouthWest and LocationBiasNorthEast required when LocationBias set to FindPlaceFromTextLocationBiasRectangle")
			}
		}
	}

	var response struct {
		Candidates       []PlacesSearchResult `json:"candidates,omitempty"`
		HTMLAttributions []string             `json:"html_attributions,omitempty"`
		commonResponse
	}

	if err := c.getJSON(ctx, findPlaceFromTextAPI, r, &response); err != nil {
		return FindPlaceFromTextResponse{}, err
	}

	if err := response.StatusError(); err != nil {
		return FindPlaceFromTextResponse{}, err
	}

	return FindPlaceFromTextResponse{response.Candidates, response.HTMLAttributions}, nil
}
