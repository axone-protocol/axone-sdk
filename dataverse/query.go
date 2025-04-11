package dataverse

import (
	"fmt"

	cgschema "github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6"
)

const W3IDPrefix = "https://w3id.org/axone/ontology/v4"

func ref[T any](v T) *T {
	return &v
}

//nolint:funlen
func buildGetResourceGovAddrRequest(resource string) cgschema.SelectQuery {
	return cgschema.SelectQuery{
		Limit: ref(1),
		Prefixes: []cgschema.Prefix{
			{
				Prefix:    "gov",
				Namespace: fmt.Sprintf("%s/schema/credential/governance/text/", W3IDPrefix),
			},
		},
		Select: []cgschema.SelectItem{
			{
				Variable: ref(cgschema.SelectItem_Variable("code")),
			},
		},
		Where: cgschema.WhereClause{
			Bgp: &cgschema.WhereClause_Bgp{
				Patterns: []cgschema.TriplePattern{
					{
						Subject: cgschema.VarOrNode{Variable: ref(cgschema.VarOrNode_Variable("credId"))},
						Predicate: cgschema.VarOrNamedNode{
							NamedNode: &cgschema.VarOrNamedNode_NamedNode{Full: &VcBodySubject},
						},
						Object: cgschema.VarOrNodeOrLiteral{
							Node: &cgschema.VarOrNodeOrLiteral_Node{
								NamedNode: &cgschema.Node_NamedNode{Full: ref(cgschema.IRI_Full(resource))},
							},
						},
					},
					{
						Subject: cgschema.VarOrNode{Variable: ref(cgschema.VarOrNode_Variable("credId"))},
						Predicate: cgschema.VarOrNamedNode{
							NamedNode: &cgschema.VarOrNamedNode_NamedNode{Full: &VcBodyType},
						},
						Object: cgschema.VarOrNodeOrLiteral{
							Node: &cgschema.VarOrNodeOrLiteral_Node{
								NamedNode: &cgschema.Node_NamedNode{Prefixed: ref(cgschema.IRI_Prefixed("gov:GovernanceTextCredential"))},
							},
						},
					},
					{
						Subject: cgschema.VarOrNode{Variable: ref(cgschema.VarOrNode_Variable("credId"))},
						Predicate: cgschema.VarOrNamedNode{
							NamedNode: &cgschema.VarOrNamedNode_NamedNode{Full: &VcBodyClaim},
						},
						Object: cgschema.VarOrNodeOrLiteral{Variable: ref(cgschema.VarOrNodeOrLiteral_Variable("claim"))},
					},
					{
						Subject: cgschema.VarOrNode{Variable: ref(cgschema.VarOrNode_Variable("claim"))},
						Predicate: cgschema.VarOrNamedNode{
							NamedNode: &cgschema.VarOrNamedNode_NamedNode{Prefixed: ref(cgschema.IRI_Prefixed("gov:isGovernedBy"))},
						},
						Object: cgschema.VarOrNodeOrLiteral{Variable: ref(cgschema.VarOrNodeOrLiteral_Variable("gov"))},
					},
					{
						Subject: cgschema.VarOrNode{Variable: ref(cgschema.VarOrNode_Variable("gov"))},
						Predicate: cgschema.VarOrNamedNode{
							NamedNode: &cgschema.VarOrNamedNode_NamedNode{Prefixed: ref(cgschema.IRI_Prefixed("gov:fromGovernance"))},
						},
						Object: cgschema.VarOrNodeOrLiteral{Variable: ref(cgschema.VarOrNodeOrLiteral_Variable("code"))},
					},
				},
			},
		},
	}
}
