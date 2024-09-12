package dataverse

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/axone-protocol/axone-sdk/tx"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/piprate/json-gold/ld"
)

func (t *txClient) SubmitClaims(ctx context.Context, vc *verifiable.Credential) error {
	rdf, err := credentialToRDF(vc)
	if err != nil {
		return NewDVError(ErrConvertRDF, err)
	}

	msg, err := json.Marshal(map[string]interface{}{
		"submit_claims": map[string]interface{}{
			"claims": base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s", rdf))),
		},
	})
	if err != nil {
		return NewDVError(ErrMarshalJSON, err)
	}

	msgExec := &wasmtypes.MsgExecuteContract{
		Sender:   t.signer.Addr(),
		Contract: t.contractAddr,
		Msg:      msg,
		Funds:    nil,
	}

	_, err = t.txClient.SendTx(ctx, tx.NewTransaction(t.txConfig,
		tx.WithMsgs(msgExec),
		tx.WithSigner(t.signer),
		tx.WithGasLimit(2000000),
	))
	if err != nil {
		return NewDVError(ErrSendTx, err)
	}

	return nil
}

func credentialToRDF(vc *verifiable.Credential) (interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/n-quads"

	vcRaw, err := vc.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var vcJSON interface{}
	err = json.Unmarshal(vcRaw, &vcJSON)
	if err != nil {
		return nil, err
	}

	return proc.ToRDF(vcJSON, options)
}
