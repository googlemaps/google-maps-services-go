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

// More information about Google Address Validation API is available on
// https://developers.google.com/maps/documentation/address-validation

package maps

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

var addressValidationAPI = &apiConfig{
	host:             "https://addressvalidation.googleapis.com",
	path:             "/v1:validateAddress",
	acceptsClientID:  true,
	acceptsSignature: false,
}

// ValidateAddress makes a Address Validation API request
func (c *Client) ValidateAddress(ctx context.Context, r *AddressValidationRequest) (*AddressValidationResult, error) {
	if r.Address == nil {
		return nil, errors.New("maps: address is missing")
	}

	var response struct {
		AddressValidationResult
		ErrorResponse
	}
	if err := c.getJSON(ctx, addressValidationAPI, r, &response); err != nil {
		return nil, err
	}

	if response.Err != nil {
		return nil, &response.ErrorResponse
	}

	return &response.AddressValidationResult, nil
}

// Represents a postal address, e.g. for postal delivery or payments addresses.
// Given a postal address, a postal service can deliver items to a premise, P.O.
// Box or similar.
// It is not intended to model geographical locations (roads, towns,
// mountains).
//
// In typical usage an address would be created via user input or from importing
// existing data, depending on the type of process.
//
// Advice on address input / editing:
//   - Use an i18n-ready address widget such as
//     https://github.com/google/libaddressinput)
//   - Users should not be presented with UI elements for input or editing of
//     fields outside countries where that field is used.
//
// For more guidance on how to use this schema, please see:
// https://support.google.com/business/answer/6397478
type PostalAddress struct {
	// The schema revision of the `PostalAddress`. This must be set to 0, which is
	// the latest revision.
	//
	// All new revisions **must** be backward compatible with old revisions.
	Revision int32 `json:"revision,omitempty"`
	// Optional. CLDR region code of the country/region of the address. See https://cldr.unicode.org/ and
	// https://www.unicode.org/cldr/charts/30/supplemental/territory_information.html for details.
	// Example "CH" for Switzerland. If the region code is not provided, it will be inferred from the address.
	// For best performance, it is recommended to include the region code if you know it. Having
	// inconsistent or repeated regions can lead to poor performance, for example, if the addressLines
	// already includes the region, do not provide the region code again in this field. Supported regions can
	// be found in the FAQ.
	RegionCode string ` json:"regionCode,omitempty"`
	// Optional. BCP-47 language code of the contents of this address (if
	// known). This is often the UI language of the input form or is expected
	// to match one of the languages used in the address' country/region, or their
	// transliterated equivalents.
	// This can affect formatting in certain countries, but is not critical
	// to the correctness of the data and will never affect any validation or
	// other non-formatting related operations.
	//
	// If this value is not known, it should be omitted (rather than specifying a
	// possibly incorrect default).
	//
	// Examples: "zh-Hant", "ja", "ja-Latn", "en".
	LanguageCode string ` json:"languageCode,omitempty"`
	// Optional. Postal code of the address. Not all countries use or require
	// postal codes to be present, but where they are used, they may trigger
	// additional validation with other parts of the address (e.g. state/zip
	// validation in the U.S.A.).
	PostalCode string ` json:"postalCode,omitempty"`
	// Optional. Additional, country-specific, sorting code. This is not used
	// in most regions. Where it is used, the value is either a string like
	// "CEDEX", optionally followed by a number (e.g. "CEDEX 7"), or just a number
	// alone, representing the "sector code" (Jamaica), "delivery area indicator"
	// (Malawi) or "post office indicator" (e.g. Côte d'Ivoire).
	SortingCode string ` json:"sortingCode,omitempty"`
	// Optional. Highest administrative subdivision which is used for postal
	// addresses of a country or region.
	// For example, this can be a state, a province, an oblast, or a prefecture.
	// Specifically, for Spain this is the province and not the autonomous
	// community (e.g. "Barcelona" and not "Catalonia").
	// Many countries don't use an administrative area in postal addresses. E.g.
	// in Switzerland this should be left unpopulated.
	AdministrativeArea string ` json:"administrativeArea,omitempty"`
	// Optional. Generally refers to the city/town portion of the address.
	// Examples: US city, IT comune, UK post town.
	// In regions of the world where localities are not well defined or do not fit
	// into this structure well, leave locality empty and use address_lines.
	Locality string ` json:"locality,omitempty"`
	// Optional. Sublocality of the address.
	// For example, this can be neighborhoods, boroughs, districts.
	Sublocality string ` json:"sublocality,omitempty"`
	// Unstructured address lines describing the lower levels of an address.
	//
	// Because values in address_lines do not have type information and may
	// sometimes contain multiple values in a single field (e.g.
	// "Austin, TX"), it is important that the line order is clear. The order of
	// address lines should be "envelope order" for the country/region of the
	// address. In places where this can vary (e.g. Japan), address_language is
	// used to make it explicit (e.g. "ja" for large-to-small ordering and
	// "ja-Latn" or "en" for small-to-large). This way, the most specific line of
	// an address can be selected based on the language.
	//
	// The minimum permitted structural representation of an address consists
	// of a region_code with all remaining information placed in the
	// address_lines. It would be possible to format such an address very
	// approximately without geocoding, but no semantic reasoning could be
	// made about any of the address components until it was at least
	// partially resolved.
	//
	// Creating an address only containing a region_code and address_lines, and
	// then geocoding is the recommended way to handle completely unstructured
	// addresses (as opposed to guessing which parts of the address should be
	// localities or administrative areas).
	AddressLines []string ` json:"addressLines,omitempty"`
	// Optional. The recipient at the address.
	// This field may, under certain circumstances, contain multiline information.
	// For example, it might contain "care of" information.
	Recipients []string ` json:"recipients,omitempty"`
	// Optional. The name of the organization at the address.
	Organization string ` json:"organization,omitempty"`
}

