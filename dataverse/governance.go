package dataverse

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6"
	lsschema "github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v6"
)

func (c *queryClient) GetResourceGovAddr(ctx context.Context, resourceDID string) (string, error) {
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

	codeURI := string(*code.Value.Full)
	addr := codeURI
	if i := strings.LastIndex(string(*code.Value.Full), ":"); i != -1 {
		addr = codeURI[i+1:]
	}

	return addr, nil
}

func (c *queryClient) GovCode(ctx context.Context, addr string) (string, error) {
	gov, err := c.lawStoneFactory(addr)
	if err != nil {
		return "", fmt.Errorf("failed to create law-stone client: %w", err)
	}

	code, err := gov.ProgramCode(ctx, &lsschema.QueryMsg_ProgramCode{})
	if err != nil {
		return "", fmt.Errorf("failed to query law-stone contract: %w", err)
	}
	if code == nil {
		return "", nil
	}

	decodedCode, err := base64.StdEncoding.DecodeString(*code)
	if err != nil {
		return "", fmt.Errorf("failed to decode law-stone code: %w", err)
	}

	return string(decodedCode), nil
}

func (c *queryClient) AskGovPermittedActions(ctx context.Context, addr, did string) ([]string, error) {
	gov, err := c.lawStoneFactory(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create law-stone client: %w", err)
	}

	response, err := gov.Ask(ctx, &lsschema.QueryMsg_Ask{Query: fmt.Sprintf("tell_permitted_actions('%s',Actions).", did)})
	if err != nil {
		return nil, fmt.Errorf("failed to query law-stone contract: %w", err)
	}

	if len(response.Answer.Results) != 1 {
		return nil, nil
	}
	if len(response.Answer.Results[0].Substitutions) != 1 {
		return nil, nil
	}

	result := response.Answer.Results[0].Substitutions[0].Expression
	result = result[1 : len(result)-1]
	actions := make([]string, 0)
	for _, action := range strings.Split(result, ",") {
		actions = append(actions, strings.Trim(action, "'"))
	}

	return actions, nil
}

func (c *queryClient) AskGovTellAction(ctx context.Context, addr, did, action string) (bool, error) {
	gov, err := c.lawStoneFactory(addr)
	if err != nil {
		return false, fmt.Errorf("failed to create law-stone client: %w", err)
	}

	response, err := gov.Ask(ctx, &lsschema.QueryMsg_Ask{Query: fmt.Sprintf("tell('%s','%s',Result,_).", did, action)})
	if err != nil {
		return false, fmt.Errorf("failed to query law-stone contract: %w", err)
	}

	if len(response.Answer.Results) != 1 {
		return false, nil
	}
	if len(response.Answer.Results[0].Substitutions) != 1 {
		return false, nil
	}

	return response.Answer.Results[0].Substitutions[0].Expression == "permitted", nil
}
