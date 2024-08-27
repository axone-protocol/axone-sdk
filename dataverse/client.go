package dataverse

import (
	"context"
	"fmt"

	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"google.golang.org/grpc"
)

type Client interface {
	GetResourceGovAddr(context.Context, string) (string, error)
	ExecGov(context.Context, string, string) (interface{}, error)
}

type client struct {
	dataverseClient dvschema.QueryClient
	cognitariumAddr string
}

func NewDataverseClient(ctx context.Context, dataverseClient dvschema.QueryClient) (Client, error) {
	cognitariumAddr, err := getCognitariumAddr(ctx, dataverseClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get cognitarium address: %w", err)
	}

	return &client{
		dataverseClient,
		cognitariumAddr,
	}, nil
}

func NewClient(ctx context.Context,
	grpcAddr, contractAddr string,
	opts ...grpc.DialOption,
) (Client, error) {
	dataverseClient, err := dvschema.NewQueryClient(grpcAddr, contractAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataverse client: %w", err)
	}

	return NewDataverseClient(ctx, dataverseClient)
}

func (c *client) GetResourceGovAddr(_ context.Context, _ string) (string, error) {
	panic("not implemented")
}

func (c *client) ExecGov(_ context.Context, _ string, _ string) (interface{}, error) {
	panic("not implemented")
}

func getCognitariumAddr(ctx context.Context, dvClient dvschema.QueryClient) (string, error) {
	query := dvschema.QueryMsg_Dataverse{}
	resp, err := dvClient.Dataverse(ctx, &query)
	if err != nil {
		return "", err
	}

	return string(resp.TriplestoreAddress), nil
}