// LanguageOptions is the language options for the Address Validation API.
type LanguageOptions struct {
	// Preview: Return a Address in English.
	ReturnEnglishLatinAddress bool `json:"returnEnglishLatinAddress,omitempty"`
}

// AddressValidationRequest is the request format for the Address Validation API.
type AddressValidationRequest struct {
	// Required. The address being validated. Unformatted addresses should be submitted via
	// addressLines.
	// The total length of the fields in this input must not exceed 280 characters.
	// Supported regions can be found [here](https://developers.google.com/maps/documentation/address-validation/coverage).
	// The languageCode value in the input address is reserved for future uses and is ignored today. The
	// validated address result will be populated based on the preferred language for the given address, as
	// identified by the system.
	// The Address Validation API ignores the values in recipients and organization. Any values in
	// those fields will be discarded and not returned. Please do not set them.
	Address *PostalAddress `json:"address,omitempty"`
	// This field must be empty for the first address validation request. If
	// more requests are necessary to fully validate a single address (for
	// example if the changes the user makes after the initial validation need to
	// be re-validated), then each followup request must populate this field with
	// the response_id from the very first response in the validation sequence.
	PreviousResponseId string `json:"previousResponse_id,omitempty"`
	// Optional. Preview: This feature is in Preview (pre-GA). Pre-GA products and features might have
	// limited support, and changes to pre-GA products and features might not be compatible With Other
	// pre-GA versions. Pre-GA Offerings are covered by the Google Maps Platform Service Specific Terms.
	// For more information, see the launch stage descriptions.
	// Enables the Address Validation API to include additional information in the response.
	LanguageOptions LanguageOptions `json:"languageOptions,omitempty"`
}

func (r *AddressValidationRequest) params() url.Values {
	q := make(url.Values)

	return q
}

// AddressValidationResult is the result format for the Address Validation API.
type AddressValidationResult struct {
	// The result of the address validation.
	Result *ValidationResult `json:"result"`
	// The UUID that identifies this response. If the address needs to be re-validated, this UUID must accompany the new request.
	ResponseId string `json:"responseId"`
}

// The result of validating an address.
type ValidationResult struct {
	// Overall verdict flags
	Verdict *Verdict `json:"verdict"`
	// Information about the address itself as opposed to the geocode.
	Address *AddressValidationAddress `json:"address"`
	// Information about the location and place that the address geocoded to.
	Geocode *AddressValidationGeocode ` json:"geocode"`
	// Other information relevant to deliverability. `metadata` is not guaranteed
	// to be fully populated for every address sent to the Address Validation API.
	Metadata *AddressMetadata ` json:"metadata"`
}

// The various granularities that an address or a geocode can have.
// When used to indicate granularity for an *address*, these values indicate
// with how fine a granularity the address identifies a mailing destination.
// For example, an address such as "123 Main Street, Redwood City, CA, 94061"
// identifies a `PREMISE` while something like "Redwood City, CA, 94061"
// identifies a `LOCALITY`. However, if we are unable to find a geocode for
// "123 Main Street" in Redwood City, the geocode returned might be of
// `LOCALITY` granularity even though the address is more granular.
type VerdictGranularity string

