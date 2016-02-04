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

import "time"

// DateTime is the public API representation of a point in time.
type DateTime struct {
	// Text is the time specified as a string. The time is displayed in
	// the corresponding TimeZone.
	Text string `json:"text"`

	// TimeZone is the name of the time zone in the IANA Time Zone Database. For
	// example, "America/New_York" or "Australia/Sydney".
	TimeZone string `json:"time_zone"`

	// Value is the number of seconds since midnight 01 January, 1970 UTC.
	Value int64 `json:"value"`
}

// Time returns the time.Time for this DateTime.
func (dt *DateTime) Time() time.Time {
	if dt == nil {
		return time.Time{}
	}

	loc, err := time.LoadLocation(dt.TimeZone)
	t := time.Unix(dt.Value, 0)
	if err == nil && loc != nil {
		t = t.In(loc)
	}
	return t
}

// NewDateTime builds a DateTime from the given time.Time. This will be nil
// if time.Time is the zero time.
func NewDateTime(t time.Time) *DateTime {
	if t.IsZero() {
		return nil
	}

	loc := t.Location()
	return &DateTime{
		Text:     t.Format(time.RFC1123), // TODO(samthor): better format
		TimeZone: loc.String(),
		Value:    t.UnixNano() / int64(time.Second),
	}
}

// Duration is the public API representation of a duration.
type Duration struct {
	// Value indicates the duration, in seconds.
	Value int64 `json:"value"`

	// Text contains a human-readable representation of the duration.
	Text string `json:"text"`
}

// Duration returns the time.Duration for this internal Duration.
func (d *Duration) Duration() time.Duration {
	if d == nil {
		return time.Duration(0)
	}
	return time.Duration(d.Value) * time.Second
}

// NewDuration builds an internal Duration from the passed time.Duration.
func NewDuration(d time.Duration) *Duration {
	if d == 0 {
		return &Duration{}
	}
	return &Duration{
		Value: int64(d / time.Second),
		Text:  d.String(),
	}
}

// Location is the Roads API+ representation of a location on the Earth. It
// differs only in the encoded names, which are the longer forms of 'latitude'
// and 'longitude'.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
