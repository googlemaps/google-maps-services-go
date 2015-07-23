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
	"bytes"
	"io"
	"log"
)

// Polyline represents a list of lat,lng points encoded as a byte array.
// See: https://developers.google.com/maps/documentation/utilities/polylinealgorithm
type Polyline struct {
	Points string `json:"points"`
}

// DecodePolyline converts a polyline encoded string to an array of LatLng objects.
func DecodePolyline(poly string) []LatLng {
	p := &Polyline{
		Points: poly,
	}
	return p.Decode()
}

// Decode converts this encoded Polyline to an array of LatLng objects.
func (p *Polyline) Decode() []LatLng {
	input := bytes.NewBufferString(p.Points)

	var lat, lng int64
	path := make([]LatLng, 0, len(p.Points)/2)
	for {
		dlat, _ := decodeInt(input)
		dlng, err := decodeInt(input)
		if err == io.EOF {
			return path
		}
		if err != nil {
			log.Fatal("unexpected err decoding polyline", err)
		}

		lat, lng = lat+dlat, lng+dlng
		path = append(path, LatLng{
			Lat: float64(lat) * 1e-5,
			Lng: float64(lng) * 1e-5,
		})
	}
}

// Encode returns a new encoded Polyline from the given path.
func Encode(path []LatLng) string {
	var prevLat, prevLng int64

	out := new(bytes.Buffer)
	out.Grow(len(path) * 4)

	for _, point := range path {
		lat := int64(point.Lat * 1e5)
		lng := int64(point.Lng * 1e5)

		encodeInt(lat-prevLat, out)
		encodeInt(lng-prevLng, out)

		prevLat, prevLng = lat, lng
	}

	return out.String()
}

// decodeInt reads an encoded int64 from the passed io.ByteReader.
func decodeInt(r io.ByteReader) (int64, error) {
	result := int64(1)
	var shift uint8

	for {
		raw, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		b := raw - 63 - 1
		result += int64(b) << shift
		shift += 5

		if b < 0x1f {
			bit := result & 1
			result >>= 1
			if bit != 0 {
				result = ^result
			}
			return result, nil
		}
	}
}

// encodeInt writes an encoded int64 to the passed io.ByteWriter.
func encodeInt(v int64, w io.ByteWriter) {
	if v < 0 {
		v = ^(v << 1)
	} else {
		v <<= 1
	}
	for v >= 0x20 {
		w.WriteByte((0x20 | (byte(v) & 0x1f)) + 63)
		v >>= 5
	}
	w.WriteByte(byte(v) + 63)
}
