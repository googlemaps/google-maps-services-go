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
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

// mockServerForQueryWithImage returns a mock server that only responds to a particular query string, and responds with an encoded Image.
func mockServerForQueryWithImage(query string, code int, img image.Image) *countingServer {
	server := &countingServer{}

	server.s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if query != "" && r.URL.RawQuery != query {
			server.failed = append(server.failed, r.URL.RawQuery)
			http.Error(w, fmt.Sprintf("Expected '%s', got '%s'", query, r.URL.RawQuery), 999)
			return
		}
		server.successful++

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
	}))

	return server
}

func TestStaticMode(t *testing.T) {

	response := image.NewRGBA(image.Rect(0, 0, 640, 400))

	server := mockServerForQueryWithImage("center=Brooklyn+Bridge%2CNew+York%2CNY&format=PNG&key=AIzaNotReallyAnAPIKey&language=EN-us&maptype=roadmap&region=US&scale=2&size=600x300&zoom=13", 200, response)
	defer server.s.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.s.URL))
	r := &StaticMapRequest{
		Center:   "Brooklyn Bridge,New York,NY",
		Size:     "600x300",
		Zoom:     13,
		Scale:    2,
		Language: "EN-us",
		Format:   "PNG",
		Region:   "US",
		MapType:  MapType("roadmap"),
	}

	resp, err := c.StaticMap(context.Background(), r)
	if err != nil {
		t.Errorf("r.StaticMap returned non nil error: %+v", err)
	}

	if resp.Bounds().Min.X != 0 || resp.Bounds().Min.Y != 0 || resp.Bounds().Max.X != 640 || resp.Bounds().Max.Y != 400 {
		t.Errorf("Response image not of the correct dimensions")
	}
}

func TestMapStylesMultiStyles(t *testing.T) {
	r := StaticMapRequest{
		MapStyles: MapStyle{
			"landscape.natural": Elements{
				"geometry": StyleRules{
					"color": "0x0000FF",
					"width": "50",
				},
			},
		},
	}

	values := r.params()
	t.Log(values)
	if c := strings.Count(values.Encode(), "style"); c != 1 {
		t.Errorf("Generate query string does not contain sufficient Style parameters (found %d)", c)
	}
}

func TestMapStylesFeaturesAndElements(t *testing.T) {
	r := StaticMapRequest{
		MapStyles: MapStyle{
			"water": Elements{
				"geometry.fill": StyleRules{
					"color": "0xFF0000",
				},
			},
			"landscape.natural": Elements{
				"geometry": StyleRules{
					"color": "0x0000FF",
					"width": "50",
				},
			},
		},
	}

	values := r.params()
	t.Log(values)
	if c := strings.Count(values.Encode(), "style"); c != 2 {
		t.Errorf("Generate query string does not contain sufficient Style parameters (found %d)", c)
	}
}

func TestMapStylesMultiElementsPerFeature(t *testing.T) {
	r := StaticMapRequest{
		MapStyles: MapStyle{
			"all": Elements{
				"geometry": StyleRules{
					"color":      "0xeeaabb",
					"visibility": "simplified",
				},
				"labels.icon": StyleRules{
					"saturation": "-50",
					"weight":     "2",
				},
			},
		},
	}

	values := r.params()
	t.Log(values)
	if c := strings.Count(values.Encode(), "style"); c != 2 {
		t.Errorf("Generate query string does not contain sufficient Style parameters (found %d)", c)
	}
}
