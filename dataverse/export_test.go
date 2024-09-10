package dataverse

import (
	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
)

func NewDataverseClient(
	dataverseClient dvschema.QueryClient,
	cognitariumClient cgschema.QueryClient,
	lawStoneFactory LawStoneFactory,
) Client {
	return &client{
		dataverseClient,
		cognitariumClient,
		lawStoneFactory,
	}
}

var GetCognitariumAddr = getCognitariumAddr
