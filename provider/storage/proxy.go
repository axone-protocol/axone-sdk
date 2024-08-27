package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/axone-protocol/axone-sdk/auth"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/axone-protocol/axone-sdk/keys"
)

const (
	readAction  = "read"
	storeAction = "store"
)

type Proxy struct {
	key       *keys.Key
	dvClient  dataverse.Client
	authProxy auth.Proxy

	// given a resource id return its stream
	readFn func(context.Context, string) (io.Reader, error)
	// store the given resource given its id and stream
	storeFn func(context.Context, string, io.Reader) error
}

func NewProxy(
	ctx context.Context,
	key *keys.Key,
	serviceID string,
	dvClient dataverse.Client,
	authParser credential.Parser[*credential.AuthClaim],
	readFn func(context.Context, string) (io.Reader, error),
	storeFn func(context.Context, string, io.Reader) error,
) (*Proxy, error) {
	gov, err := dvClient.GetResourceGovAddr(ctx, key.DID)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		key:       key,
		dvClient:  dvClient,
		authProxy: auth.NewProxy(gov, serviceID, dvClient, authParser),
		readFn:    readFn,
		storeFn:   storeFn,
	}, nil
}

func (p *Proxy) Authenticate(ctx context.Context, credential []byte) (*auth.Identity, error) {
	return p.authProxy.Authenticate(ctx, credential)
}

func (p *Proxy) Read(ctx context.Context, id *auth.Identity, resourceID string) (io.Reader, error) {
	// check authenticated identity resolved authorized actions
	if !id.Can(readAction) {
		return nil, errors.New("unauthorized")
	}

	// get resource gov addr
	govAddr, err := p.dvClient.GetResourceGovAddr(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	// exec resource gov
	_, err = p.dvClient.ExecGov(ctx, govAddr, fmt.Sprintf("can('%s','%s').", readAction, p.key.DID))
	if err != nil {
		return nil, err
	}

	// fetch resource data
	return p.readFn(ctx, resourceID)
}

func (p *Proxy) Store(ctx context.Context, id *auth.Identity, resourceID string, src io.Reader) (io.Reader, error) {
	// check authenticated identity resolved authorized actions
	if !id.Can(storeAction) {
		return nil, errors.New("unauthorized")
	}

	if err := p.storeFn(ctx, resourceID, src); err != nil {
		return nil, err
	}

	return bytes.NewReader([]byte("publication VC")), nil
}
