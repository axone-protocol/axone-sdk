package auth_test

import (
	ctx "context"
	"fmt"
	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAuthProxy_Authenticate(t *testing.T) {
	tests := []struct {
		name             string
		credential       []byte
		expectedIdentity *auth.Identity
		wantErr          error
	}{
		{
			name:       "valid token",
			credential: []byte("valid"),
			expectedIdentity: &auth.Identity{
				DID:               "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
				AuthorizedActions: nil,
			},
		},
		{
			name:       "valid token",
			credential: nil,
			wantErr:    fmt.Errorf("failed to parse credential: nil"),
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
				mockDataverse.EXPECT().ExecGov(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string(nil), nil).MaxTimes(1)

				aProxy := auth.NewProxy(
					"did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
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
