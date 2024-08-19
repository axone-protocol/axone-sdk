package jwt_test

import (
    "bytes"
    "context"
    "fmt"
    "github.com/axone-protocol/axone-sdk/auth"
    "github.com/axone-protocol/axone-sdk/auth/jwt"
    "github.com/axone-protocol/axone-sdk/testutil"
    . "github.com/smartystreets/goconvey/convey"
    "go.uber.org/mock/gomock"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestFactory_HTTPAuthHandler(t *testing.T) {
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
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            Convey("Given a JWT factory and mocked auth proxy on mocked http server", t, func() {
                factory := jwt.NewFactory([]byte("secret"), "issuer", 5*time.Second)

                controller := gomock.NewController(t)
                defer controller.Finish()

                mockAuthProxy := testutil.NewMockProxy(controller)
                handler := factory.HTTPAuthHandler(mockAuthProxy)

                mockAuthProxy.EXPECT().Authenticate(gomock.Any(), gomock.Any()).
                    DoAndReturn(func(ctx context.Context, body []byte) (interface{}, error) {
                        if string(body) == "valid" {
                            return test.identity, nil
                        }
                        return nil, fmt.Errorf("authentication failed")
                    }).
                    Times(1)

                Convey("When the HTTPAuthHandler is called", func() {
                    req, err := http.NewRequest("POST", "/", bytes.NewBuffer(test.body))
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
