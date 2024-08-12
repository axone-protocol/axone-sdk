package credential_test

import (
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
						So(err, ShouldResemble, test.wantErr)
						So(authClaim, ShouldResemble, test.result)
					})
				})
			})
		})
	}
}
