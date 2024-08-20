//nolint:lll
package credential_test

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/axone-sdk/credential"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Secp256k1PubKeyFetcher(t *testing.T) {
	tests := []struct {
		name        string
		issuerID    string
		keyID       string
		wantErr     error
		wantSuccess bool
	}{
		{
			name:        "valid key",
			issuerID:    "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			keyID:       "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:     nil,
			wantSuccess: true,
		},
		{
			name:        "invalid key fingerprint",
			issuerID:    "did:key:invalid",
			keyID:       "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:     fmt.Errorf("failed to get key fingerprint: unknown key encoding"),
			wantSuccess: false,
		},
		{
			name:        "invalid did method",
			issuerID:    "did:invalid:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			keyID:       "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:     fmt.Errorf("invalid did:key method: invalid"),
			wantSuccess: false,
		},
		{
			name:        "invalid did",
			issuerID:    "did:invalid",
			keyID:       "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:     fmt.Errorf("failed to parse DID document: invalid did: did:invalid. Make sure it conforms to the DID syntax: https://w3c.github.io/did-core/#did-syntax"),
			wantSuccess: false,
		},
		{
			name:        "invalid pubkey algorithm",
			issuerID:    "did:key:z6MkpwdnLPAm4apwcrRYQ6fZ3rAcqjLZR4AMk14vimfnozqY",
			keyID:       "#z6MkpwdnLPAm4apwcrRYQ6fZ3rAcqjLZR4AMk14vimfnozqY",
			wantErr:     fmt.Errorf("unsupported key algorithm"),
			wantSuccess: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given an issuerID and keyID", t, func() {
				Convey("When the Secp256k1PubKeyFetcher is called", func() {
					pubKey, err := credential.Secp256k1PubKeyFetcher(test.issuerID, test.keyID)

					Convey("Then the result should be as expected", func() {
						if test.wantErr != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
						} else {
							So(err, ShouldBeNil)
							So(pubKey, ShouldNotBeNil)
						}
					})
				})
			})
		})
	}
}

func TestVDRKeyResolverWithSecp256k1_PublicKeyFetcher(t *testing.T) {
	tests := []struct {
		name           string
		issuerID       string
		keyID          string
		wantErr        error
		wantSuccessAlg string
	}{
		{
			name:           "valid secp256k1 key",
			issuerID:       "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			keyID:          "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:        nil,
			wantSuccessAlg: "EcdsaSecp256k1VerificationKey2019",
		},
		{
			name:           "valid ed25519 key",
			issuerID:       "did:key:z6MkpwdnLPAm4apwcrRYQ6fZ3rAcqjLZR4AMk14vimfnozqY",
			keyID:          "#z6MkpwdnLPAm4apwcrRYQ6fZ3rAcqjLZR4AMk14vimfnozqY",
			wantErr:        nil,
			wantSuccessAlg: "Ed25519VerificationKey2018",
		},
		{
			name:           "invalid did",
			issuerID:       "did:invalid",
			keyID:          "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:        fmt.Errorf("failed to parse DID document: invalid did: did:invalid. Make sure it conforms to the DID syntax: https://w3c.github.io/did-core/#did-syntax"),
			wantSuccessAlg: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given secp256k1 fetcher and vdr key resolver", t, func() {
				secp256k1Fetcher := credential.Secp256k1PubKeyFetcher
				vdrKeyResolver := credential.NewVDRKeyResolverWithSecp256k1(secp256k1Fetcher)

				Convey("When the PublicKeyFetcher is called", func() {
					pubKey, err := vdrKeyResolver.PublicKeyFetcher(test.issuerID, test.keyID)

					Convey("Then the result should be as expected", func() {
						if test.wantErr != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
						} else {
							So(err, ShouldBeNil)
							So(pubKey, ShouldNotBeNil)
							So(pubKey.Type, ShouldEqual, test.wantSuccessAlg)
						}
					})
				})
			})
		})
	}
}
