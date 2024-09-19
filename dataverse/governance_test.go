package dataverse_test

import (
	"context"
	"fmt"
	"testing"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	lsschema "github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v5"
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
			name:        "ask for good did response with addr in uri",
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
									Value: cgschema.IRI{Full: toAddress(cgschema.IRI_Full("contract:law-stone:foo"))},
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

				client := dataverse.NewDataverseQueryClient(
					mockDataverseClient,
					mockCognitarium,
					nil,
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

func TestClient_AskGovPermittedActions(t *testing.T) {
	tests := []struct {
		name          string
		addr          string
		did           string
		response      *lsschema.AskResponse
		responseError error
		wantErr       error
		wantResult    []string
	}{
		{
			name:          "law stone client new error",
			addr:          "error",
			did:           "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response:      nil,
			responseError: nil,
			wantErr:       fmt.Errorf("failed to create law-stone client: error"),
			wantResult:    nil,
		},
		{
			name:          "law stone client ask error",
			addr:          "foo",
			did:           "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response:      nil,
			responseError: fmt.Errorf("error"),
			wantErr:       fmt.Errorf("failed to query law-stone contract: error"),
			wantResult:    nil,
		},
		{
			name: "no results in response",
			addr: "foo",
			did:  "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    nil,
		},
		{
			name: "no substitutions in response",
			addr: "foo",
			did:  "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    nil,
		},
		{
			name: "single action in response",
			addr: "foo",
			did:  "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{
								{
									Expression: "['read']",
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    []string{"read"},
		},
		{
			name: "multiple actions in response",
			addr: "foo",
			did:  "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{
								{
									Expression: "['read','store']",
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    []string{"read", "store"},
		},
		{
			name: "quoted/unquoted actions in response",
			addr: "foo",
			did:  "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{
								{
									Expression: "['read',store]",
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    []string{"read", "store"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mocked dataverse client", t, func() {
				controller := gomock.NewController(t)
				defer controller.Finish()

				lawStoneMock := testutil.NewMockLawStoneQueryClient(controller)
				if test.addr != "error" {
					lawStoneMock.EXPECT().
						Ask(gomock.Any(), gomock.Eq(&lsschema.QueryMsg_Ask{Query: fmt.Sprintf("tell_permitted_actions('%s',Actions).", test.did)})).
						Return(test.response, test.responseError).
						Times(1)
				}

				client := dataverse.NewDataverseQueryClient(
					testutil.NewMockDataverseQueryClient(controller),
					testutil.NewMockCognitariumQueryClient(controller),
					func(addr string) (lsschema.QueryClient, error) {
						if addr == "error" {
							return nil, fmt.Errorf("error")
						}

						return lawStoneMock, nil
					},
				)

				Convey("When AskGovPermittedActions is called", func() {
					actions, err := client.AskGovPermittedActions(context.Background(), test.addr, test.did)

					Convey("Then the permitted actions should be returned", func() {
						if test.wantErr == nil {
							So(err, ShouldBeNil)
							So(actions, ShouldResemble, test.wantResult)
						} else {
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(actions, ShouldBeNil)
						}
					})
				})
			})
		})
	}
}

func TestClient_AskGovTellAction(t *testing.T) {
	tests := []struct {
		name          string
		addr          string
		did           string
		action        string
		response      *lsschema.AskResponse
		responseError error
		wantErr       error
		wantResult    bool
	}{
		{
			name:          "law stone client new error",
			addr:          "error",
			did:           "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			action:        "read",
			response:      nil,
			responseError: nil,
			wantErr:       fmt.Errorf("failed to create law-stone client: error"),
			wantResult:    false,
		},
		{
			name:          "law stone client ask error",
			addr:          "foo",
			did:           "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			action:        "read",
			response:      nil,
			responseError: fmt.Errorf("error"),
			wantErr:       fmt.Errorf("failed to query law-stone contract: error"),
			wantResult:    false,
		},
		{
			name:   "no results in response",
			addr:   "foo",
			did:    "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			action: "read",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    false,
		},
		{
			name:   "no substitutions in response",
			addr:   "foo",
			did:    "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			action: "read",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    false,
		},
		{
			name:   "permitted",
			addr:   "foo",
			did:    "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			action: "read",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{
								{
									Expression: "permitted",
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    true,
		},
		{
			name:   "prohibited",
			addr:   "foo",
			did:    "did:key:zQ3shuwMJWYXRi64qiGojsV9bPN6Dtugz5YFM2ESPtkaNxTZ5",
			action: "store",
			response: &lsschema.AskResponse{
				Answer: &lsschema.Answer{
					Results: []lsschema.Result{
						{
							Substitutions: []lsschema.Substitution{
								{
									Expression: "prohibited",
								},
							},
						},
					},
				},
			},
			responseError: nil,
			wantErr:       nil,
			wantResult:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mocked dataverse client", t, func() {
				controller := gomock.NewController(t)
				defer controller.Finish()

				lawStoneMock := testutil.NewMockLawStoneQueryClient(controller)
				if test.addr != "error" {
					lawStoneMock.EXPECT().
						Ask(gomock.Any(), gomock.Eq(&lsschema.QueryMsg_Ask{Query: fmt.Sprintf("tell('%s','%s',Result,_).", test.did, test.action)})).
						Return(test.response, test.responseError).
						Times(1)
				}

				client := dataverse.NewDataverseQueryClient(
					testutil.NewMockDataverseQueryClient(controller),
					testutil.NewMockCognitariumQueryClient(controller),
					func(addr string) (lsschema.QueryClient, error) {
						if addr == "error" {
							return nil, fmt.Errorf("error")
						}

						return lawStoneMock, nil
					},
				)

				Convey("When AskGovTellAction is called", func() {
					result, err := client.AskGovTellAction(context.Background(), test.addr, test.did, test.action)

					Convey("Then it should indicate if the action is allowed", func() {
						if test.wantErr == nil {
							So(err, ShouldBeNil)
						} else {
							So(err.Error(), ShouldEqual, test.wantErr.Error())
						}
						So(result, ShouldResemble, test.wantResult)
					})
				})
			})
		})
	}
}
