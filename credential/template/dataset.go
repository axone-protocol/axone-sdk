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

//go:embed vc-desc-tpl.jsonld
var datasetTemplate string

var _ credential.Descriptor = &DatasetDescriptor{}

// DatasetDescriptor is a descriptor for generate a dataset description VC.
// See https://docs.axone.xyz/ontology/next/schemas/credential-dataset-description
type DatasetDescriptor struct {
	id           string
	datasetDID   string
	title        string
	description  string
	format       string
	tags         []string
	topic        string
	issuanceDate *time.Time
}

func NewDataset(datasetDID, title string, opts ...Option[*DatasetDescriptor]) *DatasetDescriptor {
	t := time.Now().UTC()
	d := &DatasetDescriptor{
		id:           uuid.New().String(),
		datasetDID:   datasetDID,
		title:        title,
		issuanceDate: &t,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *DatasetDescriptor) setID(id string) {
	d.id = id
}

func (d *DatasetDescriptor) setDescription(description string) {
	d.description = description
}

func (d *DatasetDescriptor) setFormat(format string) {
	d.format = format
}

func (d *DatasetDescriptor) setTags(tags []string) {
	d.tags = tags
}

func (d *DatasetDescriptor) setTopic(topic string) {
	d.topic = topic
}

func (d *DatasetDescriptor) setIssuanceDate(t time.Time) {
	d.issuanceDate = &t
}

func (d *DatasetDescriptor) IssuedAt() *time.Time {
	return d.issuanceDate
}

func (d *DatasetDescriptor) ProofPurpose() string {
	return "assertionMethod"
}

func (d *DatasetDescriptor) Generate() (*bytes.Buffer, error) {
	tpl, err := gotemplate.New("datasetDescriptionVC").Parse(datasetTemplate)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	err = tpl.Execute(&buf, map[string]any{
		"NamespacePrefix": dataverse.W3IDPrefix,
		"CredID":          d.id,
		"DatasetDID":      d.datasetDID,
		"Title":           d.title,
		"Description":     d.description,
		"Format":          d.format,
		"Tags":            d.tags,
		"Topic":           d.topic,
		"IssuedAt":        d.issuanceDate.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
