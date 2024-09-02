package dataverse_test

import (
	"context"
	"fmt"
	"testing"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func toAddress[T any](v T) *T {
	return &v
}

func TestClient_GetResourceGovAddr(t *testing.T) {
	tests := []struct {
		name          string
		resourceDID   string
		response      *cgschema.SelectResponse
		responseError error
		wantErr       error
		wantResult    string
	}{
		{
			name:        "ask for good did response",
			resourceDID: "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &cgschema.SelectResponse{
				Head: cgschema.Head{
					Vars: []string{"code"},
				},
				Results: cgschema.Results{
					Bindings: []map[string]cgschema.Value{
						{
							"code": cgschema.Value{
								ValueType: cgschema.URI{
									Type:  "uri",
									Value: cgschema.IRI{Full: toAddress(cgschema.IRI_Full("foo"))},
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    "foo",
		},
		{
			name:          "grpc error",
			resourceDID:   "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response:      nil,
			responseError: fmt.Errorf("gRPC: connection refused"),
			wantErr:       fmt.Errorf("gRPC: connection refused"),
			wantResult:    "",
		},
		{
			name:        "invalid variable binding in response",
			resourceDID: "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &cgschema.SelectResponse{
				Head: cgschema.Head{
					Vars: []string{"code"},
				},
				Results: cgschema.Results{
					Bindings: []map[string]cgschema.Value{
						{
							"invalid": cgschema.Value{},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       dataverse.NewDVError(dataverse.ErrVarNotFound, nil),
			wantResult:    "",
		},
		{
			name:        "no binding in response",
			resourceDID: "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &cgschema.SelectResponse{
				Head: cgschema.Head{
					Vars: []string{"code"},
				},
				Results: cgschema.Results{
					Bindings: []map[string]cgschema.Value{},
				},
			},
			responseError: nil,
			wantErr:       dataverse.NewDVError(dataverse.ErrNoResult, nil),
			wantResult:    "",
		},
		{
			name:        "invalid value type in response",
			resourceDID: "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &cgschema.SelectResponse{
				Head: cgschema.Head{
					Vars: []string{"code"},
				},
				Results: cgschema.Results{
					Bindings: []map[string]cgschema.Value{
						{
							"code": cgschema.Value{
								ValueType: cgschema.BlankNode{
									Type:  "blank_node",
									Value: "foo",
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       dataverse.NewDVError(dataverse.ErrType, fmt.Errorf("expected URI, got %T", cgschema.BlankNode{})),
			wantResult:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mocked dataverse client", t, func() {
				controller := gomock.NewController(t)
				defer controller.Finish()

				mockDataverseClient := testutil.NewMockDataverseQueryClient(controller)

				mockCognitarium := testutil.NewMockCognitariumQueryClient(controller)
				mockCognitarium.
					EXPECT().
					Select(gomock.Any(), gomock.Any()).
					Return(test.response, test.responseError).
					Times(1)

				client := dataverse.NewDataverseClient(
					mockDataverseClient,
					mockCognitarium,
				)

				Convey("When GetResourceGovAddr is called", func() {
					addr, err := client.GetResourceGovAddr(context.Background(), test.resourceDID)

					Convey("Then the resource governance address should be returned", func() {
						if test.wantErr == nil {
							So(err, ShouldBeNil)
							So(addr, ShouldEqual, test.wantResult)
						} else {
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(addr, ShouldEqual, "")
						}
					})
				})
			})
		})
	}
}