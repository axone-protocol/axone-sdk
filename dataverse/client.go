package dataverse

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/axone-protocol/axone-sdk/keys"
	"github.com/axone-protocol/axone-sdk/tx"
	"github.com/axone-protocol/axone-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
	dvschema "github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5"
	lsschema "github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v5"
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

    SubmitClaims(ctx context.Context, credential *verifiable.Credential, signer keys.Keyring) error
}

type LawStoneFactory func(string) (lsschema.QueryClient, error)

type client struct {
	dataverseClient   dvschema.QueryClient
	cognitariumClient cgschema.QueryClient
	txClient          tx.Client
	contractAddr      string
	lawStoneFactory   LawStoneFactory
}

func NewClient(ctx context.Context,
	grpcAddr, contractAddr string,
	txClient tx.Client,
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
		txClient,
		contractAddr,
        func(addr string) (lsschema.QueryClient, error) {
            return lsschema.NewQueryClient(grpcAddr, addr, opts...)
        },
	}, nil
}

func (c *client) SubmitClaims(ctx context.Context, vc *verifiable.Credential, signer keys.Keyring) error {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/n-quads"
	options.Explicit = true

	vcRaw, err := vc.MarshalJSON()
	if err != nil {
		return err
	}

	var vcJSON interface{}
	err = json.Unmarshal(vcRaw, &vcJSON)
	if err != nil {
		return err
	}
	rdf, err := proc.ToRDF(vcJSON, options)
	if err != nil {
		return err
	}
	fmt.Printf("rdf: %s\n", rdf)

	msg, err := json.Marshal(map[string]interface{}{
		"submit_claims": map[string]interface{}{
			"claims": base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", rdf))),
		},
	})
	if err != nil {
		return err
	}
	addr, err := sdk.AccAddressFromHexUnsafe(signer.PubKey().Address().String())
	if err != nil {
		return err
	}
	msgExec := &wasmtypes.MsgExecuteContract{
		Sender:   addr.String(),
		Contract: c.contractAddr,
		Msg:      msg,
		Funds:    nil,
	}

	sendTx, err := c.txClient.SendTx(types.Context{}.WithContext(ctx).WithChainID("axone-localnet"), []sdk.Msg{msgExec},
		tx.WithSigner(signer),
		tx.WithGasLimit(2000000),
	)
	if err != nil {
		return err
	}

	fmt.Printf("sendTx: %v\n", sendTx)
	return nil
}

func getCognitariumAddr(ctx context.Context, dvClient dvschema.QueryClient) (string, error) {
	query := dvschema.QueryMsg_Dataverse{}
	resp, err := dvClient.Dataverse(ctx, &query)
	if err != nil {
		return "", err
	}

	return string(resp.TriplestoreAddress), nil
}
