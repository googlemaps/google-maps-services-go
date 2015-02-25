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

package maps // import "google.golang.org/maps"

import (
  "bytes"
  "io"
)

// Polyline represents a list of lat,lng points, encoded as a string.
// See: https://developers.google.com/maps/documentation/utilities/polylinealgorithm
type Polyline struct {
	Points string `json:"points"`
}

func decodeString(s string, ch chan int64) {
  result := int64(1)
  var shift uint8

  for _, c := range s {
    b := c - 63 - 1
    result += int64(b) << shift
    shift += 5
    if b >= 0x1f {
      continue
    }

    bit := result & 1
    result >>= 1
    if bit != 0 {
      result = ^result
    }
    ch <- result
    result = 1
    shift = 0
  }
  close(ch)
}

// Decode converts this encoded Polyline to an array of LatLng objects.
func (p *Polyline) Decode() []LatLng {
  ch := make(chan int64)
  go decodeString(p.Points, ch)

  var lat, lng int64
  path := make([]LatLng, 0, len(p.Points) / 2)
  for {
    dlat, _ := <-ch
    dlng, ok := <-ch
    if !ok {
      return path
    }
    lat, lng = lat + dlat, lng + dlng
    path = append(path, LatLng{
      Lat: float64(lat) * 1e-5,
      Lng: float64(lng) * 1e-5,
    })
  }
  panic("should not get here")
}

func encode(v int64, w io.Writer) {
  if v < 0 {
    v = ^(v << 1)
  } else {
    v <<= 1
  }
  for v >= 0x20 {
    w.Write([]byte{(0x20 | (byte(v) & 0x1f)) + 63})
    v >>= 5
  }
  w.Write([]byte{byte(v) + 63})
}

func Encode(path []LatLng) *Polyline {
  var plat, plng int64

  out := new(bytes.Buffer)
  for _, point := range path {
    lat, lng := int64(point.Lat * 1e5), int64(point.Lng * 1e5)
    dlat, dlng := lat - plat, lng - plng

    encode(dlat, out)
    encode(dlng, out)

    plat, plng = lat, lng
  }

  return &Polyline{Points: out.String()}
}
