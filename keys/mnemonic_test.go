package keys

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewKeyFromMnemonic(t *testing.T) {
	tests := []struct {
		name       string
		mnemonic   string
		wantDID    string
		wantDIDKey string
		wantAddr   string
		wantErr    error
	}{
		{
			name:       "valid mnemonic",
			mnemonic:   "code ceiling reduce repeat unfold intact cloud marriage nut remove illegal eternal pool frame mask rate buzz vintage pulp suggest loan faint snake spoon",
			wantDID:    "did:key:zQ3shVMdXcC6eGDC1UGHDzzvtrZVwmqHtYaAW6BDKqPNR569S",
			wantDIDKey: "did:key:zQ3shVMdXcC6eGDC1UGHDzzvtrZVwmqHtYaAW6BDKqPNR569S#zQ3shVMdXcC6eGDC1UGHDzzvtrZVwmqHtYaAW6BDKqPNR569S",
			wantAddr:   "axone14u8n76zahep9xkfr9gc3zxv5c7rf3x8wx3fdjl",
			wantErr:    nil,
		},
		{
			name:     "invalid mnemonic",
			mnemonic: "invalid",
			wantErr:  fmt.Errorf("Invalid mnemonic"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mnemonic", t, func() {
				Convey("When NewKeyFromMnemonic is called", func() {
					key, err := NewKeyFromMnemonic(test.mnemonic)
					Convey("Then the key should be created", func() {
						if test.wantErr != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(key, ShouldBeNil)
						} else {
							So(err, ShouldBeNil)
							So(key, ShouldNotBeNil)
							So(key.privKey, ShouldNotBeNil)
							So(key.DID(), ShouldEqual, test.wantDID)
							So(key.DIDKeyID(), ShouldEqual, test.wantDIDKey)
							So(key.Addr(), ShouldEqual, test.wantAddr)
							So(key.PubKey(), ShouldNotBeNil)
						}
					})
				})
			})
		})
	}
}
