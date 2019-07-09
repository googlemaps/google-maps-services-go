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

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	_ "image/jpeg" // Loaded for image decoder
	_ "image/png"  // Loaded for image decoder
)

var staticMapAPI = &apiConfig{
	host:             "https://maps.googleapis.com",
	path:             "/maps/api/staticmap",
	acceptsClientID:  true,
	acceptsSignature: true,
}

// MapType (optional) defines the type of map to construct. There are several possible
// maptype values, including roadmap, satellite, hybrid, and terrain
type MapType string

// Format defines the format of the resulting image
type Format string

// MarkerSize specifies the size of marker from the set {tiny, mid, small}
type MarkerSize string

// Anchor sets how the icon is placed in relation to the specified markers locations
type Anchor string

const (
	// RoadMap (default) specifies a standard roadmap image, as is normally shown on the
	// Google Maps website. If no maptype value is specified, the Google Static Maps API
	// serves roadmap tiles by default.
	RoadMap MapType = "roadmap"
	//Satellite specifies a satellite image.
	Satellite MapType = "satellite"
	//Terrain specifies a physical relief map image, showing terrain and vegetation.
	Terrain MapType = "terrain"
	// Hybrid specifies a hybrid of the satellite and roadmap image, showing a
	// transparent layer of major streets and place names on the satellite image.
	Hybrid MapType = "hybrid"
	//PNG8 or png (default) specifies the 8-bit PNG format.
	PNG8 Format = "png8"
	// PNG32 specifies the 32-bit PNG format.
	PNG32 Format = "png32"
	// GIF specifies the GIF format.
	GIF Format = "gif"
	// JPG specifies the JPEG compression format.
	JPG Format = "jpg"
	// JPGBaseline specifies a non-progressive JPEG compression format.
	JPGBaseline Format = "jpg-baseline"

	// Tiny Marker size
	Tiny MarkerSize = "tiny"
	// Mid Marker size
	Mid MarkerSize = "mid"
	// Small Marker size
	Small MarkerSize = "small"

	// Top Marker anchor position
	Top Anchor = "top"
	// Bottom Marker anchor position
	Bottom Anchor = "Bottom"
	// Left Marker anchor position
	Left Anchor = "left"
	// Right Marker anchor position
	Right Anchor = "right"
	// Center Marker anchor position
	Center Anchor = "center"
	// Topleft Marker anchor position
	Topleft Anchor = "topleft"
	// Topright Marker anchor position
	Topright Anchor = "topright"
	// Bottomleft Marker anchor position
	Bottomleft Anchor = "bottomleft"
	// Bottomright Marker anchor position
	Bottomright Anchor = "bottomright"
)

// CustomIcon replace the default Map Pin
type CustomIcon struct {
	// IconURL is th icon URL
	IconURL string
	// Anchor sets how the icon is placed in relation to the specified markers locations
	Anchor Anchor
	// Scale is the custom icon scale
	Scale int
}

func (c CustomIcon) String() string {
	var r []string

	if c.IconURL != "" {
		r = append(r, fmt.Sprintf("icon:%s", c.IconURL))
	}

	if c.Anchor != "" {
		r = append(r, fmt.Sprintf("anchor:%s", c.Anchor))
	}

	if c.Scale != 0 {
		r = append(r, fmt.Sprintf("scale:%d", c.Scale))
	}

	return strings.Join(r, "|")
}

// Marker is a Map pin
type Marker struct {
	// Color specifies a 24-bit color (example: color=0xFFFFCC) or a predefined color
	// from the set {black, brown, green, purple, yellow, blue, gray, orange, red,
	// white}.
	Color string
	// Label specifies a single uppercase alphanumeric character from the set {A-Z, 0-9}
	Label string
	// MarkerSize specifies the size of marker from the set {tiny, mid, small}
	Size string
	// CustomIcon replace the default Map Pin
	CustomIcon CustomIcon
	// Location is the Marker position
	Location []LatLng
	// LocationAddress is the Marker position as a postal address or other geocodable location.
	LocationAddress string
}

func (m Marker) String() string {
	var r []string

	if m.CustomIcon != (CustomIcon{}) {
		r = append(r, m.CustomIcon.String())
	} else {
		if m.Color != "" {
			r = append(r, fmt.Sprintf("color:%s", m.Color))
		}

		if m.Label != "" {
			r = append(r, fmt.Sprintf("label:%s", m.Label))
		}

		if m.Size != "" {
			r = append(r, fmt.Sprintf("size:%s", m.Size))
		}
	}

	for _, l := range m.Location {
		r = append(r, l.String())
	}
	if m.LocationAddress != "" {
		r = append(r, m.LocationAddress)
	}

	return strings.Join(r, "|")
}

