package tx

import (
	"context"
	"errors"
	"fmt"
	"github.com/axone-protocol/axone-sdk/keys"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type Transaction interface {
	Sender() string
	GetSignedTx(ctx context.Context, accNum, accSeq uint64, chainID string) ([]byte, error)
}

var _ Transaction = &transaction{}

type transaction struct {
	txConfig  sdkclient.TxConfig
	txBuilder sdkclient.TxBuilder

	msgs      []types.Msg
	signer    keys.Keyring
	memo      string
	gasLimit  uint64
	feeAmount types.Coins
}

type Option func(*transaction)

func NewTransaction(txConfig sdkclient.TxConfig, opts ...Option) Transaction {
	tx := &transaction{
		txConfig: txConfig,
	}
	for _, opt := range opts {
		opt(tx)
	}
	return tx
}

func WithMsgs(msgs ...types.Msg) Option {
	return func(tx *transaction) {
		tx.msgs = msgs
	}
}

func WithMemo(memo string) Option {
	return func(tx *transaction) {
		tx.memo = memo
	}
}

func WithGasLimit(limit uint64) Option {
	return func(tx *transaction) {
		tx.gasLimit = limit
	}
}

func WithFeeAmount(amount types.Coins) Option {
	return func(tx *transaction) {
		tx.feeAmount = amount
	}
}

func WithSigner(signer keys.Keyring) Option {
	return func(tx *transaction) {
		tx.signer = signer
	}
}

func (t *transaction) sign(ctx context.Context,
	accNum, accSeq uint64,
	chainID string) error {
	if t.signer == nil {
		return errors.New("no signer provided")
	}

	signerData := authsigning.SignerData{
		Address:       t.signer.Addr(),
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      accSeq,
		PubKey:        t.signer.PubKey(),
	}

	sig := signingtypes.SignatureV2{
		PubKey: t.signer.PubKey(),
		Data: &signingtypes.SingleSignatureData{
			SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: signerData.Sequence,
	}

	if err := t.txBuilder.SetSignatures(sig); err != nil {
		return err
	}

	bytesToSign, err := authsigning.GetSignBytesAdapter(ctx,
		t.txConfig.SignModeHandler(),
		signingtypes.SignMode_SIGN_MODE_DIRECT,
		signerData,
		t.txBuilder.GetTx())
	if err != nil {
		return err
	}

	sigBytes, err := t.signer.Sign(bytesToSign)
	if err != nil {
		return err
	}

	sig.Data = &signingtypes.SingleSignatureData{
		SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}

	return t.txBuilder.SetSignatures(sig)
}

func (t *transaction) GetSignedTx(ctx context.Context,
	accNum, accSeq uint64,
	chainID string) ([]byte, error) {

	t.txBuilder = t.txConfig.NewTxBuilder()

	if err := t.txBuilder.SetMsgs(t.msgs...); err != nil {
		return nil, err
	}
	t.txBuilder.SetGasLimit(t.gasLimit)
	t.txBuilder.SetFeeAmount(t.feeAmount)
	t.txBuilder.SetMemo(t.memo)

	if err := t.sign(ctx, accNum, accSeq, chainID); err != nil {
		return nil, fmt.Errorf("could not sign transaction: %w", err)
	}

	return t.txConfig.TxEncoder()(t.txBuilder.GetTx())
}

func (t *transaction) Sender() string {
	return t.signer.Addr()
}
