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

package directions

import (
	"encoding/json"
	"google.golang.org/maps/internal"
)

// safeLeg is a raw version of Leg that does not have custom encoding or
// decoding methods applied.
type safeLeg Leg

// encodedLeg is the actual encoded version of Leg as per the Maps APIs.
type encodedLeg struct {
	safeLeg
	EncDuration      *internal.Duration `json:"duration"`
	EncArrivalTime   *internal.DateTime `json:"arrival_time"`
	EncDepartureTime *internal.DateTime `json:"departure_time"`
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
