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

// More information about Google Distance Matrix API is available on
// https://developers.google.com/maps/documentation/distancematrix/

package maps

import (
	"fmt"
	"strings"
	"time"
)

// Mode is for specifying travel mode.
type Mode string

// Avoid is for specifying routes that avoid certain features.
type Avoid string

// Units specifies which units system to return human readable results in.
type Units string

// TransitMode is for specifying a transit mode for a request
type TransitMode string

// TransitRoutingPreference biases which routes are returned
type TransitRoutingPreference string

// Travel mode preferences.
const (
	TravelModeDriving   = Mode("driving")
	TravelModeWalking   = Mode("walking")
	TravelModeBicycling = Mode("bicycling")
	TravelModeTransit   = Mode("transit")
)

// Features to avoid.
const (
	AvoidTolls    = Avoid("tolls")
	AvoidHighways = Avoid("highways")
	AvoidFerries  = Avoid("ferries")
)

// Units to use on human readable distances.
const (
	UnitsMetric   = Units("metric")
	UnitsImperial = Units("imperial")
)

// Transit mode of directions or distance matrix request.
const (
	TransitModeBus    = TransitMode("bus")
	TransitModeSubway = TransitMode("subway")
	TransitModeTrain  = TransitMode("train")
	TransitModeTram   = TransitMode("tram")
	TransitModeRail   = TransitMode("rail")
)

// Transit Routing preferences for transit mode requests
const (
	TransitRoutingPreferenceLessWalking    = TransitRoutingPreference("less_walking")
	TransitRoutingPreferenceFewerTransfers = TransitRoutingPreference("fewer_transfers")
)

// Distance is the API representation for a distance between two points.
type Distance struct {
	// HumanReadable is the human friendly distance. This is rounded and in an
	// appropriate unit for the request. The units can be overriden with a request
	// parameter.
	HumanReadable string `json:"text"`
	// Meters is the numeric distance, always in meters. This is intended to be used
	// only in algorithmic situations, e.g. sorting results by some user specified
	// metric.
	Meters int `json:"value"`
}

// TrafficModel specifies traffic prediction model when requesting future directions.
type TrafficModel string

// Traffic prediction model when requesting future directions.
const (
	TrafficModelBestGuess   = TrafficModel("best_guess")
	TrafficModelOptimistic  = TrafficModel("optimistic")
	TrafficModelPessimistic = TrafficModel("pessimistic")
)

// PriceLevel is the Price Levels for Places API
type PriceLevel string

// Price Levels for the Places API
const (
	PriceLevelFree          = PriceLevel("0")
	PriceLevelInexpensive   = PriceLevel("1")
	PriceLevelModerate      = PriceLevel("2")
	PriceLevelExpensive     = PriceLevel("3")
	PriceLevelVeryExpensive = PriceLevel("4")
)

// OpeningHours describes the opening hours for a Place Details result.
type OpeningHours struct {
	// OpenNow is a boolean value indicating if the place is open at the current time.
	// Please note, this field will be null if it isn't present in the response.
	OpenNow *bool `json:"open_now,omitempty"`
	// Periods is an array of opening periods covering seven days, starting from Sunday,
	// in chronological order.
	Periods []OpeningHoursPeriod `json:"periods,omitempty"`
	// weekdayText is an array of seven strings representing the formatted opening hours
	// for each day of the week, for example "Monday: 8:30 am – 5:30 pm".
	WeekdayText []string `json:"weekday_text,omitempty"`
	// PermanentlyClosed indicates that the place has permanently shut down. Please
	// note, this field will be null if it isn't present in the response.
	PermanentlyClosed *bool `json:"permanently_closed,omitempty"`
}

// OpeningHoursPeriod is a single OpeningHours day describing when the place opens
// and closes.
type OpeningHoursPeriod struct {
	// Open is when the place opens.
	Open OpeningHoursOpenClose `json:"open"`
	// Close is when the place closes.
	Close OpeningHoursOpenClose `json:"close"`
}

// OpeningHoursOpenClose describes when the place is open.
type OpeningHoursOpenClose struct {
	// Day is a number from 0–6, corresponding to the days of the week, starting on
	// Sunday. For example, 2 means Tuesday.
	Day time.Weekday `json:"day"`
	// Time contains a time of day in 24-hour hhmm format. Values are in the range
	// 0000–2359. The time will be reported in the place’s time zone.
	Time string `json:"time"`
}

