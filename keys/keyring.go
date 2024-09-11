package keys

import (
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

// Keyring defines the interface for a keyring that can sign messages.
type Keyring interface {
	Sign(msg []byte) ([]byte, error)
	PubKey() types.PubKey
	Alg() string
	// DID return the DID of the key.
	DID() string
	// DIDKeyID returns the DID key ID of the key.
	DIDKeyID() string
	// Addr returns the bech32 address of the key.
	Addr() string
}
