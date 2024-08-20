package credential_test

import (
	"fmt"
	"github.com/axone-protocol/axone-sdk/credential"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
			wantErr:     fmt.Errorf("pub:key vdr Read: failed to get key fingerPrint: unknown key encoding"),
			wantSuccess: false,
		},
		{
			name:        "invalid did method",
			issuerID:    "did:invalid:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			keyID:       "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:     fmt.Errorf("vdr Read: invalid did:key method: invalid"),
			wantSuccess: false,
		},
		{
			name:        "invalid did",
			issuerID:    "did:invalid",
			keyID:       "#zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			wantErr:     fmt.Errorf("pub:key vdr Read: failed to parse DID document: invalid did: did:invalid. Make sure it conforms to the DID syntax: https://w3c.github.io/did-core/#did-syntax"),
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
