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
	"reflect"
	"testing"
)

const (
	epsilon     = 0.0001
	routeSydMel = "rvumEis{y[`NsfA~tAbF`bEj^h{@{KlfA~eA~`AbmEghAt~D|e@jlRpO~yH_\\v}LjbBh~FdvCxu@`nCplDbcBf_B|w" +
		"BhIfhCnqEb~D~jCn_EngApdEtoBbfClf@t_CzcCpoEr_Gz_DxmAphDjjBxqCviEf}B|pEvsEzbE~qGfpExjBlqCx}" +
		"BvmLb`FbrQdpEvkAbjDllD|uDldDj`Ef|AzcEx_Gtm@vuI~xArwD`dArlFnhEzmHjtC~eDluAfkC|eAdhGpJh}N_m" +
		"ArrDlr@h|HzjDbsAvy@~~EdTxpJje@jlEltBboDjJdvKyZpzExrAxpHfg@pmJg[tgJuqBnlIarAh}DbN`hCeOf_Ib" +
		"xA~uFt|A|xEt_ArmBcN|sB|h@b_DjOzbJ{RlxCcfAp~AahAbqG~Gr}AerA`dCwlCbaFo]twKt{@bsG|}A~fDlvBvz" +
		"@tw@rpD_r@rqB{PvbHek@vsHlh@ptNtm@fkD[~xFeEbyKnjDdyDbbBtuA|~Br|Gx_AfxCt}CjnHv`Ew\\lnBdrBfq" +
		"BraD|{BldBxpG|]jqC`mArcBv]rdAxgBzdEb{InaBzyC}AzaEaIvrCzcAzsCtfD~qGoPfeEh]h`BxiB`e@`kBxfAv" +
		"^pyA`}BhkCdoCtrC~bCxhCbgEplKrk@tiAteBwAxbCwuAnnCc]b{FjrDdjGhhGzfCrlDruBzSrnGhvDhcFzw@n{@z" +
		"xAf}Fd{IzaDnbDjoAjqJjfDlbIlzAraBxrB}K~`GpuD~`BjmDhkBp{@r_AxCrnAjrCx`AzrBj{B|r@~qBbdAjtDnv" +
		"CtNzpHxeApyC|GlfM`fHtMvqLjuEtlDvoFbnCt|@xmAvqBkGreFm~@hlHw|AltC}NtkGvhBfaJ|~@riAxuC~gErwC" +
		"ttCzjAdmGuF`iFv`AxsJftD|nDr_QtbMz_DheAf~Buy@rlC`i@d_CljC`gBr|H|nAf_Fh{G|mE~kAhgKviEpaQnu@" +
		"zwAlrA`G~gFnvItz@j{Cng@j{D{]`tEftCdcIsPz{DddE~}PlnE|dJnzG`eG`mF|aJdqDvoAwWjzHv`H`wOtjGzeX" +
		"hhBlxErfCf{BtsCjpEjtD|}Aja@xnAbdDt|ErMrdFh{CzgAnlCnr@`wEM~mE`bA`uD|MlwKxmBvuFlhB|sN`_@fvB" +
		"p`CxhCt_@loDsS|eDlmChgFlqCbjCxk@vbGxmCjbMba@rpBaoClcCk_DhgEzYdzBl\\vsA_JfGztAbShkGtEhlDzh" +
		"C~w@hnB{e@yF}`D`_Ayx@~vGqn@l}CafC"
	routeWith0b = "ynkrFq|zfE?sCnBpA"
)

var expectedWith0b = []LatLng{{Lat: 39.87709, Lng: 32.74713}, {Lat: 39.87709, Lng: 32.74787}, {Lat: 39.87653, Lng: 32.74746}}

func TestPolylineDecode(t *testing.T) {
	decoded := DecodePolyline(routeSydMel)
	l := len(decoded)

	if expected, actual := (&LatLng{-33.86746, 151.207090}), decoded[0]; !expected.AlmostEqual(&actual, epsilon) {
		t.Errorf("first point was %v, expected %v", decoded[0], expected)
	}
	if expected, actual := (&LatLng{-37.814130, 144.963180}), decoded[l-1]; !expected.AlmostEqual(&actual, epsilon) {
		t.Errorf("last point was %v, expected %v", decoded[l-1], expected)
	}

	decoded = DecodePolyline(routeWith0b)
	if len(decoded) != len(expectedWith0b) {
		t.Errorf("expected equal encoding for 0 change in one direction length mismatch, was len %v, expected len %v", len(decoded), len(expectedWith0b))
	}
	for i := range decoded {
		if expected, actual := (expectedWith0b[i]), (decoded[i]); !expected.AlmostEqual(&actual, epsilon) {
			t.Errorf("expected equal encoding for 0 change in one direction mismatch coordinate, was %+v, expected %+v", decoded, expectedWith0b)
			break
		}
	}
}

func TestPolylineEncode(t *testing.T) {
	decoded := DecodePolyline(routeSydMel)

	encoded := Encode(decoded)
	if !reflect.DeepEqual(encoded, routeSydMel) {
		t.Errorf("expected equal encoding, was len %v, expected len %v", len(routeSydMel), len(encoded))
	}

	decoded2 := DecodePolyline(routeWith0b)

	encoded2 := Encode(decoded2)
	if !reflect.DeepEqual(encoded2, routeWith0b) {
		t.Errorf("expected 2 equal encoding, was len %v, expected len %v", len(routeWith0b), len(encoded2))
	}
}
