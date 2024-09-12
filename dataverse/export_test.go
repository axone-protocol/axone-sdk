package dataverse

import (
	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	"github.com/axone-protocol/axone-sdk/keys"
	"github.com/axone-protocol/axone-sdk/tx"
	"github.com/cosmos/cosmos-sdk/client"
)

func NewDataverseQueryClient(
	dataverseClient dvschema.QueryClient,
	cognitariumClient cgschema.QueryClient,
	lawStoneFactory LawStoneFactory,
) QueryClient {
	return &queryClient{
		dataverseClient,
		cognitariumClient,
		"",
		lawStoneFactory,
	}
}

func NewDataverseTxClient(
	dataverseClient dvschema.QueryClient,
	cognitariumClient cgschema.QueryClient,
	lawStoneFactory LawStoneFactory,
	client tx.Client,
	txConfig client.TxConfig,
	signer keys.Keyring,
) TxClient {
	return &txClient{
		queryClient: &queryClient{
			dataverseClient,
			cognitariumClient,
			"",
			lawStoneFactory,
		},
		txClient: client,
		txConfig: txConfig,
		signer:   signer,
	}
}

var GetCognitariumAddr = getCognitariumAddr
