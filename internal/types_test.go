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

package internal

import (
	"testing"
	"time"
)

func TestDateTime(t *testing.T) {
	if blank := NewDateTime(time.Time{}); blank != nil {
		t.Errorf("expected nil DateTime from zero time, was %v", blank)
	}

	loc, _ := time.LoadLocation("Australia/Sydney")
	orig := time.Date(2015, time.February, 25, 9, 9, 41, 0, loc)
	dt := NewDateTime(orig)
	if expected := "Australia/Sydney"; dt.TimeZone != expected {
		t.Errorf("expected timezone %v, was %v", expected, dt.TimeZone)
	}

	if actual := dt.Time(); !actual.Equal(orig) {
		t.Errorf("expected known time %v, was %v", orig, actual)
	}

	var blank *DateTime
	if expected := (time.Time{}); blank.Time() != expected {
		t.Errorf("expected nil DateTime to be zero time, was %v", blank.Time())
	}
}

func TestDuration(t *testing.T) {
	if empty := NewDuration(time.Duration(0)); empty.Text != "" || empty.Value != 0 {
		t.Errorf("expected empty duration, was %v", empty)
	}

	orig := time.Second * time.Duration(133)
	d := NewDuration(orig)
	if expected := "2m13s"; d.Text != expected {
		t.Errorf("expected text duration %v, was %v", expected, d.Text)
	}

	if actual := d.Duration(); actual != orig {
		t.Errorf("expected known duration %v, was %v", d, actual)
	}

	var blank *Duration
	if expected := time.Duration(0); blank.Duration() != expected {
		t.Errorf("expected nil Duration to be zero duration, was %v", blank.Duration())
	}
}
