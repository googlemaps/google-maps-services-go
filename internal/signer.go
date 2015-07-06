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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// GenerateSignature builds the digital signature for a key and a message.
// The key is assumed to be in URL safe base64 encoding.
// See: https://developers.google.com/maps/documentation/business/webservices/auth#digital_signatures
func GenerateSignature(key, message string) (string, error) {
	k, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha1.New, k)
	mac.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil)), nil
}