const (
	// Default value. This value is unused.
	VerdictGranularityGranularityUnspecified = VerdictGranularity("GRANULARITY_UNSPECIFIED")
	// Below-building level result, such as an apartment.
	VerdictGranularitySubPremise = VerdictGranularity("SUB_PREMISE")
	// Building-level result.
	VerdictGranularityPremise = VerdictGranularity("PREMISE")
	// A geocode that should be very close to the building-level location of
	// the address.
	VerdictGranularityPremiseProximity = VerdictGranularity("PREMISE_PROXIMITY")
	// The address or geocode indicates a block. Only used in regions which
	// have block-level addressing, such as Japan.
	VerdictGranularityBlock = VerdictGranularity("BLOCK")
	// The geocode or address is granular to route, such as a street, road, or
	// highway.
	VerdictGranularityRoute = VerdictGranularity("ROUTE")
	// All other granularities, which are bucketed together since they are not
	// deliverable.
	VerdictGranularityOther = VerdictGranularity("OTHER")
)

// High level overview of the address validation result and geocode.
type Verdict struct {
	// The granularity of the **input** address. This is the result of parsing the
	// input address and does not give any validation signals. For validation
	// signals, refer to `validationGranularity` below.
	//
	// For example, if the input address includes a specific apartment number,
	// then the `inputGranularity` here will be `SUB_PREMISE`. If we cannot match
	// the apartment number in the databases or the apartment number is invalid,
	// the `validationGranularity` will likely be `PREMISE` or below.
	InputGranularity VerdictGranularity `json:"inputGranularity"`
	// The granularity level that the API can fully **validate** the address to.
	// For example, an `validationGranularity` of `PREMISE` indicates all address
	// components at the level of `PREMISE` or more coarse can be validated.
	ValidationGranularity VerdictGranularity `json:"validationGranularity"`
	// Information about the granularity of the `geocode`
	// This can be understood as the semantic meaning of how coarse or fine the
	// geocoded location is.
	//
	// This can differ from the `validationGranularity` above occasionally. For
	// example, our database might record the existence of an apartment number but
	// do not have a precise location for the apartment within a big apartment
	// complex. In that case, the `validationGranularity` will be `SUB_PREMISE`
	// but the `geocodeGranularity` will be `PREMISE`.
	GeocodeGranularity VerdictGranularity `json:"geocodeGranularity"`
	// The address is considered complete if there are no unresolved tokens, no
	// unexpected or missing address components.
	AddressComplete bool `json:"addressComplete"`
	// At least one address component cannot be categorized or validated.
	HasUnconfirmedComponents bool `json:"hasUnconfirmedComponents"`
	// At least one address component was inferred (added) that wasn't in the input
	HasInferredComponents bool `json:"hasInferredComponents"`
	// At least one address component was replaced
	HasReplacedComponents bool `json:"hasReplacedComponents"`
}

// A wrapper for the name of the component.
type AddressValidationComponentName struct {
	// The name text. For example, "5th Avenue" for a street name or "1253" for a
	// street number.
	Text string `json:"text"`
	// The BCP-47 language code. This will not be present if the component name is
	// not associated with a language, such as a street number.
	LanguageCode string `json:"languageCode"`
}

// The different possible values for confirmation levels.
type AddressValidationAddressComponentConfirmationLevel string

const (
	// Default value. This value is unused.
	AddressComponentConfirmationLevelConfirmationLevelUnspecified = AddressValidationAddressComponentConfirmationLevel("CONFIRMATION_LEVEL_UNSPECIFIED")
	// We were able to verify that this component exists and makes sense in the
	// context of the rest of the address.
	AddressComponentConfirmationLevelConfirmed = AddressValidationAddressComponentConfirmationLevel("CONFIRMED")
	// This component could not be confirmed, but it is plausible that it
	// exists. For example, a street number within a known valid range of
	// numbers on a street where specific house numbers are not known.
	AddressComponentConfirmationLevelUnconfirmedButPlausible = AddressValidationAddressComponentConfirmationLevel("UNCONFIRMED_BUT_PLAUSIBLE")
	// This component was not confirmed and is likely to be wrong. For
	// example, a neighborhood that does not fit the rest of the address.
	AddressComponentConfirmationLevelUnconfirmedAndSuspicous = AddressValidationAddressComponentConfirmationLevel("UNCONFIRMED_AND_SUSPICIOUS")
)

