package credential

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/hyperledger/aries-framework-go/component/models/ld/proof"
	"time"

	"github.com/btcsuite/btcd/btcec"
	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk/jwksupport"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"
	"github.com/piprate/json-gold/ld"
)

type Claim interface {
	From(vc *verifiable.Credential) error
}

type Parser[T Claim] interface {
	ParseSigned(raw []byte) (T, error)
}

type credentialParser struct {
	documentLoader ld.DocumentLoader
}

func (cp *credentialParser) parseSigned(raw []byte) (*verifiable.Credential, error) {
	publicKeyFetcher := verifiable.NewVDRKeyResolver(vdr.New(vdr.WithVDR(key.New()))).PublicKeyFetcher()

	vc, err := verifiable.ParseCredential(
		raw,
		verifiable.WithJSONLDValidation(),
		verifiable.WithPublicKeyFetcher(func(issuerID, keyID string) (*verifier.PublicKey, error) {
			// HACK: as the publicKeyFetcher doesn't manage `EcdsaSecp256k1VerificationKey2019` as verification method
			// we got to manage it ourselves.
			pubKey, err := mayResolveSecp256k1PubKey(issuerID, keyID)
			if err != nil {
				return nil, err
			}

			if pubKey != nil {
				return pubKey, nil
			}

			return publicKeyFetcher(issuerID, keyID)
		}),
		verifiable.WithEmbeddedSignatureSuites(
			ed25519signature2018.New(suite.WithVerifier(ed25519signature2018.NewPublicKeyVerifier())),
			ed25519signature2020.New(suite.WithVerifier(ed25519signature2020.NewPublicKeyVerifier())),
			ecdsasecp256k1signature2019.New(suite.WithVerifier(ecdsasecp256k1signature2019.NewPublicKeyVerifier())),
		),
		verifiable.WithJSONLDDocumentLoader(cp.documentLoader),
	)
	if err != nil {
		return nil, err
	}
	return withCheck(vc)
}

func withCheck(vc *verifiable.Credential) (*verifiable.Credential, error) {
	if vc.Expired != nil && time.Now().After(vc.Expired.Time) {
		return nil, NewVCError(ErrExpired, fmt.Errorf("%s", vc.Expired.Time))
	}

	if vc.Issued != nil && time.Now().Before(vc.Issued.Time) {
		return nil, NewVCError(ErrIssued, fmt.Errorf("%s", vc.Issued.Time))
	}

	if _, err := extractProof(vc); err != nil {
		return nil, err
	}

	return vc, nil
}

// Hack helper to resolve a did key as a `EcdsaSecp256k1VerificationKey2019` if it is, as the `PublicKeyFetcher` we use
// doesn't.
func mayResolveSecp256k1PubKey(issuerID, keyID string) (*verifier.PublicKey, error) {
	issuerDid, err := did.Parse(issuerID)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to parse DID document: %w", err)
	}

	if issuerDid.Method != "key" {
		return nil, fmt.Errorf("vdr Read: invalid did:key method: %s", issuerDid.Method)
	}

	pubKeyBytes, code, err := fingerprint.PubKeyFromFingerprint(issuerDid.MethodSpecificID)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to get key fingerPrint: %w", err)
	}

	if code == 0xe7 && fmt.Sprintf("#%s", issuerDid.MethodSpecificID) == keyID {
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
	return nil, fmt.Errorf("pub:key vdr Read: invalid fingerprint")
}

func extractCustomStringClaim(claim *verifiable.Subject, key string) (string, error) {
	field, ok := claim.CustomFields[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}

	strField, ok := field.(string)
	if !ok {
		return "", fmt.Errorf("key '%s' is not a string", key)
	}
	return strField, nil
}

func extractProof(vc *verifiable.Credential) (*proof.Proof, error) {
	if len(vc.Proofs) == 0 {
		return nil, NewVCError(ErrMissingProof, nil)
	}

	pf, err := proof.NewProof(vc.Proofs[0])
	if err != nil {
		return nil, NewVCError(ErrInvalidProof, err)
	}
	return pf, nil
}
