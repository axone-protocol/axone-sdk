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

func TestClient_NewClient(t *testing.T) {
	tests := []struct {
		name         string
		grpcAddr     string
		contractAddr string
		wantErr      error
	}{
		{
			name:         "should call get cognitarium address with invalid grpc address",
			grpcAddr:     "invalid",
			contractAddr: "did:key:zQ3shuwMJ",
			wantErr:      fmt.Errorf("failed to get cognitarium address"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a gRPC addr and contract addr", t, func() {
				Convey("When Client is created", func() {
					client, err := dataverse.NewClient(context.Background(), test.grpcAddr, test.contractAddr)

					Convey("The client should be created or return an error", func() {
						So(err.Error(), ShouldContainSubstring, test.wantErr.Error())
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

func Test_GetCognitariumAddr(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		err        error
		wantErr    error
		wantResult string
	}{
		{
			name:       "valid dataverse client",
			address:    "addr",
			err:        nil,
			wantErr:    nil,
			wantResult: "addr",
		},
		{
			name:       "dataverse client return error",
			address:    "",
			err:        fmt.Errorf("error"),
			wantErr:    fmt.Errorf("error"),
			wantResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a mocked dataverse client", t, func() {
				controller := gomock.NewController(t)
				mockDataverse := testutil.NewMockDataverseQueryClient(controller)
				if test.err != nil {
					mockDataverse.EXPECT().
						Dataverse(gomock.Any(), gomock.Any()).
						Return(nil, test.err).
						Times(1)
				} else {
					mockDataverse.EXPECT().
						Dataverse(gomock.Any(), gomock.Any()).
						Return(&dvschema.DataverseResponse{TriplestoreAddress: dvschema.Addr(test.address)}, nil).
						Times(1)
				}

				Convey("When getCognitariumAddr is called", func() {
					addr, err := dataverse.GetCognitariumAddr(context.Background(), mockDataverse)

					Convey("Then should return expected result or error", func() {
						if test.err == nil {
							So(addr, ShouldEqual, test.wantResult)
							So(err, ShouldBeNil)
						} else {
							So(addr, ShouldEqual, "")
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
						}
					})
				})
			})
		})
	}
}
