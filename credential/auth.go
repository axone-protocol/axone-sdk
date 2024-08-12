package credential

import (
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
)

const ClaimToService = "toService"

var _ Claim = (*AuthClaim)(nil)

type AuthClaim struct {
	ID        string
	ToService string
}

func (ac *AuthClaim) From(vc *verifiable.Credential) error {
	claims, ok := vc.Subject.([]verifiable.Subject)
	if !ok {
		return fmt.Errorf("malformed vc subject")
	}

	if len(claims) != 1 {
		return fmt.Errorf("expected a single vc claim")
	}

	toService, err := extractCustomStringClaim(&claims[0], ClaimToService)
	if err != nil {
		return err
	}

	ac.ID = claims[0].ID
	ac.ToService = toService

	return nil
}

var _ Parser[*AuthClaim] = (*AuthParser)(nil)

type AuthParser struct{}

func NewAuthParser() *AuthParser {
	return &AuthParser{}
}

func (ap *AuthParser) ParseSigned(raw []byte) (*AuthClaim, error) {
	cred, err := parseWithSignVerification(raw)

    authClaim := &AuthClaim{}
	err = authClaim.From(cred)
	if err != nil {
		return nil, fmt.Errorf("malformed auth claim: %w", err)
	}

	return authClaim, nil
}
