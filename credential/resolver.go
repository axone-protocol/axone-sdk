package credential

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk/jwksupport"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"
)

var _ verifiable.PublicKeyFetcher = Secp256k1PubKeyFetcher

var Secp256k1PubKeyFetcher = resolve

var ErrKeyAlgorithm = fmt.Errorf("unsupported key algorithm")

func resolve(issuerID, _ string) (*verifier.PublicKey, error) {
	parsed, err := did.Parse(issuerID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DID document: %w", err)
	}

	if parsed.Method != "key" {
		return nil, fmt.Errorf("invalid did:key method: %s", parsed.Method)
	}

	pubKeyBytes, code, err := fingerprint.PubKeyFromFingerprint(parsed.MethodSpecificID)
	if err != nil {
		return nil, fmt.Errorf("failed to get key fingerprint: %w", err)
	}

	if code != 0xe7 { // secp256k1 algorithm code
		return nil, ErrKeyAlgorithm
	}

	pubKey, err := secp.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	j, err := jwksupport.JWKFromKey(&ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     pubKey.X(),
		Y:     pubKey.Y(),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating JWK: %w", err)
	}

	return &verifier.PublicKey{
		Type: "EcdsaSecp256k1VerificationKey2019",
		JWK:  j,
	}, nil
}

// VDRKeyResolverWithSecp256k1 is a VDR key resolver including a secp256k1 public key fetcher as is
// not available in the default VDRKeyResolver.
// It's a hack to include this algorithm in the resolver.
type VDRKeyResolverWithSecp256k1 struct {
	vdrKeyResolver         *verifiable.VDRKeyResolver
	secp256k1PubKeyFetcher verifiable.PublicKeyFetcher
}

func NewVDRKeyResolverWithSecp256k1(secp256k1PubKeyFetcher verifiable.PublicKeyFetcher) *VDRKeyResolverWithSecp256k1 {
	return &VDRKeyResolverWithSecp256k1{
		verifiable.NewVDRKeyResolver(vdr.New(vdr.WithVDR(key.New()))),
		secp256k1PubKeyFetcher,
	}
}

func (r *VDRKeyResolverWithSecp256k1) PublicKeyFetcher(issuerDID, keyID string) (*verifier.PublicKey, error) {
	pubKey, err := r.secp256k1PubKeyFetcher(issuerDID, keyID)
	if err != nil && !errors.Is(err, ErrKeyAlgorithm) {
		return nil, err
	} else if !errors.Is(err, ErrKeyAlgorithm) {
		return pubKey, nil
	}

	return r.vdrKeyResolver.PublicKeyFetcher()(issuerDID, keyID)
}