// Represents an address component, such as a street, city, or state.
type AddressValidationComponent struct {
	// The name for this component.
	ComponentName *AddressValidationComponentName `json:"componentName"`
	// The type of the address component. See
	// [Table 2: Additional types returned by the Places
	// service](https://developers.google.com/places/web-service/supported_types#table2)
	// for a list of possible types.
	ComponentType string `json:"componentType"`
	// Indicates the level of certainty that we have that the component
	// is correct.
	ConfirmationLevel AddressValidationAddressComponentConfirmationLevel `json:"confirmationLevel"`
	// Indicates that the component was not part of the input, but we
	// inferred it for the address location and believe it should be provided
	// for a complete address.
	Inferred bool `json:"inferred"`
	// Indicates the spelling of the component name was corrected in a minor way,
	// for example by switching two characters that appeared in the wrong order.
	// This indicates a cosmetic change.
	SpellCorrected bool `json:"spellCorrected"`
	// Indicates the name of the component was replaced with a completely
	// different one, for example a wrong postal code being replaced with one that
	// is correct for the address. This is not a cosmetic change, the input
	// component has been changed to a different one.
	Replaced bool `json:"replaced"`
	// Indicates an address component that is not expected to be present in a
	// postal address for the given region. We have retained it only because it
	// was part of the input.
	Unexpected bool `json:"unexpected"`
}

// Details of the post-processed address. Post-processing includes
// correcting misspelled parts of the address, replacing incorrect parts, and
// inferring missing parts.
type AddressValidationAddress struct {
	// The post-processed address, formatted as a single-line address following
	// the address formatting rules of the region where the address is located.
	FormattedAddress string `json:"formattedAddress"`
	// The post-processed address represented as a postal address.
	PostalAddress *PostalAddress `json:"postalAddress"`
	// Unordered list. The individual address components of the formatted and
	// corrected address, along with validation information. This provides
	// information on the validation status of the individual components.
	//
	// Address components are not ordered in a particular way. Do not make any
	// assumptions on the ordering of the address components in the list.
	AddressComponents []*AddressValidationComponent `json:"addressComponents"`
	// The types of components that were expected to be present in a correctly
	// formatted mailing address but were not found in the input AND could
	// not be inferred. Components of this type are not present in
	// `formattedAddress`, `postalAddress`, or `addressComponents`. An
	// example might be `['streetNumber', 'route']` for an input like
	// "Boulder, Colorado, 80301, USA". The list of possible types can be found
	// [here](https://developers.google.com/maps/documentation/geocoding/requests-geocoding#Types).
	MissingComponentTypes []string `json:"missingComponentTypes"`
	// The types of the components that are present in the `address_components`
	// but could not be confirmed to be correct. This field is provided for the
	// sake of convenience: its contents are equivalent to iterating through the
	// `addressComponents` to find the types of all the components where the
	// confirmationLevel is not CONFIRMED or the inferred flag is not set to `true`.
	// The list of possible types can be found
	// [here](https://developers.google.com/maps/documentation/geocoding/requests-geocoding#Types).
	UnconfirmedComponentTypes []string `json:"unconfirmedComponentTypes"`
	// Any tokens in the input that could not be resolved. This might be an
	// input that was not recognized as a valid part of an address (for example
	// in an input like "123235253253 Main St, San Francisco, CA, 94105", the
	// unresolved tokens may look like `["123235253253"]` since that does not
	// look like a valid street number.
	UnresolvedTokens []string `json:"unresolvedTokens"`
}

// Plus code (http://plus.codes) is a location reference with two formats:
// global code defining a 14mx14m (1/8000th of a degree) or smaller rectangle,
// and compound code, replacing the prefix with a reference location.
type AddressValidationPlusCode struct {
	// Place's global (full) code, such as "9FWM33GV+HQ", representing an
	// 1/8000 by 1/8000 degree area (~14 by 14 meters).
	GlobalCode string `json:"globalCode"`
	// Place's compound code, such as "33GV+HQ, Ramberg, Norway", containing
	// the suffix of the global code and replacing the prefix with a formatted
	// name of a reference entity.
	CompoundCode string `json:"compoundCode"`
}

