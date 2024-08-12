package credential_test

import (
	"fmt"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/piprate/json-gold/ld"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestAuthParser_ParseSigned(t *testing.T) {
	tests := []struct {
		name      string
		serviceId string
		file      string
		wantErr   error
		result    *credential.AuthClaim
	}{
		{
			name:      "valid credential",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/valid.jsonld",
			wantErr:   nil,
			result: &credential.AuthClaim{
				ID:        "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
				ToService: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			},
		},
		{
			name:      "credential not signed",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_not-signed.jsonld",
			wantErr:   fmt.Errorf("missing verifiable credential proof"),
			result:    nil,
		},
		{
			name:      "credential not signed",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_wrong-signature.jsonld",
			wantErr:   fmt.Errorf("decode new credential: check embedded proof: check linked data proof: ecdsa: invalid signature"),
			result:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a credential", t, func() {
				raw, err := os.ReadFile(test.file)
				So(err, ShouldBeNil)

				parser := credential.NewAuthParser(test.serviceId, ld.NewDefaultDocumentLoader(nil))

				Convey("When parsing the credential", func() {
					authClaim, err := parser.ParseSigned(raw)

					Convey("Then the result should be as expected", func() {
						if err != nil {
							So(err.Error(), ShouldResemble, test.wantErr.Error())
						} else {
							So(err, ShouldBeNil)
						}
						So(authClaim, ShouldResemble, test.result)
					})
				})
			})
		})
	}
}
