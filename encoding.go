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
	"encoding/json"
	"net/url"

	"googlemaps.github.io/maps/internal"
)

// safeLeg is a raw version of Leg that does not have custom encoding or
// decoding methods applied.
type safeLeg Leg

// encodedLeg is the actual encoded version of Leg as per the Maps APIs.
type encodedLeg struct {
	safeLeg
	EncDuration          *internal.Duration `json:"duration"`
	EncDurationInTraffic *internal.Duration `json:"duration_in_traffic"`
	EncArrivalTime       *internal.DateTime `json:"arrival_time"`
	EncDepartureTime     *internal.DateTime `json:"departure_time"`
}

// UnmarshalJSON implements json.Unmarshaler for Leg. This decodes the API
// representation into types useful for Go developers.
func (leg *Leg) UnmarshalJSON(data []byte) error {
	x := encodedLeg{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*leg = Leg(x.safeLeg)

	leg.Duration = x.EncDuration.Duration()
	leg.DurationInTraffic = x.EncDurationInTraffic.Duration()
	leg.ArrivalTime = x.EncArrivalTime.Time()
	leg.DepartureTime = x.EncDepartureTime.Time()

	return nil
}

// MarshalJSON implements json.Marshaler for Leg. This encodes Go types back to
// the API representation.
func (leg *Leg) MarshalJSON() ([]byte, error) {
	x := encodedLeg{}
	x.safeLeg = safeLeg(*leg)

	x.EncDuration = internal.NewDuration(leg.Duration)
	x.EncDurationInTraffic = internal.NewDuration(leg.DurationInTraffic)
	x.EncArrivalTime = internal.NewDateTime(leg.ArrivalTime)
	x.EncDepartureTime = internal.NewDateTime(leg.DepartureTime)

	return json.Marshal(x)
}

// safeStep is a raw version of Step that does not have custom encoding or
// decoding methods applied.
type safeStep Step

// encodedStep is the actual encoded version of Step as per the Maps APIs.
type encodedStep struct {
	safeStep
	EncDuration *internal.Duration `json:"duration"`
}

// UnmarshalJSON implements json.Unmarshaler for Step. This decodes the API
// representation into types useful for Go developers.
func (step *Step) UnmarshalJSON(data []byte) error {
	x := encodedStep{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*step = Step(x.safeStep)

	step.Duration = x.EncDuration.Duration()

	return nil
}

// MarshalJSON implements json.Marshaler for Step. This encodes Go types back to
// the API representation.
func (step *Step) MarshalJSON() ([]byte, error) {
	x := encodedStep{}
	x.safeStep = safeStep(*step)

	x.EncDuration = internal.NewDuration(step.Duration)

	return json.Marshal(x)
}

// safeTransitDetails is a raw version of TransitDetails that does not have
// custom encoding or decoding methods applied.
type safeTransitDetails TransitDetails

// encodedTransitDetails is the actual encoded version of TransitDetails as per
// the Maps APIs
type encodedTransitDetails struct {
	safeTransitDetails
	EncArrivalTime   *internal.DateTime `json:"arrival_time"`
	EncDepartureTime *internal.DateTime `json:"departure_time"`
}

// UnmarshalJSON implements json.Unmarshaler for TransitDetails. This decodes
// the API representation into types useful for Go developers.
func (transitDetails *TransitDetails) UnmarshalJSON(data []byte) error {
	x := encodedTransitDetails{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*transitDetails = TransitDetails(x.safeTransitDetails)

	transitDetails.ArrivalTime = x.EncArrivalTime.Time()
	transitDetails.DepartureTime = x.EncDepartureTime.Time()

	return nil
}

// MarshalJSON implements json.Marshaler for TransitDetails. This encodes Go
// types back to the API representation.
func (transitDetails *TransitDetails) MarshalJSON() ([]byte, error) {
	x := encodedTransitDetails{}
	x.safeTransitDetails = safeTransitDetails(*transitDetails)

	x.EncArrivalTime = internal.NewDateTime(transitDetails.ArrivalTime)
	x.EncDepartureTime = internal.NewDateTime(transitDetails.DepartureTime)

	return json.Marshal(x)
}

// safeTransitLine is the raw version of TransitLine that does not have custom
// encoding or decoding methods applied.
type safeTransitLine TransitLine

// encodedTransitLine is the actual encoded version of TransitLine as per the
// Maps APIs
type encodedTransitLine struct {
	safeTransitLine
	EncURL  string `json:"url"`
	EncIcon string `json:"icon"`
}

// UnmarshalJSON imlpements json.Unmarshaler for TransitLine. This decodes the
// API representation into types useful for Go developers.
func (transitLine *TransitLine) UnmarshalJSON(data []byte) error {
	x := encodedTransitLine{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*transitLine = TransitLine(x.safeTransitLine)

	transitLine.URL, err = url.Parse(x.EncURL)
	if err != nil {
		return err
	}
	transitLine.Icon, err = url.Parse(x.EncIcon)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements json.Marshaler for TransitLine. This encodes Go
// types back to the API representation.
func (transitLine *TransitLine) MarshalJSON() ([]byte, error) {
	x := encodedTransitLine{}
	x.safeTransitLine = safeTransitLine(*transitLine)

	x.EncURL = transitLine.URL.String()
	x.EncIcon = transitLine.Icon.String()

	return json.Marshal(x)
}

// safeTransitAgency is the raw version of TransitAgency that does not have
// custom encoding or decoding methods applied.
type safeTransitAgency TransitAgency

// encodedTransitAgency is the actual encoded version of TransitAgency as per the
// Maps APIs
type encodedTransitAgency struct {
	safeTransitAgency
	EncURL string `json:"url"`
}

// UnmarshalJSON imlpements json.Unmarshaler for TransitAgency. This decodes the
// API representation into types useful for Go developers.
func (transitAgency *TransitAgency) UnmarshalJSON(data []byte) error {
	x := encodedTransitAgency{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*transitAgency = TransitAgency(x.safeTransitAgency)

	transitAgency.URL, err = url.Parse(x.EncURL)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements json.Marshaler for TransitAgency. This encodes Go
// types back to the API representation.
func (transitAgency *TransitAgency) MarshalJSON() ([]byte, error) {
	x := encodedTransitAgency{}
	x.safeTransitAgency = safeTransitAgency(*transitAgency)

	x.EncURL = transitAgency.URL.String()

	return json.Marshal(x)
}

// safeTransitLineVehicle is the raw version of TransitLineVehicle that does not
// have custom encoding or decoding methods applied.
type safeTransitLineVehicle TransitLineVehicle

// encodedTransitLineVehicle is the actual encoded version of TransitLineVehicle
// as per the Maps APIs
type encodedTransitLineVehicle struct {
	safeTransitLineVehicle
	EncIcon string `json:"icon"`
}

// UnmarshalJSON imlpements json.Unmarshaler for TransitLineVehicle. This
// decodes the API representation into types useful for Go developers.
func (transitLineVehicle *TransitLineVehicle) UnmarshalJSON(data []byte) error {
	x := encodedTransitLineVehicle{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*transitLineVehicle = TransitLineVehicle(x.safeTransitLineVehicle)

	transitLineVehicle.Icon, err = url.Parse(x.EncIcon)
	if err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements json.Marshaler for TransitLineVehicle. This encodes
// Go types back to the API representation.
func (transitLineVehicle *TransitLineVehicle) MarshalJSON() ([]byte, error) {
	x := encodedTransitLineVehicle{}
	x.safeTransitLineVehicle = safeTransitLineVehicle(*transitLineVehicle)

	x.EncIcon = transitLineVehicle.Icon.String()

	return json.Marshal(x)
}

// safeDistanceMatrixElement is a raw version of DistanceMatrixElement that
// does not have custom encoding or decoding methods applied.
type safeDistanceMatrixElement DistanceMatrixElement

// encodedDistanceMatrixElement is the actual encoded version of
// DistanceMatrixElement as per the Maps APIs.
type encodedDistanceMatrixElement struct {
	safeDistanceMatrixElement
	EncDuration          *internal.Duration `json:"duration"`
	EncDurationInTraffic *internal.Duration `json:"duration_in_traffic"`
}

// UnmarshalJSON implements json.Unmarshaler for DistanceMatrixElement. This
// decodes the API representation into types useful for Go developers.
func (dme *DistanceMatrixElement) UnmarshalJSON(data []byte) error {
	x := encodedDistanceMatrixElement{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*dme = DistanceMatrixElement(x.safeDistanceMatrixElement)

	dme.Duration = x.EncDuration.Duration()
	dme.DurationInTraffic = x.EncDurationInTraffic.Duration()

	return nil
}

// MarshalJSON implements json.Marshaler for DistanceMatrixElement. This encodes
// Go types back to the API representation.
func (dme *DistanceMatrixElement) MarshalJSON() ([]byte, error) {
	x := encodedDistanceMatrixElement{}
	x.safeDistanceMatrixElement = safeDistanceMatrixElement(*dme)

	x.EncDuration = internal.NewDuration(dme.Duration)
	x.EncDurationInTraffic = internal.NewDuration(dme.DurationInTraffic)

	return json.Marshal(x)
}

// safeSnappedPoint is a raw version of SnappedPoint that does not have custom
// encoding or decoding methods applied.
type safeSnappedPoint SnappedPoint

// encodedSnappedPoint is the actual encoded version of SnappedPoint as per the
// Roads API.
type encodedSnappedPoint struct {
	safeSnappedPoint
	EncLocation internal.Location `json:"location"`
}

// UnmarshalJSON implements json.Unmarshaler for SnappedPoint. This decode the
// API representation into types useful for Go developers.
func (sp *SnappedPoint) UnmarshalJSON(data []byte) error {
	x := encodedSnappedPoint{}
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	*sp = SnappedPoint(x.safeSnappedPoint)

	sp.Location.Lat = x.EncLocation.Latitude
	sp.Location.Lng = x.EncLocation.Longitude

	return nil
}

// MarshalJSON implements json.Marshaler for SnappedPoint. This encodes Go
// types back to the API representation.
func (sp *SnappedPoint) MarshalJSON() ([]byte, error) {
	x := encodedSnappedPoint{}
	x.safeSnappedPoint = safeSnappedPoint(*sp)

	x.EncLocation.Latitude = sp.Location.Lat
	x.EncLocation.Longitude = sp.Location.Lng

	return json.Marshal(x)
}
