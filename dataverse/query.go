package dataverse

import (
	"fmt"
	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5"
)

const W3IDPrefix = "https://w3id.org/axone/ontology/v4"

func buildGetResourceGovAddrRequest(resource string) cgschema.SelectQuery {
	limit := 1
	selectVar := cgschema.SelectItem_Variable("code")
	codeVarOrNodeOrLit := cgschema.VarOrNodeOrLiteral_Variable("code")
	credId := cgschema.VarOrNode_Variable("credId")
	resourceIRI := cgschema.IRI_Full(resource)
	govType := cgschema.IRI_Prefixed("gov:GovernanceTextCredential")
	claimVar := cgschema.VarOrNode_Variable("claim")
	claimVarOrNode := cgschema.VarOrNodeOrLiteral_Variable("claim")
	isGovernedBy := cgschema.IRI_Prefixed("gov:isGovernedBy")
	fromGovernance := cgschema.IRI_Prefixed("gov:fromGovernance")
	govVarOrNodeOrLit := cgschema.VarOrNodeOrLiteral_Variable("gov")
	govVar := cgschema.VarOrNode_Variable("gov")

	return cgschema.SelectQuery{
		Limit: &limit,
		Prefixes: []cgschema.Prefix{
			{
				Prefix:    "gov",
				Namespace: fmt.Sprintf("%s/schema/credential/governance/text/", W3IDPrefix),
			},
		},
		Select: []cgschema.SelectItem{
			{
				Variable: &selectVar,
			},
		},
		Where: cgschema.WhereClause{
			Bgp: &cgschema.WhereClause_Bgp{
				Patterns: []cgschema.TriplePattern{
					{
						Subject:   cgschema.VarOrNode{Variable: &credId},
						Predicate: cgschema.VarOrNamedNode{NamedNode: &cgschema.VarOrNamedNode_NamedNode{Full: &VcBodySubject}},
						Object:    cgschema.VarOrNodeOrLiteral{Node: &cgschema.VarOrNodeOrLiteral_Node{NamedNode: &cgschema.Node_NamedNode{Full: &resourceIRI}}},
					},
					{
						Subject:   cgschema.VarOrNode{Variable: &credId},
						Predicate: cgschema.VarOrNamedNode{NamedNode: &cgschema.VarOrNamedNode_NamedNode{Full: &VcBodyType}},
						Object:    cgschema.VarOrNodeOrLiteral{Node: &cgschema.VarOrNodeOrLiteral_Node{NamedNode: &cgschema.Node_NamedNode{Prefixed: &govType}}},
					},
					{
						Subject:   cgschema.VarOrNode{Variable: &credId},
						Predicate: cgschema.VarOrNamedNode{NamedNode: &cgschema.VarOrNamedNode_NamedNode{Full: &VcBodyClaim}},
						Object:    cgschema.VarOrNodeOrLiteral{Variable: &claimVarOrNode},
					},
					{
						Subject:   cgschema.VarOrNode{Variable: &claimVar},
						Predicate: cgschema.VarOrNamedNode{NamedNode: &cgschema.VarOrNamedNode_NamedNode{Prefixed: &isGovernedBy}},
						Object:    cgschema.VarOrNodeOrLiteral{Variable: &govVarOrNodeOrLit},
					},
					{
						Subject:   cgschema.VarOrNode{Variable: &govVar},
						Predicate: cgschema.VarOrNamedNode{NamedNode: &cgschema.VarOrNamedNode_NamedNode{Prefixed: &fromGovernance}},
						Object:    cgschema.VarOrNodeOrLiteral{Variable: &codeVarOrNodeOrLit},
					},
				},
			},
		},
	}
}
