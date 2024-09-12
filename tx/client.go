package tx

import (
	"context"
	"cosmossdk.io/x/tx/signing"
	"fmt"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectype "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
)

type Client interface {
	SendTx(ctx context.Context, transaction Transaction) (*sdk.TxResponse, error)
}

type client struct {
	authClient      authtypes.QueryClient
	txServiceClient tx.ServiceClient
	chainID         string
}

func NewClient(authClient authtypes.QueryClient, txServiceClient tx.ServiceClient, chainID string) Client {
	return &client{
		authClient,
		txServiceClient,
		chainID,
	}
}

func (c *client) SendTx(ctx context.Context, transaction Transaction) (*sdk.TxResponse, error) {
	accNum, accSeq, err := c.getAccountNumberSequence(ctx, transaction.Sender())
	if err != nil {
		return nil, fmt.Errorf("failed to get account number and sequence: %w", err)
	}

	txEncoded, err := transaction.GetSignedTx(ctx, accNum, accSeq, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed build a signed tx: %w", err)
	}

	resp, err := c.txServiceClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{TxBytes: txEncoded, Mode: tx.BroadcastMode_BROADCAST_MODE_SYNC})
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast tx: %w", err)
	}
	return resp.TxResponse, nil
}

func (c *client) getAccountNumberSequence(ctx context.Context, addr string) (uint64, uint64, error) {
	resp, err := c.authClient.Account(ctx, &authtypes.QueryAccountRequest{Address: addr})
	if err != nil {
		return 0, 0, err
	}

	var account authtypes.BaseAccount
	if err := account.Unmarshal(resp.GetAccount().Value); err != nil {
		return 0, 0, err
	}

	return account.AccountNumber, account.Sequence, nil
}

func MakeDefaultTxConfig() (sdkclient.TxConfig, error) {
	interfaceRegistry, err := codectype.NewInterfaceRegistryWithOptions(codectype.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return authtx.NewTxConfig(
		codec.NewProtoCodec(interfaceRegistry),
		authtx.DefaultSignModes,
	), nil
}
