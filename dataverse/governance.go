package dataverse

import (
	"context"
	"fmt"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
)

func (c *client) GetResourceGovAddr(ctx context.Context, resourceDID string) (string, error) {
	query := buildGetResourceGovAddrRequest(resourceDID)
	response, err := c.cognitariumClient.Select(ctx, &cgschema.QueryMsg_Select{Query: query})
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
