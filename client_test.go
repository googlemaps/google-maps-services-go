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
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClientChannelIsConfigured(t *testing.T) {
	_, err := NewClient(WithAPIKey("AIza-Maps-API-Key"), WithChannel("Test-Channel"))
	if err != nil {
		t.Errorf("Unable to create client with channel")
	}
}

func TestClientWithExperienceId(t *testing.T) {
	ids := []string{"foo", "bar"}
	c, err := NewClient(WithAPIKey("AIza-Maps-API-Key"), WithExperienceId(ids...))
	assert.Nil(t, err)
	assert.Equal(t, c.experienceId, ids)
}

func TestClientSetExperienceId(t *testing.T) {
	ids := []string{"foo", "bar"}
	c, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"))

	c.setExperienceId(ids...)
	assert.Equal(t, c.experienceId, ids)
}

func TestClientGetExperienceId(t *testing.T) {
	ids := []string{"foo", "bar"}
	c, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"))

	c.experienceId = ids
	assert.Equal(t, c.getExperienceId(), ids)
}

func TestClientClearExperienceId(t *testing.T) {
	ids := []string{"foo", "bar"}
	c, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"))

	c.experienceId = ids
	c.clearExperienceId()
	assert.Nil(t, c.experienceId)
}

func TestClientSetExperienceIdHeader(t *testing.T) {
	ids := []string{"foo", "bar"}
	c, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"))

	// slice has two elements
	c.experienceId = ids
	req, _ := http.NewRequest("GET", "/", nil)
	c.setExperienceIdHeader(context.Background(), req)
	assert.Equal(t, strings.Join(ids, ","), req.Header.Get(ExperienceIdHeaderName))

	// slice is nil
	c.experienceId = nil
	req, _ = http.NewRequest("GET", "/", nil)
	c.setExperienceIdHeader(context.Background(), req)
	assert.Equal(t, "", req.Header.Get(ExperienceIdHeaderName))

	// slice is empty
	c.experienceId = []string{}
	req, _ = http.NewRequest("GET", "/", nil)
	c.setExperienceIdHeader(context.Background(), req)
	assert.Equal(t, "", req.Header.Get(ExperienceIdHeaderName))

	var ctx context.Context
	// context has one element
	c.experienceId = []string{}
	ctx = context.Background()
	ctx = ExperienceIdContext(ctx, "foo")
	req, _ = http.NewRequest("GET", "/", nil)
	c.setExperienceIdHeader(ctx, req)
	assert.Equal(t, "foo", req.Header.Get(ExperienceIdHeaderName))

	// context has two elements
	c.experienceId = []string{}
	ctx = context.Background()
	ctx = ExperienceIdContext(ctx, ids...)
	req, _ = http.NewRequest("GET", "/", nil)
	c.setExperienceIdHeader(ctx, req)
	assert.Equal(t, strings.Join(ids, ","), req.Header.Get(ExperienceIdHeaderName))

	// context has two elements and client has two elements
	c.experienceId = ids
	ctx = context.Background()
	ctx = ExperienceIdContext(ctx, ids...)
	req, _ = http.NewRequest("GET", "/", nil)
	c.setExperienceIdHeader(ctx, req)
	assert.Equal(t, strings.Join(ids, ",") + "," + strings.Join(ids, ","), req.Header.Get(ExperienceIdHeaderName))
}

func TestClientExperienceIdSample(t *testing.T) {
	// [START maps_experience_id]
	experienceId := uuid.New().String()

	// instantiate client with experience id
	client, _ := NewClient(WithAPIKey("AIza-Maps-API-Key"), WithExperienceId("foo"))

	// clear the current experience id
	client.clearExperienceId()

	// set a new experience id
	otherExperienceId := uuid.New().String()
	client.setExperienceId(experienceId, otherExperienceId)

	// make API request, the client will set the header
	// X-GOOG-MAPS-EXPERIENCE-ID: experienceId,otherExperienceId

	// get current experience id
	var ids []string
	ids = client.getExperienceId()
	// [END maps_experience_id]

	assert.Equal(t, ids, []string{experienceId, otherExperienceId})
}
