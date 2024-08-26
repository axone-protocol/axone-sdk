package dataverse_test

import (
	"context"
	"fmt"
	"testing"

	schema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestClient_NewDataverseClient(t *testing.T) {
	tests := []struct {
		name        string
		returnedErr error
		wantErr     error
		wantAddr    string
	}{
		{
			name:        "receive an cognitarium address",
			returnedErr: nil,
			wantErr:     nil,
			wantAddr:    "addr",
		},
		{
			name:        "receive an error",
			returnedErr: fmt.Errorf("error"),
			wantErr:     fmt.Errorf("failed to get cognitarium address: %w", fmt.Errorf("error")),
			wantAddr:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mocked dataverse client", t, func() {
				controller := gomock.NewController(t)
				defer controller.Finish()

				mockClient := testutil.NewMockDataverseQueryClient(controller)
				mockClient.EXPECT().
					Dataverse(gomock.Any(), gomock.Any()).
					Return(
						&schema.DataverseResponse{
							TriplestoreAddress: schema.Addr(test.wantAddr),
						},
						test.returnedErr,
					).
					Times(1)

				Convey("When Client is created", func() {
					client, err := dataverse.NewDataverseClient(context.Background(), mockClient)

					Convey("Then the client should be created if no error on dataverse client", func() {
						So(err, ShouldEqual, test.wantErr)
						if test.wantErr == nil {
							So(client, ShouldNotBeNil)
						} else {
							So(client, ShouldBeNil)
						}
					})
				})
			})
		})
	}
}
