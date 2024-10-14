package jwt_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/auth/jwt"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestIssuer_HTTPAuthHandler(t *testing.T) {
	tests := []struct {
		name                 string
		body                 []byte
		identity             *auth.Identity
		expectedStatus       int
		expectedBodyContains string
	}{
		{
			name: "valid request",
			body: []byte(`valid`),
			identity: &auth.Identity{
				DID:               "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
				AuthorizedActions: nil,
			},
			expectedStatus:       200,
			expectedBodyContains: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		},
		{
			name:                 "invalid authentication",
			body:                 []byte(`invalid`),
			identity:             nil,
			expectedStatus:       401,
			expectedBodyContains: "failed to authenticate: authentication failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a JWT issuer and mocked auth proxy on mocked http server", t, func() {
				issuer := jwt.NewIssuer(nil, "issuer", 5*time.Second)

				controller := gomock.NewController(t)
				defer controller.Finish()

				mockAuthProxy := testutil.NewMockProxy(controller)
				handler := issuer.HTTPAuthHandler(mockAuthProxy)

				mockAuthProxy.EXPECT().Authenticate(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, body []byte) (interface{}, error) {
						if body != nil && string(body) == "valid" {
							return test.identity, nil
						}
						return nil, fmt.Errorf("authentication failed")
					}).
					Times(1)

				Convey("When the HTTPAuthHandler is called", func() {
					req, err := http.NewRequest("POST", "/", bytes.NewBuffer(test.body)) //nolint:noctx
					So(err, ShouldBeNil)

					recorder := httptest.NewRecorder()
					handler.ServeHTTP(recorder, req)

					Convey("Then the response should match the expected status and body", func() {
						So(recorder.Code, ShouldEqual, test.expectedStatus)
						So(recorder.Body.String(), ShouldContainSubstring, test.expectedBodyContains)
					})
				})
			})
		})
	}
}

func TestIssuer_VerifyHTTPMiddleware(t *testing.T) {
	// Generate a valid token for testing purpose
	token, err := jwt.NewIssuer([]byte("secret"), "issuer", 5*time.Second).IssueJWT(&auth.Identity{
		DID:               "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
		AuthorizedActions: nil,
	})
	if err != nil {
		t.Fatalf("failed to issue a test jwt token: %v", err)
	}

	tests := []struct {
		name                 string
		authHeader           string
		expectedIdentity     *auth.Identity
		expectedStatus       int
		expectedBodyContains string
	}{
		{
			name:       "valid token",
			authHeader: fmt.Sprintf("Bearer %s", token),
			expectedIdentity: &auth.Identity{
				DID: "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
			},
			expectedStatus:       200,
			expectedBodyContains: "success",
		},
		{
			name:                 "invalid token",
			authHeader:           "Bearer invalid_token",
			expectedStatus:       401,
			expectedBodyContains: "token contains an invalid number of segments",
		},
		{
			name:                 "missing token",
			authHeader:           "",
			expectedStatus:       401,
			expectedBodyContains: "couldn't find bearer token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a JWT issuer and mocked authenticated handler", t, func() {
				issuer := jwt.NewIssuer([]byte("secret"), "issuer", 5*time.Second)

				mockHandler := func(id *auth.Identity, w http.ResponseWriter, _ *http.Request) {
					So(id, ShouldResemble, test.expectedIdentity)
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte("success"))
					So(err, ShouldBeNil)
				}

				middleware := issuer.VerifyHTTPMiddleware(mockHandler)

				Convey("When the VerifyHTTPMiddleware is called", func() {
					req, err := http.NewRequest("GET", "/", nil) //nolint:noctx
					So(err, ShouldBeNil)
					if test.authHeader != "" {
						req.Header.Set("Authorization", test.authHeader)
					}

					recorder := httptest.NewRecorder()
					middleware.ServeHTTP(recorder, req)

					Convey("Then the response should match the expected status and body", func() {
						So(recorder.Code, ShouldEqual, test.expectedStatus)
						So(recorder.Body.String(), ShouldContainSubstring, test.expectedBodyContains)
					})
				})
			})
		})
	}
}
