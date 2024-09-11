package tx

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name          string
		opts          []Option
		wantGasLimit  uint64
		wantFeeAmount types.Coins
		wantMemo      string
		wantMsgs      []types.Msg
	}{
		{
			name: "without option",
			opts: nil,
		},
		{
			name: "with gas limit",
			opts: []Option{
				WithGasLimit(1000),
			},
			wantGasLimit: 1000,
		},
		{
			name: "with fee amount",
			opts: []Option{
				WithFeeAmount(types.NewCoins(types.NewInt64Coin("uaxone", 1000))),
			},
			wantFeeAmount: types.NewCoins(types.NewInt64Coin("uaxone", 1000)),
		},
		{
			name: "with memo",
			opts: []Option{
				WithMemo("memo"),
			},
			wantMemo: "memo",
		},
		{
			name: "with msgs",
			opts: []Option{
				WithMsgs(&wasmtypes.MsgInstantiateContract{}),
			},
			wantMsgs: []types.Msg{
				&wasmtypes.MsgInstantiateContract{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a transaction", t, func() {
				txConfig, err := MakeDefaultTxConfig()
				if err != nil {
					So(err, ShouldBeNil)
				}

				tx := NewTransaction(txConfig, test.opts...)

				Convey("Then the transaction should be created", func() {
					So(tx.(*transaction).gasLimit, ShouldEqual, test.wantGasLimit)
					So(tx.(*transaction).feeAmount, ShouldResemble, test.wantFeeAmount)
					So(tx.(*transaction).memo, ShouldEqual, test.wantMemo)
					So(tx.(*transaction).msgs, ShouldResemble, test.wantMsgs)
				})
			})
		})
	}
}
