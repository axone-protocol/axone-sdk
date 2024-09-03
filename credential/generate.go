package credential

import (
	"bytes"
	"time"

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
//	vc, err := credential.New(
//	    NewGovernanceVC().
//	        WithDatasetDID("did:key:...").
//	        WithGovAddr("axone1234..."),
//	).
//	WithParser(parser).
//	WithSignature(signer, "did:key:..."). // Signature is optional and Generate a not signed VC if not provided.
//	Generate()
func New(descriptor Descriptor) *Generator {
	return &Generator{
		vc: descriptor,
	}
}

func (generator *Generator) WithParser(parser *DefaultParser) *Generator {
	generator.parser = parser
	return generator
}

func (generator *Generator) WithSignature(signer verifiable.Signer, did string) *Generator {
	generator.signer = signer
	generator.signerDID = did
	return generator
}

func (generator *Generator) Generate() (*verifiable.Credential, error) {
	raw, err := generator.vc.Generate()
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