// Photo describes a photo available with a Search Result.
type Photo struct {
	// PhotoReference is used to identify the photo when you perform a Photo request.
	PhotoReference string `json:"photo_reference"`
	// Height is the maximum height of the image.
	Height int `json:"height"`
	// Width is the maximum width of the image.
	Width int `json:"width"`
	// htmlAttributions contains any required attributions.
	HTMLAttributions []string `json:"html_attributions"`
}

// Component specifies a key for the parts of a structured address. See
// https://developers.google.com/maps/documentation/geocoding/intro#ComponentFiltering
// for more detail.
type Component string

const (
	// ComponentRoute matches long or short name of a route
	ComponentRoute = Component("route")
	// ComponentLocality matches against both locality and sublocality types
	ComponentLocality = Component("locality")
	// ComponentAdministrativeArea matches all the administrative_area levels
	ComponentAdministrativeArea = Component("administrative_area")
	// ComponentPostalCode matches postal_code and postal_code_prefix
	ComponentPostalCode = Component("postal_code")
	// ComponentCountry matches a country name or a two letter ISO 3166-1 country code
	ComponentCountry = Component("country")
)

// RankBy specifies the order in which results are listed.
type RankBy string

// RankBy options for Places Search.
const (
	RankByProminence = RankBy("prominence")
	RankByDistance   = RankBy("distance")
)

// PlaceType restricts Place API search to the results to places matching the
// specified type.
type PlaceType string

