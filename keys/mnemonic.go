package keys

import (
	"github.com/axone-protocol/axoned/v9/x/logic/util"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	k "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Keyring = &Key{}

type Key struct {
	privKey  types.PrivKey
	did      string
	didKeyID string
	addr     string
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

	// TODO: Maybe make this configurable...
	sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

	return &Key{
		privKey:  pkey,
		did:      did,
		didKeyID: didKeyID,
		addr:     sdk.AccAddress(pkey.PubKey().Address()).String(),
	}, nil
}

func (k *Key) PubKey() types.PubKey {
	return k.privKey.PubKey()
}

func (k *Key) Sign(msg []byte) ([]byte, error) {
	return k.privKey.Sign(msg)
}

func (k *Key) Alg() string {
	return "secp256k1"
}

func (k *Key) DID() string {
	return k.did
}

func (k *Key) DIDKeyID() string {
	return k.didKeyID
}

func (k *Key) Addr() string {
	return k.addr
}

func parseMnemonic(mnemonic string) (types.PrivKey, error) {
	algo, err := k.NewSigningAlgoFromString("secp256k1", k.SigningAlgoList{hd.Secp256k1})
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
