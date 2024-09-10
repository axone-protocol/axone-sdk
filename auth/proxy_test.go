//nolint:lll
package auth_test

import (
	ctx "context"
	"fmt"
	"testing"

	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestAuthProxy_Authenticate(t *testing.T) {
	tests := []struct {
		name             string
		credential       []byte
		serviceID        string
		expectedIdentity *auth.Identity
		wantErr          error
	}{
		{
			name:       "valid credential",
			credential: []byte("valid"),
			serviceID:  "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			expectedIdentity: &auth.Identity{
				DID:               "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
				AuthorizedActions: nil,
			},
		},
		{
			name:       "check error returned by authParser",
			credential: nil,
			serviceID:  "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			wantErr:    fmt.Errorf("failed to parse credential: nil"),
		},
		{
			name:       "valid credential that target wrong service",
			credential: []byte("valid"),
			serviceID:  "did:key:wrong",
			wantErr:    fmt.Errorf("credential not intended for this service: `did:key:wrong` (target: `did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz`)"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mock auth parser and AuthProxy", t, func() {
				controller := gomock.NewController(t)
				mockAuthParser := testutil.NewMockParser[*credential.AuthClaim](controller)

				if test.credential != nil {
					mockAuthParser.EXPECT().
						ParseSigned([]byte("valid")).
						Return(&credential.AuthClaim{
							ID:        "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
							ToService: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
						}, nil).
						Times(1)
				} else {
					mockAuthParser.EXPECT().ParseSigned(nil).Return(nil, fmt.Errorf("nil")).Times(1)
				}

				mockDataverse := testutil.NewMockClient(controller)
				mockDataverse.EXPECT().AskGovPermittedActions(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string(nil), nil).MaxTimes(1)

				aProxy := auth.NewProxy(
					"did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
					test.serviceID,
					mockDataverse,
					mockAuthParser)

				Convey("When authenticating", func() {
					identity, err := aProxy.Authenticate(ctx.Background(), test.credential)

					Convey("Then the result should be as expected", func() {
						if test.wantErr != nil {
							So(err.Error(), ShouldEqual, test.wantErr.Error())
						} else {
							So(err, ShouldBeNil)
						}
						So(identity, ShouldResemble, test.expectedIdentity)
					})
				})
			})
		})
	}
}
