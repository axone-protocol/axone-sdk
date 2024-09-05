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

type Proxy struct {
	key       keys.Keyring
	baseURL   string
	dvClient  dataverse.Client
	authProxy auth.Proxy
	vcParser  *credential.DefaultParser

	// given a resource id return its stream
	readFn func(context.Context, string) (io.Reader, error)
	// store the given resource given its id and stream
	storeFn func(context.Context, string, io.Reader) error
}

func NewProxy(
	ctx context.Context,
	key keys.Keyring,
	baseURL string,
	dvClient dataverse.Client,
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

func (p *Proxy) Authenticate(ctx context.Context, credential []byte) (*auth.Identity, error) {
	return p.authProxy.Authenticate(ctx, credential)
}

func (p *Proxy) Read(ctx context.Context, id *auth.Identity, resourceID string) (io.Reader, error) {
	if !id.Can(readAction) {
		return nil, errors.New("unauthorized")
	}

	govAddr, err := p.dvClient.GetResourceGovAddr(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	ok, err := p.dvClient.AskGovTellAction(ctx, govAddr, p.key.DID(), readAction)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("unauthorized")
	}

	return p.readFn(ctx, resourceID)
}

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
		credential.WithSigner(p.key, p.key.DID()),
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
