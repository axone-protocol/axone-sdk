package credential

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
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

type AuthParser struct {
	*credentialParser
	ServiceID string
}

func NewAuthParser(serviceID string, documentLoader ld.DocumentLoader) *AuthParser {
	return &AuthParser{
		credentialParser: &credentialParser{documentLoader: documentLoader},
		ServiceID:        serviceID,
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
		return nil, fmt.Errorf("malformed auth claim: %w", err)
	}

	if authClaim.ToService != ap.ServiceID {
		return nil, fmt.Errorf("auth claim target doesn't match current service id: %s (target: %s)", ap.ServiceID, authClaim.ToService)
	}

	if cred.Issuer.ID != authClaim.ID {
		return nil, fmt.Errorf("auth claim subject differs from issuer (subject: %s, issuer: %s)", authClaim.ID, cred.Issuer.ID)
	}
	return authClaim, nil
}
