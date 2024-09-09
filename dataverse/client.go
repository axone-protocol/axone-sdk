package dataverse

import (
	"context"
	"fmt"
	lsschema "github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v5"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"google.golang.org/grpc"
)

type Client interface {
	// GetResourceGovAddr returns the governance address of a resource.
	// It queries the cognitarium to get the governance address (law-stone contract address)
	// of a resource. The resource is identified by its DID.
	GetResourceGovAddr(context.Context, string) (string, error)

	// AskGovPermittedActions returns the permitted actions for a resource identified by its DID.
	// It queries the law-stone contract to get the permitted actions for a resource using the following predicate:
	// ```prolog
	// tell_permitted_actions(DID, Actions).
	// ```
	AskGovPermittedActions(context.Context, string, string) ([]string, error)

	// AskGovTellAction queries the law-stone contract to check if a given action is permitted for a resource.
	// It uses the following predicate:
	// ```prolog
	// tell(DID, Action, Result, Evidence).
	// ```
	// The function returns true if Result is 'permitted', false otherwise.
	AskGovTellAction(context.Context, string, string, string) (bool, error)
}

type LawStoneFactory func(string) (lsschema.QueryClient, error)

type client struct {
	dataverseClient   dvschema.QueryClient
	cognitariumClient cgschema.QueryClient
	lawStoneFactory   LawStoneFactory
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
		func(addr string) (lsschema.QueryClient, error) {
			return lsschema.NewQueryClient(grpcAddr, addr, opts...)
		},
	}, nil
}

func getCognitariumAddr(ctx context.Context, dvClient dvschema.QueryClient) (string, error) {
	query := dvschema.QueryMsg_Dataverse{}
	resp, err := dvClient.Dataverse(ctx, &query)
	if err != nil {
		return "", err
	}

	return string(resp.TriplestoreAddress), nil
}
