package credential

import (
	"fmt"
	"time"

	"github.com/hyperledger/aries-framework-go/component/models/ld/proof"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2020"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
)

type Claim interface {
	From(vc *verifiable.Credential) error
}

type Parser[T Claim] interface {
	ParseSigned(raw []byte) (T, error)
}

type DefaultParser struct {
	documentLoader ld.DocumentLoader
}

func NewDefaultParser(documentLoader ld.DocumentLoader) *DefaultParser {
	return &DefaultParser{documentLoader: documentLoader}
}

func (cp *DefaultParser) Parse(raw []byte) (*verifiable.Credential, error) {
	vc, err := verifiable.ParseCredential(
		raw,
		verifiable.WithJSONLDValidation(),
		verifiable.WithPublicKeyFetcher(NewVDRKeyResolverWithSecp256k1(Secp256k1PubKeyFetcher).PublicKeyFetcher),
		verifiable.WithJSONLDDocumentLoader(cp.documentLoader),
	)
	if err != nil {
		return nil, NewVCError(ErrParse, err)
	}
	return vc, nil
}

func (cp *DefaultParser) parseSigned(raw []byte) (*verifiable.Credential, error) {
	vc, err := verifiable.ParseCredential(
		raw,
		verifiable.WithJSONLDValidation(),
		verifiable.WithPublicKeyFetcher(NewVDRKeyResolverWithSecp256k1(Secp256k1PubKeyFetcher).PublicKeyFetcher),
		verifiable.WithEmbeddedSignatureSuites(
			ed25519signature2018.New(suite.WithVerifier(ed25519signature2018.NewPublicKeyVerifier())),
			ed25519signature2020.New(suite.WithVerifier(ed25519signature2020.NewPublicKeyVerifier())),
			ecdsasecp256k1signature2019.New(suite.WithVerifier(ecdsasecp256k1signature2019.NewPublicKeyVerifier())),
		),
		verifiable.WithJSONLDDocumentLoader(cp.documentLoader),
	)
	if err != nil {
		return nil, NewVCError(ErrParse, err)
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
		return nil, NewVCError(ErrInvalidProof, err)
	}

	return vc, nil
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
		return nil, err
	}
	return pf, nil
}
