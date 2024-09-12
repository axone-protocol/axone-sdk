package tx_test

import (
	"context"
	"fmt"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/axone-protocol/axone-sdk/tx"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdktype "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestClient_SendTx(t *testing.T) {
	acc := &authtypes.BaseAccount{
		AccountNumber: 20,
		Sequence:      19,
	}
	accByte, err := acc.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		shouldAccountErr   error
		shouldSignErr      error
		shouldBroadcastErr error
		wantErr            error
	}{
		{
			name: "success",
		},
		{
			name:             "account error",
			shouldAccountErr: fmt.Errorf("account error"),
			wantErr:          fmt.Errorf("failed to get account number and sequence: account error"),
		},
		{
			name:          "signature error",
			shouldSignErr: fmt.Errorf("signature error"),
			wantErr:       fmt.Errorf("failed build a signed tx: signature error"),
		},
		{
			name:               "broadcast error",
			shouldBroadcastErr: fmt.Errorf("broadcast error"),
			wantErr:            fmt.Errorf("failed to broadcast tx: broadcast error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a client with mocked auth client and service client", t, func() {

				controller := gomock.NewController(t)
				defer controller.Finish()

				mockAuthClient := testutil.NewMockAuthQueryClient(controller)
				mockTxService := testutil.NewMockTxServiceClient(controller)
				mockTransaction := testutil.NewMockTransaction(controller)

				mockTransaction.EXPECT().Sender().Return("axone1").Times(1)

				if test.shouldAccountErr != nil {
					mockAuthClient.EXPECT().
						Account(gomock.Any(), &authtypes.QueryAccountRequest{Address: "axone1"}).
						Return(nil, test.shouldAccountErr)
				} else {
					mockAuthClient.EXPECT().
						Account(gomock.Any(), &authtypes.QueryAccountRequest{Address: "axone1"}).
						Return(&authtypes.QueryAccountResponse{Account: &types.Any{Value: accByte}}, nil)
				}

				if test.shouldSignErr != nil {
					mockTransaction.EXPECT().
						GetSignedTx(gomock.Any(), uint64(20), uint64(19), "chainID").
						Return(nil, test.shouldSignErr)
				} else if test.shouldAccountErr == nil {
					mockTransaction.EXPECT().
						GetSignedTx(gomock.Any(), uint64(20), uint64(19), "chainID").
						Return([]byte("txEncoded"), nil)
				}

				if test.shouldBroadcastErr != nil {
					mockTxService.EXPECT().
						BroadcastTx(gomock.Any(), &sdktx.BroadcastTxRequest{TxBytes: []byte("txEncoded"), Mode: sdktx.BroadcastMode_BROADCAST_MODE_SYNC}).
						Return(nil, test.shouldBroadcastErr)
				} else if test.shouldAccountErr == nil && test.shouldSignErr == nil {
					mockTxService.EXPECT().
						BroadcastTx(gomock.Any(), &sdktx.BroadcastTxRequest{TxBytes: []byte("txEncoded"), Mode: sdktx.BroadcastMode_BROADCAST_MODE_SYNC}).
						Return(&sdktx.BroadcastTxResponse{TxResponse: &sdktype.TxResponse{}}, nil)
				}

				client := tx.NewClient(mockAuthClient, mockTxService, "chainID")

				Convey("When SendTx is called", func() {
					result, err := client.SendTx(context.Background(), mockTransaction)

					Convey("Then it should return an error", func() {
						if test.wantErr != nil {
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(result, ShouldBeNil)
						} else {
							So(err, ShouldBeNil)
							So(result, ShouldNotBeNil)
						}
					})
				})
			})
		})
	}
}