// Place Types for the Places API.
const (
	PlaceTypeAccounting            = PlaceType("accounting")
	PlaceTypeAirport               = PlaceType("airport")
	PlaceTypeAmusementPark         = PlaceType("amusement_park")
	PlaceTypeAquarium              = PlaceType("aquarium")
	PlaceTypeArtGallery            = PlaceType("art_gallery")
	PlaceTypeAtm                   = PlaceType("atm")
	PlaceTypeBakery                = PlaceType("bakery")
	PlaceTypeBank                  = PlaceType("bank")
	PlaceTypeBar                   = PlaceType("bar")
	PlaceTypeBeautySalon           = PlaceType("beauty_salon")
	PlaceTypeBicycleStore          = PlaceType("bicycle_store")
	PlaceTypeBookStore             = PlaceType("book_store")
	PlaceTypeBowlingAlley          = PlaceType("bowling_alley")
	PlaceTypeBusStation            = PlaceType("bus_station")
	PlaceTypeCafe                  = PlaceType("cafe")
	PlaceTypeCampground            = PlaceType("campground")
	PlaceTypeCarDealer             = PlaceType("car_dealer")
	PlaceTypeCarRental             = PlaceType("car_rental")
	PlaceTypeCarRepair             = PlaceType("car_repair")
	PlaceTypeCarWash               = PlaceType("car_wash")
	PlaceTypeCasino                = PlaceType("casino")
	PlaceTypeCemetery              = PlaceType("cemetery")
	PlaceTypeChurch                = PlaceType("church")
	PlaceTypeCityHall              = PlaceType("city_hall")
	PlaceTypeClothingStore         = PlaceType("clothing_store")
	PlaceTypeConvenienceStore      = PlaceType("convenience_store")
	PlaceTypeCourthouse            = PlaceType("courthouse")
	PlaceTypeDentist               = PlaceType("dentist")
	PlaceTypeDepartmentStore       = PlaceType("department_store")
	PlaceTypeDoctor                = PlaceType("doctor")
	PlaceTypeElectrician           = PlaceType("electrician")
	PlaceTypeElectronicsStore      = PlaceType("electronics_store")
	PlaceTypeEmbassy               = PlaceType("embassy")
	PlaceTypeFireStation           = PlaceType("fire_station")
	PlaceTypeFlorist               = PlaceType("florist")
	PlaceTypeFuneralHome           = PlaceType("funeral_home")
	PlaceTypeFurnitureStore        = PlaceType("furniture_store")
	PlaceTypeGasStation            = PlaceType("gas_station")
	PlaceTypeGym                   = PlaceType("gym")
	PlaceTypeHairCare              = PlaceType("hair_care")
	PlaceTypeHardwareStore         = PlaceType("hardware_store")
	PlaceTypeHinduTemple           = PlaceType("hindu_temple")
	PlaceTypeHomeGoodsStore        = PlaceType("home_goods_store")
	PlaceTypeHospital              = PlaceType("hospital")
	PlaceTypeInsuranceAgency       = PlaceType("insurance_agency")
	PlaceTypeJewelryStore          = PlaceType("jewelry_store")
	PlaceTypeLaundry               = PlaceType("laundry")
	PlaceTypeLawyer                = PlaceType("lawyer")
	PlaceTypeLibrary               = PlaceType("library")
	PlaceTypeLiquorStore           = PlaceType("liquor_store")
	PlaceTypeLocalGovernmentOffice = PlaceType("local_government_office")
	PlaceTypeLocksmith             = PlaceType("locksmith")
	PlaceTypeLodging               = PlaceType("lodging")
	PlaceTypeMealDelivery          = PlaceType("meal_delivery")
	PlaceTypeMealTakeaway          = PlaceType("meal_takeaway")
	PlaceTypeMosque                = PlaceType("mosque")
	PlaceTypeMovieRental           = PlaceType("movie_rental")
	PlaceTypeMovieTheater          = PlaceType("movie_theater")
	PlaceTypeMovingCompany         = PlaceType("moving_company")
	PlaceTypeMuseum                = PlaceType("museum")
	PlaceTypeNightClub             = PlaceType("night_club")
	PlaceTypePainter               = PlaceType("painter")
	PlaceTypePark                  = PlaceType("park")
	PlaceTypeParking               = PlaceType("parking")
	PlaceTypePetStore              = PlaceType("pet_store")
	PlaceTypePharmacy              = PlaceType("pharmacy")
	PlaceTypePhysiotherapist       = PlaceType("physiotherapist")
	PlaceTypePlumber               = PlaceType("plumber")
	PlaceTypePolice                = PlaceType("police")
	PlaceTypePostOffice            = PlaceType("post_office")
	PlaceTypeRealEstateAgency      = PlaceType("real_estate_agency")
	PlaceTypeRestaurant            = PlaceType("restaurant")
	PlaceTypeRoofingContractor     = PlaceType("roofing_contractor")
	PlaceTypeRvPark                = PlaceType("rv_park")
	PlaceTypeSchool                = PlaceType("school")
	PlaceTypeShoeStore             = PlaceType("shoe_store")
	PlaceTypeShoppingMall          = PlaceType("shopping_mall")
	PlaceTypeSpa                   = PlaceType("spa")
	PlaceTypeStadium               = PlaceType("stadium")
	PlaceTypeStorage               = PlaceType("storage")
	PlaceTypeStore                 = PlaceType("store")
	PlaceTypeSubwayStation         = PlaceType("subway_station")
	PlaceTypeSupermarket           = PlaceType("supermarket")
	PlaceTypeSynagogue             = PlaceType("synagogue")
	PlaceTypeTaxiStand             = PlaceType("taxi_stand")
	PlaceTypeTrainStation          = PlaceType("train_station")
	PlaceTypeTravelAgency          = PlaceType("travel_agency")
	PlaceTypeUniversity            = PlaceType("university")
	PlaceTypeVeterinaryCare        = PlaceType("veterinary_care")
	PlaceTypeZoo                   = PlaceType("zoo")
)

