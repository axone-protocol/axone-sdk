package credential

import (
	"bytes"
	"time"

	"github.com/axone-protocol/axone-sdk/keys"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
)

// Generator is a verifiable credential generator.
type Generator struct {
	vc        Descriptor
	signer    verifiable.Signer
	signerDID string
	parser    *DefaultParser
}

// New allow to Generate a verifiable credential with the given credential descriptor.
// Example:
//
//		vc, err := credential.New(
//		    template.NewGovernance(
//					"datasetID",
//					"addr",
//					WithID[*GovernanceDescriptor]("id")
//	     ),
//		    WithParser(parser),
//		    WithSigner(signer)). // Signature is optional and Generate a not signed VC if not provided.
//		Generate()
func New(descriptor Descriptor, opts ...Option) *Generator {
	g := &Generator{
		vc: descriptor,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Option is a function that configures a Generator.
type Option func(*Generator)

func WithParser(parser *DefaultParser) Option {
	return func(g *Generator) {
		g.parser = parser
	}
}

func WithSigner(signer keys.Keyring) Option {
	return func(g *Generator) {
		g.signer = signer
		g.signerDID = signer.DIDKeyID()
	}
}

func (generator *Generator) Generate() (*verifiable.Credential, error) {
	raw, err := generator.vc.Generate()
	if err != nil {
		return nil, NewVCError(ErrGenerate, err)
	}

	if generator.parser == nil {
		return nil, NewVCError(ErrNoParser, nil)
	}
	cred, err := generator.parser.Parse(raw.Bytes())
	if err != nil {
		return nil, NewVCError(ErrParse, err)
	}

	if generator.signer != nil {
		err := cred.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
			Created:                 generator.vc.IssuedAt(),
			SignatureType:           "EcdsaSecp256k1Signature2019",
			Suite:                   ecdsasecp256k1signature2019.New(suite.WithSigner(generator.signer)),
			SignatureRepresentation: verifiable.SignatureJWS,
			VerificationMethod:      generator.signerDID,
			Purpose:                 generator.vc.ProofPurpose(),
		}, jsonld.WithDocumentLoader(generator.parser.documentLoader))
		if err != nil {
			return nil, NewVCError(ErrSign, err)
		}
	}

	return cred, nil
}

// Descriptor is an interface representing the description of a verifiable credential.
type Descriptor interface {
	IssuedAt() *time.Time
	Generate() (*bytes.Buffer, error)
	ProofPurpose() string
}
