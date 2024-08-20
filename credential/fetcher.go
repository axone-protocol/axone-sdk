package credential

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk/jwksupport"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
)

var _ verifiable.PublicKeyFetcher = Secp256k1PubKeyFetcher

var Secp256k1PubKeyFetcher = resolve

var ErrKeyAlgorithm = fmt.Errorf("unsupported key algorithm")

func resolve(issuerID, keyID string) (*verifier.PublicKey, error) {
	parsed, err := did.Parse(issuerID)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to parse DID document: %w", err)
	}

	if parsed.Method != "key" {
		return nil, fmt.Errorf("vdr Read: invalid did:key method: %s", parsed.Method)
	}

	pubKeyBytes, code, err := fingerprint.PubKeyFromFingerprint(parsed.MethodSpecificID)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to get key fingerPrint: %w", err)
	}

	if code != 0xe7 { // secp256k1 algorithm code
		return nil, ErrKeyAlgorithm
	}

	pubKey, err := secp.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to parse public key: %w", err)
	}
	j, err := jwksupport.JWKFromKey(&ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     pubKey.X(),
		Y:     pubKey.Y(),
	})
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: error creating JWK: %w", err)
	}

	return &verifier.PublicKey{
		Type: "EcdsaSecp256k1VerificationKey2019",
		JWK:  j,
	}, nil
}
