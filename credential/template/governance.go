package template

import (
	"bytes"
	_ "embed"
	"errors"
	gotemplate "text/template"
	"time"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/google/uuid"
)

//go:embed vc-gov-tpl.jsonld
var governanceTemplate string

var _ credential.Descriptor = NewGovernance()

// GovernanceDescriptor is a descriptor for a governance verifiable credential.
type GovernanceDescriptor struct {
	id           string
	datasetDID   string
	govAddr      string
	issuanceDate *time.Time
}

// NewGovernance creates a new governance verifiable credential descriptor.
// DatasetDID and GovAddr are required. If ID is not provided, it will be generated.
// If issuance date is not provided, it will be set to the current time at descriptor instantiation.
func NewGovernance() *GovernanceDescriptor {
	t := time.Now().UTC()
	return &GovernanceDescriptor{
		id:           uuid.New().String(),
		issuanceDate: &t,
	}
}

func (g *GovernanceDescriptor) WithID(id string) *GovernanceDescriptor {
	g.id = id
	return g
}

func (g *GovernanceDescriptor) WithDatasetDID(did string) *GovernanceDescriptor {
	g.datasetDID = did
	return g
}

func (g *GovernanceDescriptor) WithGovAddr(addr string) *GovernanceDescriptor {
	g.govAddr = addr
	return g
}

func (g *GovernanceDescriptor) WithIssuanceDate(t time.Time) *GovernanceDescriptor {
	g.issuanceDate = &t
	return g
}

func (g *GovernanceDescriptor) validate() error {
	if g.datasetDID == "" {
		return errors.New("dataset DID is required")
	}
	if g.govAddr == "" {
		return errors.New("governance address is required")
	}
	return nil
}

func (g *GovernanceDescriptor) IssuedAt() *time.Time {
	return g.issuanceDate
}

func (g *GovernanceDescriptor) ProofPurpose() string {
	return "authentication"
}

func (g *GovernanceDescriptor) Generate() (*bytes.Buffer, error) {
	err := g.validate()
	if err != nil {
		return nil, err
	}

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
