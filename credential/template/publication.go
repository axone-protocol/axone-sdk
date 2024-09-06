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

//go:embed vc-publication-tpl.jsonld
var publicationTemplate string

var _ credential.Descriptor = &PublicationDescriptor{}

// PublicationDescriptor is a descriptor for generate a digital resource publication VC.
// See https://docs.axone.xyz/ontology/next/schemas/credential-digital-resource-publication
type PublicationDescriptor struct {
	id           string
	datasetDID   string
	datasetURI   string
	storageDID   string
	issuanceDate *time.Time
}

func NewPublication(datasetDID, datasetURI, storageDID string,
	opts ...Option[*PublicationDescriptor],
) *PublicationDescriptor {
	t := time.Now().UTC()
	p := &PublicationDescriptor{
		id:           uuid.New().String(),
		datasetDID:   datasetDID,
		datasetURI:   datasetURI,
		storageDID:   storageDID,
		issuanceDate: &t,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (d *PublicationDescriptor) setID(id string) {
	d.id = id
}

func (d *PublicationDescriptor) setIssuanceDate(t time.Time) {
	d.issuanceDate = &t
}

func (d *PublicationDescriptor) IssuedAt() *time.Time {
	return d.issuanceDate
}

func (d *PublicationDescriptor) ProofPurpose() string {
	return "assertionMethod"
}

func (d *PublicationDescriptor) Generate() (*bytes.Buffer, error) {
	tpl, err := gotemplate.New("publicationVC").Parse(publicationTemplate)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	err = tpl.Execute(&buf, map[string]any{
		"NamespacePrefix": dataverse.W3IDPrefix,
		"CredID":          d.id,
		"DatasetDID":      d.datasetDID,
		"DatasetURI":      d.datasetURI,
		"StorageDID":      d.storageDID,
		"IssuedAt":        d.issuanceDate.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