// ParsePlaceType will parse a string representation of a PlaceType.
func ParsePlaceType(placeType string) (PlaceType, error) {
	switch strings.ToLower(placeType) {
	case "accounting":
		return PlaceTypeAccounting, nil
	case "airport":
		return PlaceTypeAirport, nil
	case "amusement_park":
		return PlaceTypeAmusementPark, nil
	case "aquarium":
		return PlaceTypeAquarium, nil
	case "art_gallery":
		return PlaceTypeArtGallery, nil
	case "atm":
		return PlaceTypeAtm, nil
	case "bakery":
		return PlaceTypeBakery, nil
	case "bank":
		return PlaceTypeBank, nil
	case "bar":
		return PlaceTypeBar, nil
	case "beauty_salon":
		return PlaceTypeBeautySalon, nil
	case "bicycle_store":
		return PlaceTypeBicycleStore, nil
	case "book_store":
		return PlaceTypeBookStore, nil
	case "bowling_alley":
		return PlaceTypeBowlingAlley, nil
	case "bus_station":
		return PlaceTypeBusStation, nil
	case "cafe":
		return PlaceTypeCafe, nil
	case "campground":
		return PlaceTypeCampground, nil
	case "car_dealer":
		return PlaceTypeCarDealer, nil
	case "car_rental":
		return PlaceTypeCarRental, nil
	case "car_repair":
		return PlaceTypeCarRepair, nil
	case "car_wash":
		return PlaceTypeCarWash, nil
	case "casino":
		return PlaceTypeCasino, nil
	case "cemetery":
		return PlaceTypeCemetery, nil
	case "church":
		return PlaceTypeChurch, nil
	case "city_hall":
		return PlaceTypeCityHall, nil
	case "clothing_store":
		return PlaceTypeClothingStore, nil
	case "convenience_store":
		return PlaceTypeConvenienceStore, nil
	case "courthouse":
		return PlaceTypeCourthouse, nil
	case "dentist":
		return PlaceTypeDentist, nil
	case "department_store":
		return PlaceTypeDepartmentStore, nil
	case "doctor":
		return PlaceTypeDoctor, nil
	case "electrician":
		return PlaceTypeElectrician, nil
	case "electronics_store":
		return PlaceTypeElectronicsStore, nil
	case "embassy":
		return PlaceTypeEmbassy, nil
	case "fire_station":
		return PlaceTypeFireStation, nil
	case "florist":
		return PlaceTypeFlorist, nil
	case "funeral_home":
		return PlaceTypeFuneralHome, nil
	case "furniture_store":
		return PlaceTypeFurnitureStore, nil
	case "gas_station":
		return PlaceTypeGasStation, nil
	case "gym":
		return PlaceTypeGym, nil
	case "hair_care":
		return PlaceTypeHairCare, nil
	case "hardware_store":
		return PlaceTypeHardwareStore, nil
	case "hindu_temple":
		return PlaceTypeHinduTemple, nil
	case "home_goods_store":
		return PlaceTypeHomeGoodsStore, nil
	case "hospital":
		return PlaceTypeHospital, nil
	case "insurance_agency":
		return PlaceTypeInsuranceAgency, nil
	case "jewelry_store":
		return PlaceTypeJewelryStore, nil
	case "laundry":
		return PlaceTypeLaundry, nil
	case "lawyer":
		return PlaceTypeLawyer, nil
	case "library":
		return PlaceTypeLibrary, nil
	case "liquor_store":
		return PlaceTypeLiquorStore, nil
	case "local_government_office":
		return PlaceTypeLocalGovernmentOffice, nil
	case "locksmith":
		return PlaceTypeLocksmith, nil
	case "lodging":
		return PlaceTypeLodging, nil
	case "meal_delivery":
		return PlaceTypeMealDelivery, nil
	case "meal_takeaway":
		return PlaceTypeMealTakeaway, nil
	case "mosque":
		return PlaceTypeMosque, nil
	case "movie_rental":
		return PlaceTypeMovieRental, nil
	case "movie_theater":
		return PlaceTypeMovieTheater, nil
	case "moving_company":
		return PlaceTypeMovingCompany, nil
	case "museum":
		return PlaceTypeMuseum, nil
	case "night_club":
		return PlaceTypeNightClub, nil
	case "painter":
		return PlaceTypePainter, nil
	case "park":
		return PlaceTypePark, nil
	case "parking":
		return PlaceTypeParking, nil
	case "pet_store":
		return PlaceTypePetStore, nil
	case "pharmacy":
		return PlaceTypePharmacy, nil
	case "physiotherapist":
		return PlaceTypePhysiotherapist, nil
	case "plumber":
		return PlaceTypePlumber, nil
	case "police":
		return PlaceTypePolice, nil
	case "post_office":
		return PlaceTypePostOffice, nil
	case "real_estate_agency":
		return PlaceTypeRealEstateAgency, nil
	case "restaurant":
		return PlaceTypeRestaurant, nil
	case "roofing_contractor":
		return PlaceTypeRoofingContractor, nil
	case "rv_park":
		return PlaceTypeRvPark, nil
	case "school":
		return PlaceTypeSchool, nil
	case "shoe_store":
		return PlaceTypeShoeStore, nil
	case "shopping_mall":
		return PlaceTypeShoppingMall, nil
	case "spa":
		return PlaceTypeSpa, nil
	case "stadium":
		return PlaceTypeStadium, nil
	case "storage":
		return PlaceTypeStorage, nil
	case "store":
		return PlaceTypeStore, nil
	case "subway_station":
		return PlaceTypeSubwayStation, nil
	case "supermarket":
		return PlaceTypeSupermarket, nil
	case "synagogue":
		return PlaceTypeSynagogue, nil
	case "taxi_stand":
		return PlaceTypeTaxiStand, nil
	case "train_station":
		return PlaceTypeTrainStation, nil
	case "travel_agency":
		return PlaceTypeTravelAgency, nil
	case "university":
		return PlaceTypeUniversity, nil
	case "veterinary_care":
		return PlaceTypeVeterinaryCare, nil
	case "zoo":
		return PlaceTypeZoo, nil
	default:
		return PlaceType(""), fmt.Errorf("Unknown PlaceType \"%v\"", placeType)
	}
}