// Path defines a single path of two or more connected points to overlay on the image
// at specified locations
type Path struct {
	// Weight (optional) specifies the thickness of the path in pixels.
	Weight int
	// Color (optional) specifies a color in HEX
	Color string
	// Fillcolor (optional) indicates both that the path marks off a polygonal area and
	// specifies the fill color to use as an overlay within that area.
	FillColor string
	// Geodesic (optional) indicates that the requested path should be interpreted as a
	// geodesic line that follows the curvature of the earth.
	Geodesic bool
	// Location two or more connected points to overlay on the image at specified
	// locations
	Location []LatLng
}

func (p Path) String() string {
	var r []string

	if p.Color != "" {
		r = append(r, fmt.Sprintf("color:%s", p.Color))
	}

	if p.FillColor != "" {
		r = append(r, fmt.Sprintf("fillcolor:%s", p.FillColor))
	}

	if p.Weight != 0 {
		r = append(r, fmt.Sprintf("weight:%d", p.Weight))
	}

	if p.Geodesic {
		r = append(r, "geodesic:true")
	}

	if len(p.Location) == 0 {
		return strings.Join(r, "|")
	}

	encodedLocationString := fmt.Sprintf("enc:%s", Encode(p.Location))
	latLngsToString := make([]string, len(p.Location))
	for i, l := range p.Location {
		latLngsToString[i] = l.String()
	}
	locationString := strings.Join(latLngsToString, "|")
	if len(locationString) > len(encodedLocationString) {
		r = append(r, encodedLocationString)
	} else {
		r = append(r, latLngsToString...)
	}
	return strings.Join(r, "|")
}

// StaticMapRequest is the functional options struct for staticMap.Get
type StaticMapRequest struct {
	// Center focus the map at the correct location
	Center string
	// Zoom (required if markers not present) defines the zoom level of the map
	Zoom int
	// Size (required) defines the rectangular dimensions of the map image. This
	// parameter takes a string of the form {horizontal_value}x{vertical_value}
	Size string
	// Scale (optional) affects the number of pixels that are returned. Accepted values
	// are 2 and 4
	Scale int
	// Format format (optional) defines the format of the resulting image. Default: PNG.
	// Accepted Values: There are several possible formats including GIF, JPEG and PNG
	// types.
	Format Format
	// Language (optional) defines the language to use for display of labels on map
	// tiles
	Language string
	// Region (optional) defines the appropriate borders to display, based on
	// geo-political sensitivities.
	Region string
	// MapType (optional) defines the type of map to construct.
	MapType MapType
	// Markers (optional) define one or more markers to attach to the image at specified
	// locations.
	Markers []Marker
	// Paths (optional) defines multiple paths of two or more connected points to
	// overlay on the image at specified locations
	Paths []Path
	// Visible specifies one or more locations that should remain visible on the map,
	// though no markers or other indicators will be displayed.
	Visible []LatLng
	// MapStyles (optional) contains map styles.
	MapStyles []string
}

func (r *StaticMapRequest) params() url.Values {
	q := make(url.Values)

	if r.Center != "" {
		q.Set("center", r.Center)
	}

	if r.Zoom > 0 {
		q.Set("zoom", strconv.Itoa(r.Zoom))
	}
	if r.Size != "" {
		q.Set("size", r.Size)
	}
	if r.Scale > 0 {
		q.Set("scale", strconv.Itoa(r.Scale))
	}
	if r.Format != "" {
		q.Set("format", string(r.Format))
	}

	if r.Language != "" {
		q.Set("language", r.Language)
	}

	if r.Region != "" {
		q.Set("region", r.Region)
	}
	if r.MapType != "" {
		q.Set("maptype", string(r.MapType))
	}

	for _, m := range r.Markers {
		q.Add("markers", m.String())
	}

	for _, ps := range r.Paths {
		q.Add("path", ps.String())
	}

	if len(r.Visible) > 0 {
		t := make([]string, len(r.Visible))
		for i, l := range r.Visible {
			t[i] = l.String()
		}
		q.Set("visible", strings.Join(t, "|"))
	}

	for _, style := range r.MapStyles {
		q.Add("style", style)
	}

	return q
}

// StaticMap makes a StaticMap API request.
func (c *Client) StaticMap(ctx context.Context, r *StaticMapRequest) (image.Image, error) {
	if len(r.Markers) == 0 && r.Center == "" && r.Zoom == 0 {
		return nil, errors.New("maps: Center & Zoom required if Markers empty")
	}
	if r.Size == "" {
		return nil, errors.New("maps: Size empty")
	}

	resp, err := c.getBinary(ctx, staticMapAPI, r)
	if err != nil {
		return nil, err
	}
	defer resp.data.Close()

	if resp.statusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.data)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Maps Static API: %d - %s", resp.statusCode, b)
	}

	img, _, err := image.Decode(resp.data)
	return img, err
}
