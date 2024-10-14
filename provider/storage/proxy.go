// Package storage provides the core logic needed to implement storage services in the Axone protocol.
package storage

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/credential/template"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/keys"
	"github.com/piprate/json-gold/ld"
)

const (
	readAction  = "read"
	storeAction = "store"
)

// Proxy serves as an authentication and authorization proxy of an Axone storage service.
// It is responsible for authenticating and authorizing identities before operations as read and store resources.
// The specific logic of reading and storing resources is delegated to the readFn and storeFn functions.
type Proxy struct {
	key       keys.Keyring
	baseURL   string
	dvClient  dataverse.QueryClient
	authProxy auth.Proxy
	vcParser  *credential.DefaultParser

	// given a resource id return its stream
	readFn func(context.Context, string) (io.Reader, error)
	// store the given resource given its id and stream
	storeFn func(context.Context, string, io.Reader) error
}

// NewProxy creates a new Proxy instance, using the provided service DID to retrieve its governance (i.e. law-stone smart
// contract address) on the dataverse.
func NewProxy(
	ctx context.Context,
	key keys.Keyring,
	baseURL string,
	dvClient dataverse.QueryClient,
	documentLoader ld.DocumentLoader,
	readFn func(context.Context, string) (io.Reader, error),
	storeFn func(context.Context, string, io.Reader) error,
) (*Proxy, error) {
	gov, err := dvClient.GetResourceGovAddr(ctx, key.DID())
	if err != nil {
		return nil, err
	}

	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	return &Proxy{
		key:       key,
		baseURL:   baseURL,
		dvClient:  dvClient,
		authProxy: auth.NewProxy(gov, key.DID(), dvClient, credential.NewAuthParser(documentLoader)),
		vcParser:  credential.NewDefaultParser(documentLoader),
		readFn:    readFn,
		storeFn:   storeFn,
	}, nil
}

// Authenticate performs the authentication of an identity from a verifiable credential returning its resolved auth.Identity.
func (p *Proxy) Authenticate(ctx context.Context, credential []byte) (*auth.Identity, error) {
	return p.authProxy.Authenticate(ctx, credential)
}

// Read reads a resource identified by its resourceID, returning its stream if the identity is authorized to do so.
//
// The identity is authorized to read a resource if both the proxied service's governance and the requested resource's
// governance allows it. To check the proxied service's governance it uses the set of resolved permissions at authentication.
// To check the requested resource's governance it retrieves it from the dataverse before querying it.
func (p *Proxy) Read(ctx context.Context, id *auth.Identity, resourceID string) (io.Reader, error) {
	if !id.Can(readAction) {
		return nil, errors.New("unauthorized")
	}

	govAddr, err := p.dvClient.GetResourceGovAddr(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	ok, err := p.dvClient.AskGovTellAction(ctx, govAddr, id.DID, readAction)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("unauthorized")
	}

	return p.readFn(ctx, resourceID)
}

// Store stores a resource identified by its resourceID, returning its publication credential. This publication credential
// is a verifiable credential that attests the publication of the resource by the proxied service, it is expected to be
// submitted to the dataverse in order to reference the resource.
//
// The identity is authorized to read a resource if the proxied service's governance allows it, it uses the set of resolved
// permissions at authentication to do so.
func (p *Proxy) Store(ctx context.Context, id *auth.Identity, resourceID string, src io.Reader) (io.Reader, error) {
	if !id.Can(storeAction) {
		return nil, errors.New("unauthorized")
	}

	if err := p.storeFn(ctx, resourceID, src); err != nil {
		return nil, err
	}

	vc, err := credential.New(
		template.NewPublication(resourceID, p.baseURL+resourceID, p.key.DID()),
		credential.WithParser(p.vcParser),
		credential.WithSigner(p.key),
	).Generate()
	if err != nil {
		return nil, err
	}

	raw, err := vc.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(raw), nil
}
