package keys

import (
	"github.com/cosmos/cosmos-sdk/crypto/types"
)

type Keyring interface {
	Sign(msg []byte) ([]byte, error)
	PubKey() types.PubKey
	Alg() string
	DID() string
}
