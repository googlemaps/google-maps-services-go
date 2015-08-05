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
	"net/http"
	"testing"
)

func TestClientTransportMutate(t *testing.T) {
	c, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"), WithHTTPClient(&http.Client{}))
	tr, ok := c.httpClient.Transport.(*transport)
	if !ok {
		t.Errorf("Transport is expected to be a maps.transport, found to be a %T", c.httpClient.Transport)
	}
	if _, ok := tr.Base.(*transport); ok {
		t.Errorf("Transport's Base shouldn't have been a maps.transport, found to be a %T", tr.Base)
	}
}

func TestDefaultClientTransportMutate(t *testing.T) {
	c, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"))

	tr, ok := c.httpClient.Transport.(*transport)
	if !ok {
		t.Errorf("Transport is expected to be a maps.transport, found to be a %T", c.httpClient.Transport)
	}
	if _, ok := tr.Base.(*transport); ok {
		t.Errorf("Transport's Base shouldn't have been a maps.transport, found to be a %T", tr.Base)
	}
}
