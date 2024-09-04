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

//go:embed vc-publication-tpl.jsonld
var publicationTemplate string

var _ credential.Descriptor = NewPublication()

// PublicationDescriptor is a descriptor for generate a digital resource publication VC.
// See https://docs.axone.xyz/ontology/next/schemas/credential-digital-resource-publication
type PublicationDescriptor struct {
	id           string
	datasetDID   string
	datasetURI   string
	storageDID   string
	issuanceDate *time.Time
}

func NewPublication() *PublicationDescriptor {
	t := time.Now().UTC()
	return &PublicationDescriptor{
		id:           uuid.New().String(),
		issuanceDate: &t,
	}
}

func (d *PublicationDescriptor) WithID(id string) *PublicationDescriptor {
	d.id = id
	return d
}

func (d *PublicationDescriptor) WithDatasetDID(did string) *PublicationDescriptor {
	d.datasetDID = did
	return d
}

func (d *PublicationDescriptor) WithDatasetURI(uri string) *PublicationDescriptor {
	d.datasetURI = uri
	return d
}

func (d *PublicationDescriptor) WithStorageDID(did string) *PublicationDescriptor {
	d.storageDID = did
	return d
}

func (d *PublicationDescriptor) WithIssuanceDate(t time.Time) *PublicationDescriptor {
	d.issuanceDate = &t
	return d
}

func (d *PublicationDescriptor) validate() error {
	if d.datasetDID == "" {
		return errors.New("dataset DID is required")
	}
	if d.datasetURI == "" {
		return errors.New("dataset URI is required")
	}
	if d.storageDID == "" {
		return errors.New("storage DID is required")
	}
	return nil
}

func (d *PublicationDescriptor) IssuedAt() *time.Time {
	return d.issuanceDate
}

func (d *PublicationDescriptor) ProofPurpose() string {
	return "assertionMethod"
}

func (d *PublicationDescriptor) Generate() (*bytes.Buffer, error) {
	err := d.validate()
	if err != nil {
		return nil, err
	}

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
