package auth

import (
	"context"
	"fmt"
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
	dvClient dataverse.Client
	govAddr  string
}

func NewProxy(govAddr string, dvClient dataverse.Client) Proxy {
	return &authProxy{
		dvClient: dvClient,
		govAddr:  govAddr,
	}
}

func (a *authProxy) Authenticate(ctx context.Context, credential []byte) (*Identity, error) {
	// parse credential
	// verify signature
	// get authorized actions from governance, ex:
	did := "did:key:example"
	res, err := a.dvClient.ExecGov(ctx, a.govAddr, fmt.Sprintf("can(Action,'%s').", did))
	if err != nil {
		return nil, err
	}

	return &Identity{
		DID:               did,
		AuthorizedActions: res.([]string),
	}, nil
}
