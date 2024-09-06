package tx

import (
	"cosmossdk.io/x/tx/signing"
	"github.com/axone-protocol/axone-sdk/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectype "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
)

type Client interface {
	SendTx(ctx types.Context, msgs []sdk.Msg, opts ...Option) (*sdk.TxResponse, error)
}

type client struct {
	authClient authtypes.QueryClient
	txConfig   sdkclient.TxConfig
	txClient   tx.ServiceClient
}

func NewClient(authClient authtypes.QueryClient, txClient tx.ServiceClient) (Client, error) {
	config, err := makeTxConfig()
	if err != nil {
		return nil, err
	}

	return &client{
		authClient,
		config,
		txClient,
	}, nil
}

func (c *client) SendTx(ctx types.Context, msgs []sdk.Msg, opts ...Option) (*sdk.TxResponse, error) {
	txBuilder := c.txConfig.NewTxBuilder()
	transaction := New(txBuilder, msgs, opts...)

	accNum, accSeq, err := c.getAccountNumberSequence(ctx, transaction.signer.PubKey().Address().String())

	signerData := authsigning.SignerData{
		Address:       transaction.signer.PubKey().Address().String(),
		ChainID:       ctx.ChainID(),
		AccountNumber: accNum,
		Sequence:      accSeq,
		PubKey:        transaction.signer.PubKey(),
	}

	sig := signingtypes.SignatureV2{
		PubKey: transaction.signer.PubKey(),
		Data: &signingtypes.SingleSignatureData{
			SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: signerData.Sequence,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	bytesToSign, err := authsigning.GetSignBytesAdapter(ctx,
		c.txConfig.SignModeHandler(),
		signingtypes.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	sigBytes, err := transaction.signer.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}

	sig.Data = &signingtypes.SingleSignatureData{
		SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	encodeTx, err := c.txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	resp, err := c.txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{TxBytes: encodeTx})
	if err != nil {
		return nil, err
	}
	return resp.TxResponse, nil
}

func (c *client) getAccountNumberSequence(ctx types.Context, addr string) (uint64, uint64, error) {
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

func makeTxConfig() (sdkclient.TxConfig, error) {
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
