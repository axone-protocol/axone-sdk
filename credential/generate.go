package credential

import (
	"bytes"
	_ "embed"
	"errors"
	"html/template"
	"time"

	"github.com/axone-protocol/axone-sdk/dataverse"
	"github.com/google/uuid"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
)

//go:embed template/vc-gov-tpl.jsonld
var governanceTemplate string

// Generator is a verifiable credential generator.
type Generator struct {
	vc        Descriptor
	signer    verifiable.Signer
	signerDID string
	parser    *credentialParser
}

// New allow to generate a verifiable credential with the given credential descriptor.
// Example:
//
//	vc, err := credential.New(
//	    NewGovernanceVC().
//	        WithDatasetDID("did:key:...").
//	        WithGovAddr("axone1234..."),
//	).
//	WithParser(parser).
//	WithSignature(signer, "did:key:..."). // Signature is optional and generate a not signed VC if not provided.
//	Generate()
func New(descriptor Descriptor) *Generator {
	return &Generator{
		vc: descriptor,
	}
}

func (generator *Generator) WithParser(parser *credentialParser) *Generator {
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
		return nil, NewVCError(ErrGenerate, err)
	}

	if generator.parser == nil {
		return nil, NewVCError(ErrNoParser, nil)
	}
	cred, err := generator.parser.parse(raw.Bytes())
	if err != nil {
		return nil, NewVCError(ErrParse, err)
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
			return nil, NewVCError(ErrSign, err)
		}
	}

	return cred, nil
}

// Descriptor is an interface representing the description of a verifiable credential.
type Descriptor interface {
	issuedAt() *time.Time
	generate() (*bytes.Buffer, error)
	proofPurpose() string
}

var _ Descriptor = NewGovernanceVC()

// GovernanceVCDescriptor is a descriptor for a governance verifiable credential.
type GovernanceVCDescriptor struct {
	id           string
	datasetDID   string
	govAddr      string
	issuanceDate *time.Time
}

// NewGovernanceVC creates a new governance verifiable credential descriptor.
// DatasetDID and GovAddr are required. If ID is not provided, it will be generated.
// If issuance date is not provided, it will be set to the current time at the generation.
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

func (g *GovernanceVCDescriptor) issuedAt() *time.Time {
	return g.issuanceDate
}

func (g *GovernanceVCDescriptor) proofPurpose() string {
	return "authentication"
}

func (g *GovernanceVCDescriptor) generate() (*bytes.Buffer, error) {
	err := g.prepare()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &buf, nil
}
