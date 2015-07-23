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
	"encoding/json"
	"reflect"
	"testing"
)

const (
	jsonSnappedPoint = `{"originalIndex":null,"placeId":"helloPlace","location":{"latitude":-33.870315,"longitude":151.196532}}`
)

func TestSnappedPoint(t *testing.T) {
	sp := SnappedPoint{
		Location: LatLng{
			Lat: -33.870315,
			Lng: 151.196532,
		},
		PlaceID: "helloPlace",
	}

	bytes, err := json.Marshal(&sp)
	if err != nil {
		t.Errorf("expected ok encode of SnappedPoint, got: %v", err)
	}
	if string(bytes) != jsonSnappedPoint {
		t.Errorf("expected encoded snappedPoint, was: %v", string(bytes))
	}

	var out SnappedPoint
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		t.Errorf("expected ok decode of SnappedPoint, got: %v", err)
	}
	if !reflect.DeepEqual(out, sp) {
		t.Errorf("expected equal snappedPoint, was %+v expected %+v", out, sp)
	}
}
