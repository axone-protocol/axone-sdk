package dataverse_test

import (
	"context"
	"fmt"
	"testing"

	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestDataverseInfo(t *testing.T) {
	tests := []struct {
		name           string
		queryResponse  *dvschema.DataverseResponse
		queryError     error
		expectedInfo   *dataverse.Info
		expectedErrMsg string
	}{
		{
			name: "ReturnsCorrectInfo",
			queryResponse: &dvschema.DataverseResponse{
				Name:               "TestDataverse",
				TriplestoreAddress: "triplestore-address",
			},
			queryError: nil,
			expectedInfo: &dataverse.Info{
				DataverseAddress: "axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
				DataverseName:    "TestDataverse",
			},
			expectedErrMsg: "",
		},
		{
			name:           "ReturnsErrorOnQueryFailure",
			queryResponse:  nil,
			queryError:     fmt.Errorf("query error"),
			expectedInfo:   nil,
			expectedErrMsg: "query error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mocked dataverse client", t, func() {
				controller := gomock.NewController(t)
				Reset(controller.Finish)

				mockDataverseClient := testutil.NewMockDataverseQueryClient(controller)
				mockDataverseClient.EXPECT().
					Dataverse(gomock.Any(), gomock.Any()).
					Return(test.queryResponse, test.queryError).
					Times(1)

				mockCognitariumClient := testutil.NewMockCognitariumQueryClient(controller)

				client := dataverse.NewDataverseQueryClient(
					mockDataverseClient,
					mockCognitariumClient,
					nil,
				)

				Convey("When DataverseInfo is called", func() {
					info, err := client.DataverseInfo(context.Background())

					Convey("Then the expected result should be returned", func() {
						if test.expectedErrMsg == "" {
							So(err, ShouldBeNil)
							So(info, ShouldResemble, test.expectedInfo)
						} else {
							So(err.Error(), ShouldEqual, test.expectedErrMsg)
							So(info, ShouldBeNil)
						}
					})
				})
			})
		})
	}
}
