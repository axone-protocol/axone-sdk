package dataverse

import (
	"context"
	"fmt"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"google.golang.org/grpc"
)

type Client interface {
	GetGovAddr(context.Context) (string, error)
	ExecGov(context.Context, string, string) (interface{}, error)
}

type client struct {
	dataverseClient dvschema.QueryClient
}

func NewDataverseClient(dataverseClient dvschema.QueryClient) Client {
	return &client{
		dataverseClient,
	}
}

func NewClient(grpcAddr, contractAddr string, opts ...grpc.DialOption) (Client, error) {
	dataverseClient, err := dvschema.NewQueryClient(grpcAddr, contractAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataverse client: %w", err)
	}

	return &client{
		dataverseClient,
	}, nil
}

func (c *client) GetGovAddr(ctx context.Context) (string, error) {
	query := dvschema.QueryMsg_Dataverse{}
	resp, err := c.dataverseClient.Dataverse(ctx, &query)
	if err != nil {
		return "", fmt.Errorf("failed to get governance address: %w", err)
	}

	return string(resp.TriplestoreAddress), nil
}

func (c *client) ExecGov(ctx context.Context, addr string, method string) (interface{}, error) {
	panic("not implemented")
}