// AutocompletePlaceType restricts Place Autocomplete API to the results to places
// matching the specified type.
type AutocompletePlaceType string

// https://developers.google.com/places/web-service/autocomplete#place_types
const (
	AutocompletePlaceTypeGeocode       = AutocompletePlaceType("geocode")
	AutocompletePlaceTypeAddress       = AutocompletePlaceType("address")
	AutocompletePlaceTypeEstablishment = AutocompletePlaceType("establishment")
	AutocompletePlaceTypeRegions       = AutocompletePlaceType("(regions)")
	AutocompletePlaceTypeCities        = AutocompletePlaceType("(cities)")
)

// ParseAutocompletePlaceType will parse a string representation of a
// AutocompletePlaceTypes.
func ParseAutocompletePlaceType(placeType string) (AutocompletePlaceType, error) {
	switch strings.ToLower(placeType) {
	case "geocode":
		return AutocompletePlaceTypeGeocode, nil
	case "address":
		return AutocompletePlaceTypeAddress, nil
	case "establishment":
		return AutocompletePlaceTypeEstablishment, nil
	case "(regions)":
		return AutocompletePlaceTypeRegions, nil
	case "(cities)":
		return AutocompletePlaceTypeCities, nil
	default:
		return AutocompletePlaceType(""), fmt.Errorf("Unknown AutocompletePlaceType \"%v\"", placeType)
	}
}

// PlaceDetailsFieldMask allows you to specify which fields are to be returned with
// a place details request. Please see the following URL for more detail:
// https://cloud.google.com/maps-platform/user-guide/product-changes/#places
type PlaceDetailsFieldMask string

