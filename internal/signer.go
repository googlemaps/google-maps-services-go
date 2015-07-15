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
	"fmt"
	"net/url"
)

// generateSignature generates the base64 URL-encoded HMAC-SHA1 signature for the key and plaintext message.
func generateSignature(key []byte, message string) (string, error) {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil)), nil
}

// SignURL signs a url with a clientID and signature.
// The signature is assumed to be in URL safe base64 encoding.
// The returned signature string is URLEncoded.
// See: https://developers.google.com/maps/documentation/business/webservices/auth#digital_signatures
func SignURL(path, clientID string, signature []byte, q url.Values) (string, error) {
	q.Set("client", clientID)
	encodedQuery := q.Encode()
	message := fmt.Sprintf("%s?%s", path, encodedQuery)
	s, err := generateSignature(signature, message)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s&signature=%s", encodedQuery, s), nil
}
