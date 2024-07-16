// Copyright 2024 Google Inc. All Rights Reserved.
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
	"fmt"
)

/**
* An enum representing the relationship in space between the landmark and the target.
*/
type SpatialRelationship string

const (
	 // This is the default relationship when nothing more specific below
    // applies.
    SPATIAL_RELATIONSHIP_NEAR                   SpatialRelationship = "NEAR"
	// The landmark has a spatial geometry and the target is within its
	// bounds.
    SPATIAL_RELATIONSHIP_WITHIN                 SpatialRelationship = "WITHIN"
	// The target is directly adjacent to the landmark or landmark's access
	// point.
    SPATIAL_RELATIONSHIP_BESIDE                 SpatialRelationship = "BESIDE"
	// The target is directly opposite the landmark on the other side of the
	// road.
    SPATIAL_RELATIONSHIP_ACROSS_THE_ROAD        SpatialRelationship = "ACROSS_THE_ROAD"
	// On the same route as the landmark but not besides or across.
    SPATIAL_RELATIONSHIP_DOWN_THE_ROAD          SpatialRelationship = "DOWN_THE_ROAD"
	// Not on the same route as the landmark but a single 'turn' away.
    SPATIAL_RELATIONSHIP_AROUND_THE_CORNER      SpatialRelationship = "AROUND_THE_CORNER"
	// Close to the landmark's structure but further away from its access
	// point.
    SPATIAL_RELATIONSHIP_BEHIND                 SpatialRelationship = "BEHIND"
)

// String method for formatted output
func (sr SpatialRelationship) String() string {
    return string(sr)
}

/**
* An enum representing the relationship in space between the area and the target.
*/
type Containment string

const (
	/**
	* Indicates an unknown containment returned by the server.
	*/
    CONTAINMENT_UNSPECIFIED Containment = "CONTAINMENT_UNSPECIFIED"
	/** The target location is within the area region, close to the center. */
    CONTAINMENT_WITHIN                  Containment = "WITHIN"
	/** The target location is within the area region, close to the edge. */
    CONTAINMENT_OUTSKIRTS               Containment = "OUTSKIRTS"
	/** The target location is outside the area region, but close by. */
    CONTAINMENT_NEAR                    Containment = "NEAR"
)

// String method for formatted output
func (c Containment) String() string {
    return string(c)
}

/**
* Localized variant of a text in a particular language.
*/
type LocalizedText struct {
	// Localized string in the language corresponding to language_code below.
    Text         string `json:"text"`
	// The text's BCP-47 language code, such as "en-US" or "sr-Latn".
    //
    // For more information, see
    // http://www.unicode.org/reports/tr35/#Unicode_locale_identifier.
    LanguageCode string `json:"language_code"`
}

// String method for formatted output
func (lt LocalizedText) String() string {
    return fmt.Sprintf("(text=%s, languageCode=%s)", lt.Text, lt.LanguageCode)
}

// Landmarks that are useful at describing a location.
type Landmark struct {
	// The Place ID of the underlying establishment serving as the landmark.
    // Can be used to resolve more information about the landmark through Place
    // Details or Place Id Lookup.
    PlaceID                     string                   `json:"place_id"`
	// The best name for the landmark.
    DisplayName                 LocalizedText            `json:"display_name"`
	// One or more values indicating the type of the returned result. Please see <a
    // href="https://developers.google.com/maps/documentation/places/web-service/supported_types">Types
    // </a> for more detail.
    Types                       []string                 `json:"types"`
	// Defines the spatial relationship between the target location and the
    // landmark.
    SpatialRelationship         SpatialRelationship  `json:"spatial_relationship"`
	// The straight line distance between the target location and one of the
    // landmark's access points.
    StraightLineDistanceMeters  float32           `json:"straight_line_distance_meters"`
	// The travel distance along the road network between the target
    // location's closest point on a road, and the landmark's closest access
    // point on a road. This can be unpopulated if the landmark is disconnected
    // from the part of the road network the target is closest to OR if the
    // target location was not actually considered to be on the road network.
    TravelDistanceMeters        float32               `json:"travel_distance_meters"`
}

// Precise regions that are useful at describing a location.
type Area struct {
	// The Place ID of the underlying area feature. Can be used to
    // resolve more information about the area through Place Details or
    // Place Id Lookup.
    PlaceID     string              `json:"place_id"`
	// The best name for the area.
    DisplayName LocalizedText       `json:"display_name"`
	/**
	* An enum representing the relationship in space between the area and the target.
	*/
    Containment Containment     `json:"containment"`
}

/**
 * Represents a descriptor of an address.
 *
 * <p>Please see <a
 * href="https://mapsplatform.google.com/demos/address-descriptors/">Address 
 * Descriptors</a> for more detail.
 */
type AddressDescriptor struct {
	// A ranked list of nearby landmarks. The most useful (recognizable and
  	// nearby) landmarks are ranked first.
	Landmarks    []Landmark   `json:"landmarks"`
	// A ranked list of containing or adjacent areas. The most useful
  // (recognizable and precise) areas are ranked first.
	Areas        []Area       `json:"areas"`
}