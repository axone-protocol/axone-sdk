package dataverse

import (
	"context"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6"
)

// Info is a struct that contains information about the dataverse.
type Info struct {
	// DataverseAddress is the address of the dataverse smart contract instance.
	DataverseAddress string
	// Name is the name of the dataverse
	DataverseName string
}

// CognitariumInfo holds information about the cognitarium instance.
type CognitariumInfo struct {
	// Address of the cognitarium smart contract.
	Address string
	// Owner (admin) of the cognitarium.
	Owner string
	// Stat holds basic statistics.
	Stat CognitariumStat
}

// CognitariumStat contains statistics about the cognitarium.
type CognitariumStat struct {
	// Total size of triples in bytes (Uint128).
	ByteSize string
	// Total number of IRI namespaces (Uint128).
	NamespaceCount string
	// Total number of triples (Uint128).
	TripleCount string
}

func (c *queryClient) DataverseInfo(ctx context.Context) (*Info, error) {
	query := dvschema.QueryMsg_Dataverse{}
	resp, err := c.dataverseClient.Dataverse(ctx, &query)
	if err != nil {
		return nil, err
	}

	return &Info{
		DataverseAddress: c.dataverseContractAddr,
		DataverseName:    resp.Name,
	}, nil
}

func (c *queryClient) CognitariumInfo(ctx context.Context) (*CognitariumInfo, error) {
	query := cgschema.QueryMsg_Store{}
	resp, err := c.cognitariumClient.Store(ctx, &query)
	if err != nil {
		return nil, err
	}

	return &CognitariumInfo{
		Address: c.cognitariumContractAddr,
		Owner:   resp.Owner,
		Stat: CognitariumStat{
			ByteSize:       string(resp.Stat.ByteSize),
			NamespaceCount: string(resp.Stat.NamespaceCount),
			TripleCount:    string(resp.Stat.TripleCount),
		},
	}, nil
}
