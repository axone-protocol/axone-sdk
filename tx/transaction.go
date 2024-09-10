package tx

import (
    "github.com/axone-protocol/axone-sdk/keys"
    sdkclient "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/types"
)

type Transaction struct {
    builder sdkclient.TxBuilder
    signer  keys.Keyring
}

type Option func(*Transaction)

func New(builder sdkclient.TxBuilder, msgs []types.Msg, opts ...Option) *Transaction {
    tx := &Transaction{
        builder: builder,
    }
    err := builder.SetMsgs(msgs...)
    if err != nil {
        panic(err)
    }

    for _, opt := range opts {
        opt(tx)
    }
    return tx
}

func WithMemo(memo string) Option {
    return func(tx *Transaction) {
        tx.builder.SetMemo(memo)
    }
}

func WithGasLimit(limit uint64) Option {
    return func(tx *Transaction) {
        tx.builder.SetGasLimit(limit)
    }
}

func WithFeeAmount(amount types.Coins) Option {
    return func(tx *Transaction) {
        tx.builder.SetFeeAmount(amount)
    }
}

func WithSigner(signer keys.Keyring) Option {
    return func(tx *Transaction) {
        tx.signer = signer
    }
}
