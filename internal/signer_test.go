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
	"encoding/base64"
	"encoding/hex"
	"testing"
)

// From http://en.wikipedia.org/wiki/Hash-based_message_authentication_code
// HMAC_SHA1("key", "The quick brown fox jumps over the lazy dog")
//     = 0xde7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9
var message = "The quick brown fox jumps over the lazy dog"
var signingKey = []byte("key")
var signature = "de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9"

func TestSigner(t *testing.T) {
	s, err := hex.DecodeString(signature)
	if err != nil {
		t.Errorf("Couldn't decode expected signature: %+v", err)
	}
	expected := base64.URLEncoding.EncodeToString(s)
	generated, err := generateSignature(signingKey, message)
	if err != nil {
		t.Errorf("Couldn't generate actual signature: %+v", err)
	}
	if expected != generated {
		t.Errorf("expected equal signature, was %s, expected %s", generated, expected)
	}
}
