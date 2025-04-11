package dataverse

import (
	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6"
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
		"axone1xt4ahzz2x8hpkc0tk6ekte9x6crw4w6u0r67cyt3kz9syh24pd7scvlt2w",
		"axone1xa8wemfrzq03tkwqxnv9lun7rceec7wuhh8x3qjgxkaaj5fl50zsmj8u0n",
		dataverseClient,
		cognitariumClient,
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
			"",
			"",
			dataverseClient,
			cognitariumClient,
			lawStoneFactory,
		},
		txClient: client,
		txConfig: txConfig,
		signer:   signer,
	}
}

var GetCognitariumAddr = getCognitariumAddr
