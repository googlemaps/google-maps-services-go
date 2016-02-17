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
	// HumanReadable is the human friendly distance. This is rounded and in an appropriate unit for the
	// request. The units can be overriden with a request parameter.
	HumanReadable string `json:"text"`
	// Meters is the numeric distance, always in meters. This is intended to be used only in
	// algorithmic situations, e.g. sorting results by some user specified metric.
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
	// OpenNow is a boolean value indicating if the place is open at the current time. Please note, this field will be null if it isn't present in the response.
	OpenNow *bool `json:"open_now"`
	// Periods is an array of opening periods covering seven days, starting from Sunday, in chronological order.
	Periods []OpeningHoursPeriod `json:"periods"`
	// weekdayText is an array of seven strings representing the formatted opening hours for each day of the week, for example "Monday: 8:30 am – 5:30 pm".
	WeekdayText []string `json:"weekday_text"`
	// PermanentlyClosed indicates that the place has permanently shut down. Please note, this field will be null if it isn't present in the response.
	PermanentlyClosed *bool `json:"permanently_closed"`
}

// OpeningHoursPeriod is a single OpeningHours day describing when the place opens and closes.
type OpeningHoursPeriod struct {
	// Open is when the place opens.
	Open OpeningHoursOpenClose `json:"open"`
	// Close is when the place closes.
	Close OpeningHoursOpenClose `json:"close"`
}

// OpeningHoursOpenClose describes when the place is open.
type OpeningHoursOpenClose struct {
	// Day is a number from 0–6, corresponding to the days of the week, starting on Sunday. For example, 2 means Tuesday.
	Day time.Weekday `json:"day"`
	// Time contains a time of day in 24-hour hhmm format. Values are in the range 0000–2359. The time will be reported in the place’s time zone.
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

// Component specifies a key for the parts of a structured address.
// See https://developers.google.com/maps/documentation/geocoding/intro#ComponentFiltering for more detail.
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

// RankBy specifies the order in which results are listed
type RankBy string

// RankBy options for Places Search
const (
	RankByProminence = RankBy("prominence")
	RankByDistance   = RankBy("distance")
)

// PlaceType restricts Place API search to the results to places matching the specified type.
type PlaceType string

// Warning: DO NOT EDIT PlaceType* - they are code generated.

// Place Types for the Places API
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
	PlaceTypeEstablishment         = PlaceType("establishment")
	PlaceTypeFinance               = PlaceType("finance")
	PlaceTypeFireStation           = PlaceType("fire_station")
	PlaceTypeFlorist               = PlaceType("florist")
	PlaceTypeFood                  = PlaceType("food")
	PlaceTypeFuneralHome           = PlaceType("funeral_home")
	PlaceTypeFurnitureStore        = PlaceType("furniture_store")
	PlaceTypeGasStation            = PlaceType("gas_station")
	PlaceTypeGeneralContractor     = PlaceType("general_contractor")
	PlaceTypeGroceryOrSupermarket  = PlaceType("grocery_or_supermarket")
	PlaceTypeGym                   = PlaceType("gym")
	PlaceTypeHairCare              = PlaceType("hair_care")
	PlaceTypeHardwareStore         = PlaceType("hardware_store")
	PlaceTypeHealth                = PlaceType("health")
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
	PlaceTypePlaceOfWorship        = PlaceType("place_of_worship")
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
	PlaceTypeSynagogue             = PlaceType("synagogue")
	PlaceTypeTaxiStand             = PlaceType("taxi_stand")
	PlaceTypeTrainStation          = PlaceType("train_station")
	PlaceTypeTravelAgency          = PlaceType("travel_agency")
	PlaceTypeUniversity            = PlaceType("university")
	PlaceTypeVeterinaryCare        = PlaceType("veterinary_care")
	PlaceTypeZoo                   = PlaceType("zoo")
)

// Warning: DO NOT EDIT ParsePlaceType() - it is generated code.

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
	case "establishment":
		return PlaceTypeEstablishment, nil
	case "finance":
		return PlaceTypeFinance, nil
	case "fire_station":
		return PlaceTypeFireStation, nil
	case "florist":
		return PlaceTypeFlorist, nil
	case "food":
		return PlaceTypeFood, nil
	case "funeral_home":
		return PlaceTypeFuneralHome, nil
	case "furniture_store":
		return PlaceTypeFurnitureStore, nil
	case "gas_station":
		return PlaceTypeGasStation, nil
	case "general_contractor":
		return PlaceTypeGeneralContractor, nil
	case "grocery_or_supermarket":
		return PlaceTypeGroceryOrSupermarket, nil
	case "gym":
		return PlaceTypeGym, nil
	case "hair_care":
		return PlaceTypeHairCare, nil
	case "hardware_store":
		return PlaceTypeHardwareStore, nil
	case "health":
		return PlaceTypeHealth, nil
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
	case "place_of_worship":
		return PlaceTypePlaceOfWorship, nil
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
		return PlaceType("Unknown PlaceType"), fmt.Errorf("Unknown PlaceType \"%v\"", placeType)
	}
}
