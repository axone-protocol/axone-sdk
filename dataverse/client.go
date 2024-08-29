package dataverse

import (
	"context"
	"fmt"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"google.golang.org/grpc"
)

type Client interface {
	GetResourceGovAddr(context.Context, string) (string, error)
	ExecGov(context.Context, string, string) (interface{}, error)
}

type client struct {
	dataverseClient   dvschema.QueryClient
	cognitariumClient cgschema.QueryClient
	cognitariumAddr   string
}

func NewDataverseClient(ctx context.Context,
	dataverseClient dvschema.QueryClient,
	cognitariumClient cgschema.QueryClient,
) (Client, error) {
	cognitariumAddr, err := getCognitariumAddr(ctx, dataverseClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get cognitarium address: %w", err)
	}

	return &client{
		dataverseClient,
		cognitariumClient,
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

	cognitariumAddr, err := getCognitariumAddr(ctx, dataverseClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get cognitarium address: %w", err)
	}

	cognitariumClient, err := cgschema.NewQueryClient(grpcAddr, cognitariumAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cognitarium client: %w", err)
	}

	return &client{
		dataverseClient,
		cognitariumClient,
		cognitariumAddr,
	}, nil
}

func (c *client) GetResourceGovAddr(_ context.Context, resourceDID string) (string, error) {
	query := buildGetResourceGovAddrRequest(resourceDID)
	response, err := c.cognitariumClient.Select(context.Background(), &cgschema.QueryMsg_Select{Query: query})
	if err != nil {
		return "", err
	}

	if len(response.Results.Bindings) != 1 {
		return "", NewDVError(ErrNoResult, nil)
	}

	codeBinding, ok := response.Results.Bindings[0]["code"]
	if !ok {
		return "", NewDVError(ErrVarNotFound, nil)
	}
	code, ok := codeBinding.ValueType.(cgschema.URI)
	if !ok {
		return "", NewDVError(ErrType, fmt.Errorf("expected URI, got %T", codeBinding.ValueType))
	}
	return string(*code.Value.Full), nil
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
