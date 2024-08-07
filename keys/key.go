package keys

import (
	"github.com/axone-protocol/axoned/v9/x/logic/util"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Key struct {
	privKey  types.PrivKey
	DID      string
	DIDKeyID string
	Addr     string
}

func NewKeyFromMnemonic(mnemonic string) (*Key, error) {
	pkey, err := parseMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return NewKeyFromPrivKey(pkey)
}

func NewKeyFromPrivKey(pkey types.PrivKey) (*Key, error) {
	did, err := util.CreateDIDKeyByPubKey(pkey.PubKey())
	if err != nil {
		return nil, err
	}

	didKeyID, err := util.CreateDIDKeyIDByPubKey(pkey.PubKey())
	if err != nil {
		return nil, err
	}

	return &Key{
		privKey:  pkey,
		DID:      did,
		DIDKeyID: didKeyID,
		Addr:     sdk.AccAddress(pkey.PubKey().Address()).String(),
	}, nil
}

func (k *Key) PubKey() types.PubKey {
	return k.privKey.PubKey()
}

func (k *Key) Sign(msg []byte) ([]byte, error) {
	return k.privKey.Sign(msg)
}

func (k *Key) Alg() string {
	return "unknown"
}

func parseMnemonic(mnemonic string) (types.PrivKey, error) {
	algo, err := keyring.NewSigningAlgoFromString("secp256k1", keyring.SigningAlgoList{hd.Secp256k1})
	if err != nil {
		return nil, err
	}

	hdPath := hd.CreateHDPath(118, 0, 0).String()

	derivedPriv, err := algo.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return nil, err
	}

	return algo.Generate()(derivedPriv), nil
}
