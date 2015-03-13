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
package maps // import "google.golang.org/maps"

import (
	"encoding/json"

	"google.golang.org/maps/internal"
)

// safeDistanceMatrixElement is a raw version of DistanceMatrixElement that
// does not have custom encoding or decoding methods applied.
type safeDistanceMatrixElement DistanceMatrixElement

// encodedDistanceMatrixElement is the actual encoded version of
// DistanceMatrixElement as per the Maps APIs.
type encodedDistanceMatrixElement struct {
	safeDistanceMatrixElement
	EncDuration *internal.Duration `json:"duration"`
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

	return nil
}

// MarshalJSON implements json.Marshaler for DistanceMatrixElement. This encodes
// Go types back to the API representation.
func (dme *DistanceMatrixElement) MarshalJSON() ([]byte, error) {
	x := encodedDistanceMatrixElement{}
	x.safeDistanceMatrixElement = safeDistanceMatrixElement(*dme)

	x.EncDuration = internal.NewDuration(dme.Duration)

	return json.Marshal(x)
}
