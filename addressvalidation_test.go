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

// More information about Google Address Validation API is available on
// https://developers.google.com/maps/documentation/address-validation

package maps

import (
	"context"
	"testing"
)

func TestAddressValidation(t *testing.T) {
	response := `{
		"result": {
			"verdict": {
				"inputGranularity": "PREMISE",
				"validationGranularity": "PREMISE",
				"geocodeGranularity": "PREMISE",
				"addressComplete": true,
				"hasInferredComponents": true
			},
			"address": {
				"formattedAddress": "1600 Amphitheatre Parkway, Mountain View, CA 94043-1351, USA",
				"postalAddress": {
					"regionCode": "US",
					"languageCode": "en",
					"postalCode": "94043-1351",
					"administrativeArea": "CA",
					"locality": "Mountain View",
					"addressLines": [
						"1600 Amphitheatre Pkwy"
					]
				},
				"addressComponents": [
					{
						"componentName": {
							"text": "1600"
						},
						"componentType": "street_number",
						"confirmationLevel": "CONFIRMED"
					},
					{
						"componentName": {
							"text": "Amphitheatre Parkway",
							"languageCode": "en"
						},
						"componentType": "route",
						"confirmationLevel": "CONFIRMED"
					},
					{
						"componentName": {
							"text": "USA",
							"languageCode": "en"
						},
						"componentType": "country",
						"confirmationLevel": "CONFIRMED"
					},
					{
						"componentName": {
							"text": "Mountain View",
							"languageCode": "en"
						},
						"componentType": "locality",
						"confirmationLevel": "CONFIRMED",
						"inferred": true
					},
					{
						"componentName": {
							"text": "94043"
						},
						"componentType": "postal_code",
						"confirmationLevel": "CONFIRMED",
						"inferred": true
					},
					{
						"componentName": {
							"text": "CA",
							"languageCode": "en"
						},
						"componentType": "administrative_area_level_1",
						"confirmationLevel": "CONFIRMED",
						"inferred": true
					},
					{
						"componentName": {
							"text": "1351"
						},
						"componentType": "postal_code_suffix",
						"confirmationLevel": "CONFIRMED",
						"inferred": true
					}
				]
			},
			"geocode": {
				"location": {
					"latitude": 37.4223878,
					"longitude": -122.0841877
				},
				"plusCode": {
					"globalCode": "849VCWC8+X8"
				},
				"bounds": {
					"low": {
						"latitude": 37.4220699,
						"longitude": -122.084958
					},
					"high": {
						"latitude": 37.4226618,
						"longitude": -122.0829302
					}
				},
				"featureSizeMeters": 116.538734,
				"placeId": "ChIJj38IfwK6j4ARNcyPDnEGa9g",
				"placeTypes": [
					"premise"
				]
			},
			"metadata": {
				"business": true,
				"poBox": false
			},
			"uspsData": {
				"standardizedAddress": {
					"firstAddressLine": "1600 AMPHITHEATRE PKWY",
					"cityStateZipAddressLine": "MOUNTAIN VIEW CA 94043-1351",
					"city": "MOUNTAIN VIEW",
					"state": "CA",
					"zipCode": "94043",
					"zipCodeExtension": "1351"
				},
				"deliveryPointCode": "00",
				"deliveryPointCheckDigit": "0",
				"dpvConfirmation": "Y",
				"dpvFootnote": "AABB",
				"dpvCmra": "N",
				"dpvVacant": "N",
				"dpvNoStat": "Y",
				"carrierRoute": "C909",
				"carrierRouteIndicator": "D",
				"postOfficeCity": "MOUNTAIN VIEW",
				"postOfficeState": "CA",
				"fipsCountyCode": "085",
				"county": "SANTA CLARA",
				"elotNumber": "0103",
				"elotFlag": "A",
				"addressRecordType": "S"
			}
		},
		"responseId": "555af283-4afe-4a26-a192-b837ceac518a"
	}`
	server := mockServer(200, response)
	defer server.Close()
	c, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))
	r := &AddressValidationRequest{
		Address: &PostalAddress{
			RegionCode: "US",
			AddressLines: []string{
				"1600 Amphitheatre Parkway",
			},
		},
	}
	resp, err := c.ValidateAddress(context.Background(), r)
	if err != nil {
		t.Fatalf("ValidateAddress returned error: %v", err)
	}
	if resp.Result.Address.FormattedAddress != "1600 Amphitheatre Parkway, Mountain View, CA 94043-1351, USA" {
		t.Errorf("ValidateAddress returned unexpected result: %+v", resp)
	}
}
