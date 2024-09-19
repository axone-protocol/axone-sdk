package template

import (
	"bytes"
	_ "embed"
	gotemplate "text/template"
	"time"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/google/uuid"
)

//go:embed vc-gov-tpl.jsonld
var governanceTemplate string

var _ credential.Descriptor = &GovernanceDescriptor{}

// GovernanceDescriptor is a descriptor for generate a credential governance text VC.
// See https://docs.axone.xyz/ontology/next/schemas/credential-governance-text
type GovernanceDescriptor struct {
	id           string
	datasetDID   string
	govAddr      string
	issuanceDate *time.Time
}

// NewGovernance creates a new governance verifiable credential descriptor.
// DatasetDID and GovAddr are required. If ID is not provided, it will be generated.
// If issuance date is not provided, it will be set to the current time at descriptor instantiation.
func NewGovernance(datasetDID, govAddr string, opts ...Option[*GovernanceDescriptor]) *GovernanceDescriptor {
	t := time.Now().UTC()
	g := &GovernanceDescriptor{
		id:           uuid.New().String(),
		datasetDID:   datasetDID,
		govAddr:      govAddr,
		issuanceDate: &t,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *GovernanceDescriptor) setID(id string) {
	g.id = id
}

func (g *GovernanceDescriptor) setIssuanceDate(t time.Time) {
	g.issuanceDate = &t
}

func (g *GovernanceDescriptor) IssuedAt() *time.Time {
	return g.issuanceDate
}

func (g *GovernanceDescriptor) ProofPurpose() string {
	return credential.ProofPurposeAssertionMethod
}

func (g *GovernanceDescriptor) Generate() (*bytes.Buffer, error) {
	tpl, err := gotemplate.New("governanceVC").Parse(governanceTemplate)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	err = tpl.Execute(&buf, map[string]string{
		"NamespacePrefix": dataverse.W3IDPrefix,
		"CredID":          g.id,
		"DatasetDID":      g.datasetDID,
		"GovAddr":         g.govAddr,
		"IssuedAt":        g.issuanceDate.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
