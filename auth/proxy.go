package auth

import (
	"context"
	"fmt"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/dataverse"
)

// Proxy acts as the entrypoint of a service and is responsible for authenticating any identity willing to conduct some
// actions against the underlying service.
// It authenticates Decentralized Identities based on a provided Verifiable Credential and resolving allowed authorized
// actions for this identity based on on-chain rules.
// It is not responsible or aware of the communication protocol, which means it only returns information on the identity
// if authentic and won't for example issue a JWT token, this is out of its scope.
type Proxy interface {
	// Authenticate verifies the authenticity and integrity of the provided credential before resolving on-chain
	// authorized actions with the proxied service by querying the service's governance.
	Authenticate(ctx context.Context, credential []byte) (*Identity, error)
}

type authProxy struct {
	dvClient   dataverse.QueryClient
	authParser credential.Parser[*credential.AuthClaim]
	govAddr    string
	serviceID  string
}

// NewProxy creates a new Proxy instance using the given service identifier and on-chain governance address (i.e. the
// law-stone smart contract instance carrying its rules).
func NewProxy(govAddr, serviceID string,
	dvClient dataverse.QueryClient,
	authParser credential.Parser[*credential.AuthClaim],
) Proxy {
	return &authProxy{
		dvClient:   dvClient,
		authParser: authParser,
		govAddr:    govAddr,
		serviceID:  serviceID,
	}
}

// Authenticate verifies the authenticity and integrity of the provided credential before resolving on-chain
// authorized actions with the proxied service by querying its governance.
func (a *authProxy) Authenticate(ctx context.Context, credential []byte) (*Identity, error) {
	authClaim, err := a.authParser.ParseSigned(credential)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credential: %w", err)
	}

	if authClaim.ToService != a.serviceID {
		return nil, fmt.Errorf("credential not intended for this service: `%s` (target: `%s`)", a.serviceID, authClaim.ToService)
	}

	actions, err := a.dvClient.AskGovPermittedActions(ctx, a.govAddr, authClaim.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query governance for permitted actions: %w", err)
	}

	return &Identity{
		DID:               authClaim.ID,
		AuthorizedActions: actions,
	}, nil
}
