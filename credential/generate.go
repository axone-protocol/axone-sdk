package credential

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"html/template"
	"time"
)

//go:embed template/vc-gov-tpl.jsonld
var governanceTemplate string

type Generator struct {
	vc        Descriptor
	signer    verifiable.Signer
	signerDID string
	parser    *CredentialParser
}

func NewGenerator(descriptor Descriptor) *Generator {
	return &Generator{
		vc: descriptor,
	}
}

func (generator *Generator) WithParser(parser *CredentialParser) *Generator {
	generator.parser = parser
	return generator
}

func (generator *Generator) WithSignature(signer verifiable.Signer, did string) *Generator {
	generator.signer = signer
	generator.signerDID = did
	return generator
}

func (generator *Generator) Generate() (*verifiable.Credential, error) {
	raw, err := generator.vc.generate()
	if err != nil {
		return nil, err // TODO: better err handler
	}

	if generator.parser == nil {
		return nil, errors.New("no parser provided")
	}
	cred, err := generator.parser.parse(raw.Bytes())
	if err != nil {
		return nil, err // TODO:
	}

	if generator.signer != nil {
		err := cred.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
			Created:                 generator.vc.issuedAt(),
			SignatureType:           "EcdsaSecp256k1Signature2019",
			Suite:                   ecdsasecp256k1signature2019.New(suite.WithSigner(generator.signer)),
			SignatureRepresentation: verifiable.SignatureJWS,
			VerificationMethod:      generator.signerDID,
			Purpose:                 generator.vc.proofPurpose(),
		}, jsonld.WithDocumentLoader(generator.parser.documentLoader))

		if err != nil {
			return nil, err // TODO:
		}
	}

	return cred, nil
}

type Descriptor interface {
	issuedAt() *time.Time
	generate() (*bytes.Buffer, error)
	proofPurpose() string
}

var _ Descriptor = NewGovernanceVC()

type GovernanceVCDescriptor struct {
	id           string
	datasetDID   string
	govAddr      string
	issuanceDate *time.Time
}

func (g *GovernanceVCDescriptor) issuedAt() *time.Time {
	return g.issuanceDate
}

func (g *GovernanceVCDescriptor) proofPurpose() string {
	return "authentication"
}

func NewGovernanceVC() *GovernanceVCDescriptor {
	return &GovernanceVCDescriptor{}
}

func (g *GovernanceVCDescriptor) WithID(id string) *GovernanceVCDescriptor {
	g.id = id
	return g
}

func (g *GovernanceVCDescriptor) WithDatasetDID(did string) *GovernanceVCDescriptor {
	g.datasetDID = did
	return g
}

func (g *GovernanceVCDescriptor) WithGovAddr(addr string) *GovernanceVCDescriptor {
	g.govAddr = addr
	return g
}

func (g *GovernanceVCDescriptor) WithIssuanceDate(t time.Time) *GovernanceVCDescriptor {
	g.issuanceDate = &t
	return g
}

func (g *GovernanceVCDescriptor) prepare() error {
	if g.id == "" {
		g.id = uuid.New().String()
	}
	if g.datasetDID == "" {
		return errors.New("dataset DID is required")
	}
	if g.govAddr == "" {
		return errors.New("governance address is required")
	}
	if g.issuanceDate == nil {
		t := time.Now().UTC()
		g.issuanceDate = &t
	}
	return nil
}

func (g *GovernanceVCDescriptor) generate() (*bytes.Buffer, error) {
	err := g.prepare()
	if err != nil {
		return nil, err // todo
	}

	tpl, err := template.New("governanceVC").Parse(governanceTemplate)
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
		return nil, err // todo
	}

	return &buf, nil
}
