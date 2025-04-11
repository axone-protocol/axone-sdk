package dataverse

import (
	"context"

	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6"
)

// Info is a struct that contains information about the dataverse.
type Info struct {
	// DataverseAddress is the address of the dataverse smart contract instance.
	DataverseAddress string
	// Name is the name of the dataverse
	DataverseName string
}

func (c *queryClient) DataverseInfo(ctx context.Context) (*Info, error) {
	query := dvschema.QueryMsg_Dataverse{}
	resp, err := c.dataverseClient.Dataverse(ctx, &query)
	if err != nil {
		return nil, err
	}

	return &Info{
		DataverseAddress: c.contractAddr,
		DataverseName:    resp.Name,
	}, nil
}
