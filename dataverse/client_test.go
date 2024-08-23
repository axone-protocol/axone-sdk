package dataverse_test

import (
	"context"
	"fmt"
	schema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestClient_GetGovAddr(t *testing.T) {
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
			wantErr:     fmt.Errorf("failed to get governance address: %w", fmt.Errorf("error")),
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

				client := dataverse.NewDataverseClient(mockClient)

				Convey("When GetGovAddr is called", func() {
					resp, err := client.GetGovAddr(context.Background())

					Convey("Then it should return the governance address", func() {
						So(err, ShouldEqual, test.wantErr)
						So(resp, ShouldEqual, test.wantAddr)
					})
				})
			})
		})
	}
}