// The individual Place Details Field Masks.
const (
	PlaceDetailsFieldMaskAddressComponent             = PlaceDetailsFieldMask("address_component")
	PlaceDetailsFieldMaskADRAddress                   = PlaceDetailsFieldMask("adr_address")
	PlaceDetailsFieldMaskAltID                        = PlaceDetailsFieldMask("alt_id")
	PlaceDetailsFieldMaskBusinessStatus               = PlaceDetailsFieldMask("business_status")
	PlaceDetailsFieldMaskFormattedAddress             = PlaceDetailsFieldMask("formatted_address")
	PlaceDetailsFieldMaskFormattedPhoneNumber         = PlaceDetailsFieldMask("formatted_phone_number")
	PlaceDetailsFieldMaskGeometry                     = PlaceDetailsFieldMask("geometry")
	PlaceDetailsFieldMaskGeometryLocation             = PlaceDetailsFieldMask("geometry/location")
	PlaceDetailsFieldMaskGeometryLocationLat          = PlaceDetailsFieldMask("geometry/location/lat")
	PlaceDetailsFieldMaskGeometryLocationLng          = PlaceDetailsFieldMask("geometry/location/lng")
	PlaceDetailsFieldMaskGeometryViewport             = PlaceDetailsFieldMask("geometry/viewport")
	PlaceDetailsFieldMaskGeometryViewportNortheast    = PlaceDetailsFieldMask("geometry/viewport/northeast")
	PlaceDetailsFieldMaskGeometryViewportNortheastLat = PlaceDetailsFieldMask("geometry/viewport/northeast/lat")
	PlaceDetailsFieldMaskGeometryViewportNortheastLng = PlaceDetailsFieldMask("geometry/viewport/northeast/lng")
	PlaceDetailsFieldMaskGeometryViewportSouthwest    = PlaceDetailsFieldMask("geometry/viewport/southwest")
	PlaceDetailsFieldMaskGeometryViewportSouthwestLat = PlaceDetailsFieldMask("geometry/viewport/southwest/lat")
	PlaceDetailsFieldMaskGeometryViewportSouthwestLng = PlaceDetailsFieldMask("geometry/viewport/southwest/lng")
	PlaceDetailsFieldMaskIcon                         = PlaceDetailsFieldMask("icon")
	PlaceDetailsFieldMaskID                           = PlaceDetailsFieldMask("id")
	PlaceDetailsFieldMaskInternationalPhoneNumber     = PlaceDetailsFieldMask("international_phone_number")
	PlaceDetailsFieldMaskName                         = PlaceDetailsFieldMask("name")
	PlaceDetailsFieldMaskOpeningHours                 = PlaceDetailsFieldMask("opening_hours")
	PlaceDetailsFieldMaskPermanentlyClosed            = PlaceDetailsFieldMask("permanently_closed")
	PlaceDetailsFieldMaskPhotos                       = PlaceDetailsFieldMask("photos")
	PlaceDetailsFieldMaskPlaceID                      = PlaceDetailsFieldMask("place_id")
	PlaceDetailsFieldMaskPriceLevel                   = PlaceDetailsFieldMask("price_level")
	PlaceDetailsFieldMaskRatings                      = PlaceDetailsFieldMask("rating")
	PlaceDetailsFieldMaskUserRatingsTotal             = PlaceDetailsFieldMask("user_ratings_total")
	PlaceDetailsFieldMaskReviews                      = PlaceDetailsFieldMask("reviews")
	PlaceDetailsFieldMaskScope                        = PlaceDetailsFieldMask("scope")
	PlaceDetailsFieldMaskTypes                        = PlaceDetailsFieldMask("types")
	PlaceDetailsFieldMaskURL                          = PlaceDetailsFieldMask("url")
	PlaceDetailsFieldMaskUTCOffset                    = PlaceDetailsFieldMask("utc_offset")
	PlaceDetailsFieldMaskVicinity                     = PlaceDetailsFieldMask("vicinity")
	PlaceDetailsFieldMaskWebsite                      = PlaceDetailsFieldMask("website")
)

// ParsePlaceDetailsFieldMask will parse a string representation of
// PlaceDetailsFieldMask.
func ParsePlaceDetailsFieldMask(placeDetailsFieldMask string) (PlaceDetailsFieldMask, error) {
	switch strings.ToLower(placeDetailsFieldMask) {
	case "address_component":
		return PlaceDetailsFieldMaskAddressComponent, nil
	case "adr_address":
		return PlaceDetailsFieldMaskADRAddress, nil
	case "alt_id":
		return PlaceDetailsFieldMaskAltID, nil
	case "business_status":
		return PlaceDetailsFieldMaskBusinessStatus, nil
	case "formatted_address":
		return PlaceDetailsFieldMaskFormattedAddress, nil
	case "formatted_phone_number":
		return PlaceDetailsFieldMaskFormattedPhoneNumber, nil
	case "geometry":
		return PlaceDetailsFieldMaskGeometry, nil
	case "geometry/location":
		return PlaceDetailsFieldMaskGeometryLocation, nil
	case "geometry/location/lat":
		return PlaceDetailsFieldMaskGeometryLocationLat, nil
	case "geometry/location/lng":
		return PlaceDetailsFieldMaskGeometryLocationLng, nil
	case "geometry/viewport":
		return PlaceDetailsFieldMaskGeometryViewport, nil
	case "geometry/viewport/northeast":
		return PlaceDetailsFieldMaskGeometryViewportNortheast, nil
	case "geometry/viewport/northeast/lat":
		return PlaceDetailsFieldMaskGeometryViewportNortheastLat, nil
	case "geometry/viewport/northeast/lng":
		return PlaceDetailsFieldMaskGeometryViewportNortheastLng, nil
	case "geometry/viewport/southwest":
		return PlaceDetailsFieldMaskGeometryViewportSouthwest, nil
	case "geometry/viewport/southwest/lat":
		return PlaceDetailsFieldMaskGeometryViewportSouthwestLat, nil
	case "geometry/viewport/southwest/lng":
		return PlaceDetailsFieldMaskGeometryViewportSouthwestLng, nil
	case "icon":
		return PlaceDetailsFieldMaskIcon, nil
	case "id":
		return PlaceDetailsFieldMaskID, nil
	case "international_phone_number":
		return PlaceDetailsFieldMaskInternationalPhoneNumber, nil
	case "name":
		return PlaceDetailsFieldMaskName, nil
	case "opening_hours":
		return PlaceDetailsFieldMaskOpeningHours, nil
	case "permanently_closed":
		return PlaceDetailsFieldMaskPermanentlyClosed, nil
	case "photos":
		return PlaceDetailsFieldMaskPhotos, nil
	case "place_id":
		return PlaceDetailsFieldMaskPlaceID, nil
	case "price_level":
		return PlaceDetailsFieldMaskPriceLevel, nil
	case "rating":
		return PlaceDetailsFieldMaskRatings, nil
	case "user_ratings_total":
		return PlaceDetailsFieldMaskUserRatingsTotal, nil
	case "reviews":
		return PlaceDetailsFieldMaskReviews, nil
	case "scope":
		return PlaceDetailsFieldMaskScope, nil
	case "types":
		return PlaceDetailsFieldMaskTypes, nil
	case "url":
		return PlaceDetailsFieldMaskURL, nil
	case "utc_offset":
		return PlaceDetailsFieldMaskUTCOffset, nil
	case "vicinity":
		return PlaceDetailsFieldMaskVicinity, nil
	case "website":
		return PlaceDetailsFieldMaskWebsite, nil
	default:
		return PlaceDetailsFieldMask(""), fmt.Errorf("Unknown PlaceDetailsFieldMask \"%v\"", placeDetailsFieldMask)
	}
}

