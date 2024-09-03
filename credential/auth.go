package credential

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
)

const (
	ClaimToService             = "toService"
	ProofPurposeAuthentication = "authentication"
)

const ErrAuthClaim MessageError = "invalid auth claim"

var _ Claim = (*AuthClaim)(nil)

type AuthClaim struct {
	ID        string
	ToService string
}

func (ac *AuthClaim) From(vc *verifiable.Credential) error {
	claims, ok := vc.Subject.([]verifiable.Subject)
	if !ok {
		return NewVCError(ErrMalformedSubject, nil)
	}

	if len(claims) != 1 {
		return NewVCError(ErrExpectSingleClaim, nil)
	}

	toService, err := extractCustomStringClaim(&claims[0], ClaimToService)
	if err != nil {
		return NewVCError(ErrExtractClaim, err)
	}

	ac.ID = claims[0].ID
	ac.ToService = toService

	return nil
}

var _ Parser[*AuthClaim] = (*AuthParser)(nil)

type AuthParser struct {
	*DefaultParser
}

func NewAuthParser(documentLoader ld.DocumentLoader) *AuthParser {
	return &AuthParser{
		DefaultParser: &DefaultParser{documentLoader: documentLoader},
	}
}

func (ap *AuthParser) ParseSigned(raw []byte) (*AuthClaim, error) {
	cred, err := ap.parseSigned(raw)
	if err != nil {
		return nil, err
	}

	authClaim := &AuthClaim{}
	err = authClaim.From(cred)
	if err != nil {
		return nil, NewVCError(ErrMalformed, err)
	}

	proof, err := extractProof(cred)
	if err != nil {
		return nil, NewVCError(ErrInvalidProof, err)
	}
	if proof.ProofPurpose != ProofPurposeAuthentication {
		return nil, NewVCError(ErrAuthClaim,
			fmt.Errorf("proof purpose not targeting `%s` (proof purpose: `%s`)", ProofPurposeAuthentication, proof.ProofPurpose))
	}

	if cred.Issuer.ID != authClaim.ID {
		return nil, NewVCError(ErrAuthClaim,
			fmt.Errorf("subject differs from issuer (subject: `%s`, issuer: `%s`)", authClaim.ID, cred.Issuer.ID))
	}
	return authClaim, nil
}