// An object that represents a latitude/longitude pair. This is expressed as a
// pair of doubles to represent degrees latitude and degrees longitude. Unless
// specified otherwise, this must conform to the
// <a href="http://www.unoosa.org/pdf/icg/2012/template/WGS_84.pdf">WGS84
// standard</a>. Values must be within normalized ranges.
type AddressValidationLatLng struct {
	// The latitude in degrees. It must be in the range [-90.0, +90.0].
	Latitude float64 `json:"latitude,omitempty"`
	// The longitude in degrees. It must be in the range [-180.0, +180.0].
	Longitude float64 `json:"longitude,omitempty"`
}

// A latitude-longitude viewport, represented as two diagonally opposite `low`
// and `high` points. A viewport is considered a closed region, i.e. it includes
// its boundary. The latitude bounds must range between -90 to 90 degrees
// inclusive, and the longitude bounds must range between -180 to 180 degrees
// inclusive. Various cases include:
//
//   - If `low` = `high`, the viewport consists of that single point.
//
//   - If `low.longitude` > `high.longitude`, the longitude range is inverted
//     (the viewport crosses the 180 degree longitude line).
//
//   - If `low.longitude` = -180 degrees and `high.longitude` = 180 degrees,
//     the viewport includes all longitudes.
//
//   - If `low.longitude` = 180 degrees and `high.longitude` = -180 degrees,
//     the longitude range is empty.
//
//   - If `low.latitude` > `high.latitude`, the latitude range is empty.
//
// Both `low` and `high` must be populated, and the represented box cannot be
// empty (as specified by the definitions above). An empty viewport will result
// in an error.
//
// For example, this viewport fully encloses New York City:
//
//	{
//	    "low": {
//	        "latitude": 40.477398,
//	        "longitude": -74.259087
//	    },
//	    "high": {
//	        "latitude": 40.91618,
//	        "longitude": -73.70018
//	    }
//	}
type AddressValidationViewport struct {
	// Required. The low point of the viewport.
	Low *AddressValidationLatLng `json:"low,omitempty"`
	// Required. The high point of the viewport.
	High *AddressValidationLatLng `json:"high,omitempty"`
}

// Contains information about the place the input was geocoded to.
type AddressValidationGeocode struct {
	// The geocoded location of the input.
	//
	// Using place IDs is preferred over using addresses,
	// latitude/longitude coordinates, or plus codes. Using coordinates when
	// routing or calculating driving directions will always result in the point
	// being snapped to the road nearest to those coordinates. This may not be a
	// road that will quickly or safely lead to the destination and may not be
	// near an access point to the property. Additionally, when a location is
	// reverse geocoded, there is no guarantee that the returned address will
	// match the original.
	Location *AddressValidationLatLng `json:"location"`
	// The plus code corresponding to the `location`.
	PlusCode *AddressValidationPlusCode `json:"plusCode"`
	// The bounds of the geocoded place.
	Bounds *AddressValidationViewport `json:"bounds"`
	// The size of the geocoded place, in meters. This is another measure of the
	// coarseness of the geocoded location, but in physical size rather than in
	// semantic meaning.
	FeatureSizeMeters float32 `json:"featureSizeMeters"`
	// The PlaceID of the place this input geocodes to.
	//
	// For more information about Place IDs see
	// [here](https://developers.google.com/maps/documentation/places/web-service/place-id).
	PlaceId string `json:"placeId"`
	// The type(s) of place that the input geocoded to. For example,
	// `['locality', 'political']`. The full list of types can be found
	// [here](https://developers.google.com/maps/documentation/geocoding/requests-geocoding#Types).
	PlaceTypes []string `json:"placeTypes"`
}

// The metadata for the address. `metadata` is not guaranteed to be fully
// populated for every address sent to the Address Validation API.
type AddressMetadata struct {
	// Indicates that this is the address of a business.
	// If unset, indicates that the value is unknown.
	Business *bool `json:"business,omitempty"`
	// Indicates that the address of a PO box.
	// If unset, indicates that the value is unknown.
	PoBox *bool `json:"poBox,omitempty"`
	// Indicates that this is the address of a residence.
	// If unset, indicates that the value is unknown.
	Residential *bool `json:"residential,omitempty"`
}

// The Error object represents a general error returned by the API.
type ErrorResponse struct {
	Err *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
		Details []struct {
			Type     string `json:"@type"`
			Reason   string `json:"reason"`
			Domain   string `json:"domain"`
			Metadata struct {
				Service string `json:"service"`
			} `json:"metadata"`
		} `json:"details"`
	} `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("maps: %s - %s", e.Err.Status, e.Err.Message)
}