// fieldsAsStringArray converts []PlaceDetailsFieldMask to []string
func placeDetailsFieldMasksAsStringArray(fields []PlaceDetailsFieldMask) []string {
	var res []string
	for _, el := range fields {
		res = append(res, string(el))
	}
	return res
}

// PlaceSearchFieldMask allows you to specify which fields are to be returned with
// a place search request. Please see the following URL for more detail:
// https://cloud.google.com/maps-platform/user-guide/product-changes/#places
type PlaceSearchFieldMask string

// The individual Place Search Field Masks.
const (
	PlaceSearchFieldMaskAltID                        = PlaceSearchFieldMask("alt_id")
	PlaceSearchFieldMaskBusinessStatus               = PlaceSearchFieldMask("business_status")
	PlaceSearchFieldMaskFormattedAddress             = PlaceSearchFieldMask("formatted_address")
	PlaceSearchFieldMaskGeometry                     = PlaceSearchFieldMask("geometry")
	PlaceSearchFieldMaskGeometryLocation             = PlaceSearchFieldMask("geometry/location")
	PlaceSearchFieldMaskGeometryLocationLat          = PlaceSearchFieldMask("geometry/location/lat")
	PlaceSearchFieldMaskGeometryLocationLng          = PlaceSearchFieldMask("geometry/location/lng")
	PlaceSearchFieldMaskGeometryViewport             = PlaceSearchFieldMask("geometry/viewport")
	PlaceSearchFieldMaskGeometryViewportNortheast    = PlaceSearchFieldMask("geometry/viewport/northeast")
	PlaceSearchFieldMaskGeometryViewportNortheastLat = PlaceSearchFieldMask("geometry/viewport/northeast/lat")
	PlaceSearchFieldMaskGeometryViewportNortheastLng = PlaceSearchFieldMask("geometry/viewport/northeast/lng")
	PlaceSearchFieldMaskGeometryViewportSouthwest    = PlaceSearchFieldMask("geometry/viewport/southwest")
	PlaceSearchFieldMaskGeometryViewportSouthwestLat = PlaceSearchFieldMask("geometry/viewport/southwest/lat")
	PlaceSearchFieldMaskGeometryViewportSouthwestLng = PlaceSearchFieldMask("geometry/viewport/southwest/lng")
	PlaceSearchFieldMaskIcon                         = PlaceSearchFieldMask("icon")
	PlaceSearchFieldMaskID                           = PlaceSearchFieldMask("id")
	PlaceSearchFieldMaskName                         = PlaceSearchFieldMask("name")
	PlaceSearchFieldMaskOpeningHours                 = PlaceSearchFieldMask("opening_hours")
	PlaceSearchFieldMaskOpeningHoursOpenNow          = PlaceSearchFieldMask("opening_hours/open_now")
	PlaceSearchFieldMaskPermanentlyClosed            = PlaceSearchFieldMask("permanently_closed")
	PlaceSearchFieldMaskPhotos                       = PlaceSearchFieldMask("photos")
	PlaceSearchFieldMaskPlaceID                      = PlaceSearchFieldMask("place_id")
	PlaceSearchFieldMaskPriceLevel                   = PlaceSearchFieldMask("price_level")
	PlaceSearchFieldMaskRating                       = PlaceSearchFieldMask("rating")
	PlaceSearchFieldMaskUserRatingsTotal             = PlaceSearchFieldMask("user_ratings_total")
	PlaceSearchFieldMaskReference                    = PlaceSearchFieldMask("reference")
	PlaceSearchFieldMaskTypes                        = PlaceSearchFieldMask("types")
	PlaceSearchFieldMaskVicinity                     = PlaceSearchFieldMask("vicinity")
)

