// Copyright 2017 Google Inc. All Rights Reserved.
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

import "testing"

func TestParseLatLng(t *testing.T) {
	expected := &LatLng{Lat: 12.34, Lng: 56.78}
	actual, err := ParseLatLng("12.34,56.78")
	if err != nil {
		t.Errorf("Failed to parse simple LatLng %+v", err)
	}

	if !actual.AlmostEqual(expected, 0.0001) {
		t.Errorf("LatLng failed to parse expected value. Actual '%+v', expected '%+v'", actual, expected)
	}
}

func TestParseLatLngList(t *testing.T) {
	expected0 := &LatLng{Lat: 12.34, Lng: 56.78}
	expected1 := &LatLng{Lat: 14.89, Lng: 123.89}

	actual, err := ParseLatLngList("12.34,56.78|14.89,123.89")
	if err != nil {
		t.Errorf("Failed to parse LatLngList %+v", err)
	}

	if !actual[0].AlmostEqual(expected0, 0.0001) {
		t.Errorf("LatLng failed to parse expected value. Actual '%+v', expected '%+v'", actual[0], expected0)
	}

	if !actual[1].AlmostEqual(expected1, 0.0001) {
		t.Errorf("LatLng failed to parse expected value. Actual '%+v', expected '%+v'", actual[1], expected1)
	}
}
