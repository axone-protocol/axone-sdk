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

//go:embed vc-desc-tpl.jsonld
var datasetTemplate string

var _ credential.Descriptor = NewDataset()

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

func NewDataset() *DatasetDescriptor {
	t := time.Now().UTC()
	return &DatasetDescriptor{
		id:           uuid.New().String(),
		issuanceDate: &t,
	}
}

func (d *DatasetDescriptor) WithID(id string) *DatasetDescriptor {
	d.id = id
	return d
}

func (d *DatasetDescriptor) WithDatasetDID(did string) *DatasetDescriptor {
	d.datasetDID = did
	return d
}

func (d *DatasetDescriptor) WithTitle(title string) *DatasetDescriptor {
	d.title = title
	return d
}

func (d *DatasetDescriptor) WithDescription(description string) *DatasetDescriptor {
	d.description = description
	return d
}

func (d *DatasetDescriptor) WithFormat(format string) *DatasetDescriptor {
	d.format = format
	return d
}

func (d *DatasetDescriptor) WithTags(tags []string) *DatasetDescriptor {
	d.tags = tags
	return d
}

func (d *DatasetDescriptor) WithTopic(topic string) *DatasetDescriptor {
	d.topic = topic
	return d
}

func (d *DatasetDescriptor) WithIssuanceDate(t time.Time) *DatasetDescriptor {
	d.issuanceDate = &t
	return d
}

func (d *DatasetDescriptor) validate() error {
	if d.datasetDID == "" {
		return errors.New("dataset DID is required")
	}
	if d.title == "" {
		return errors.New("title is required")
	}
	return nil
}

func (d *DatasetDescriptor) IssuedAt() *time.Time {
	return d.issuanceDate
}

func (d *DatasetDescriptor) ProofPurpose() string {
	return "assertionMethod"
}

func (d *DatasetDescriptor) Generate() (*bytes.Buffer, error) {
	err := d.validate()
	if err != nil {
		return nil, err
	}

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
