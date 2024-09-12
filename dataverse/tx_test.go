package dataverse_test

import (
	"context"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/credential/template"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/axone-protocol/axone-sdk/tx"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
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

				client := dataverse.NewDataverseTxClient(
					mockDataverseClient,
					mockCognitarium,
					nil,
					mockTxClient,
					txConfig,
					mockKeyring,
				)

				Convey("When SubmitClaims is called", func() {
					err := client.SubmitClaims(context.Background(), test.credential)

					Convey("Then should return expected error", func() {
						if test.wantErr == nil {
							So(err, ShouldBeNil)
						} else {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
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