// ParsePlaceSearchFieldMask will parse a string representation of
// PlaceSearchFieldMask.
func ParsePlaceSearchFieldMask(placeSearchFieldMask string) (PlaceSearchFieldMask, error) {
	switch strings.ToLower(placeSearchFieldMask) {
	case "alt_id":
		return PlaceSearchFieldMaskAltID, nil
	case "formatted_address":
		return PlaceSearchFieldMaskFormattedAddress, nil
	case "geometry":
		return PlaceSearchFieldMaskGeometry, nil
	case "geometry/location":
		return PlaceSearchFieldMaskGeometryLocation, nil
	case "geometry/location/lat":
		return PlaceSearchFieldMaskGeometryLocationLat, nil
	case "geometry/location/lng":
		return PlaceSearchFieldMaskGeometryLocationLng, nil
	case "geometry/viewport":
		return PlaceSearchFieldMaskGeometryViewport, nil
	case "geometry/viewport/northeast":
		return PlaceSearchFieldMaskGeometryViewportNortheast, nil
	case "geometry/viewport/northeast/lat":
		return PlaceSearchFieldMaskGeometryViewportNortheastLat, nil
	case "geometry/viewport/northeast/lng":
		return PlaceSearchFieldMaskGeometryViewportNortheastLng, nil
	case "geometry/viewport/southwest":
		return PlaceSearchFieldMaskGeometryViewportSouthwest, nil
	case "geometry/viewport/southwest/lat":
		return PlaceSearchFieldMaskGeometryViewportSouthwestLat, nil
	case "geometry/viewport/southwest/lng":
		return PlaceSearchFieldMaskGeometryViewportSouthwestLng, nil
	case "icon":
		return PlaceSearchFieldMaskIcon, nil
	case "id":
		return PlaceSearchFieldMaskID, nil
	case "name":
		return PlaceSearchFieldMaskName, nil
	case "opening_hours":
		return PlaceSearchFieldMaskOpeningHours, nil
	case "opening_hours/open_now":
		return PlaceSearchFieldMaskOpeningHoursOpenNow, nil
	case "permanently_closed":
		return PlaceSearchFieldMaskPermanentlyClosed, nil
	case "photos":
		return PlaceSearchFieldMaskPhotos, nil
	case "place_id":
		return PlaceSearchFieldMaskPlaceID, nil
	case "price_level":
		return PlaceSearchFieldMaskPriceLevel, nil
	case "rating":
		return PlaceSearchFieldMaskRating, nil
	case "user_ratings_total":
		return PlaceSearchFieldMaskUserRatingsTotal, nil
	case "reference":
		return PlaceSearchFieldMaskReference, nil
	case "types":
		return PlaceSearchFieldMaskTypes, nil
	case "vicinity":
		return PlaceSearchFieldMaskVicinity, nil
	default:
		return PlaceSearchFieldMask(""), fmt.Errorf("Unknown PlaceSearchFieldMask \"%v\"", placeSearchFieldMask)
	}
}

func placeSearchFieldMasksAsStringArray(fields []PlaceSearchFieldMask) []string {
	var res []string
	for _, el := range fields {
		res = append(res, string(el))
	}
	return res
}
