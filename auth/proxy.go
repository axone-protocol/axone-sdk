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
	dvClient   dataverse.Client
	govAddr    string
	authParser credential.Parser[*credential.AuthClaim]
}

func NewProxy(govAddr string, dvClient dataverse.Client, authParser credential.Parser[*credential.AuthClaim]) Proxy {
	return &authProxy{
		dvClient:   dvClient,
		govAddr:    govAddr,
		authParser: authParser,
	}
}

func (a *authProxy) Authenticate(ctx context.Context, credential []byte) (*Identity, error) {
	authClaim, err := a.authParser.ParseSigned(credential)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credential: %w", err)
	}

	// TODO: get authorized actions from governance, ex:
	did := "did:key:example"
	res, err := a.dvClient.ExecGov(ctx, a.govAddr, fmt.Sprintf("can(Action,'%s').", did))
	if err != nil {
		return nil, err
	}

	return &Identity{
		DID:               authClaim.ID,
		AuthorizedActions: res.([]string),
	}, nil
}
