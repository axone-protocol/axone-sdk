package dataverse_test

import (
	"context"
	"testing"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/credential/template"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/axone-protocol/axone-sdk/tx"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestClient_SubmitClaims(t *testing.T) {
	tests := []struct {
		name       string
		credential *verifiable.Credential
		wantErr    error
	}{
		{
			name:       "valid credential",
			credential: generateVC(),
			wantErr:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a client", t, func() {
				controller := gomock.NewController(t)
				defer controller.Finish()
				txConfig, err := tx.MakeDefaultTxConfig()
				if err != nil {
					So(err, ShouldBeNil)
				}

				mockDataverseClient := testutil.NewMockDataverseQueryClient(controller)
				mockCognitarium := testutil.NewMockCognitariumQueryClient(controller)
				mockTxClient := testutil.NewMockTxClient(controller)
				mockKeyring := testutil.NewMockKeyring(controller)

				mockKeyring.EXPECT().Addr().Return("addr").AnyTimes()
				mockTxClient.EXPECT().SendTx(gomock.Any(), gomock.Any()).Return(&types.TxResponse{}, nil)

				client := dataverse.NewDataverseTxClient(
					mockDataverseClient,
					mockCognitarium,
					nil,
					mockTxClient,
					txConfig,
					mockKeyring,
				)

				Convey("When SubmitClaims is called", func() {
					r, err := client.SubmitClaims(context.Background(), test.credential)

					Convey("Then should return expected error", func() {
						if test.wantErr == nil {
							So(err, ShouldBeNil)
							So(r, ShouldNotBeNil)
						} else {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(r, ShouldBeNil)
						}
					})
				})
			})
		})
	}
}

func generateVC() *verifiable.Credential {
	loader, _ := testutil.MockDocumentLoader()
	vc, err := credential.New(
		template.NewGovernance("datasetID", "addr"),
		credential.WithParser(credential.NewDefaultParser(loader)),
	).
		Generate()
	if err != nil {
		panic(err)
	}
	return vc
}
