package testutil

import (
	_ "embed"

	jld "github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ldcontext"
	mockldstore "github.com/hyperledger/aries-framework-go/pkg/mock/ld"
	mockprovider "github.com/hyperledger/aries-framework-go/pkg/mock/provider"
)

var (
	//go:embed contexts/credentials-v1.jsonld
	mockCredentialsV1JSONLD []byte
	//go:embed contexts/authentication-v4.jsonld
	mockAuthenticationV4JSONLD []byte
	//go:embed contexts/governance-text-v4.jsonld
	mockGovernanceTextV4JSONLD []byte
	//go:embed contexts/security-v2.jsonld
	mockSecurityV2JSONLD []byte
	//go:embed contexts/dataset-v4.jsonld
	mockDatasetV4JSONLD []byte
	//go:embed contexts/publication-v4.jsonld
	mockPublicationV4JSONLD []byte
)

func MockDocumentLoader() (*jld.DocumentLoader, error) {
	return jld.NewDocumentLoader(createMockCtxProvider(), jld.WithExtraContexts(
		ldcontext.Document{
			URL:     "https://w3id.org/axone/ontology/v4/schema/credential/digital-service/authentication/",
			Content: mockAuthenticationV4JSONLD,
		},
		ldcontext.Document{
			URL:     "https://w3id.org/axone/ontology/v4/schema/credential/governance/text/",
			Content: mockGovernanceTextV4JSONLD,
		},
		ldcontext.Document{
			URL:     "https://www.w3.org/2018/credentials/v1",
			Content: mockCredentialsV1JSONLD,
		},
		ldcontext.Document{
			URL:     "https://w3id.org/security/v2",
			Content: mockSecurityV2JSONLD,
		},
		ldcontext.Document{
			URL:     "https://w3id.org/axone/ontology/v4/schema/credential/dataset/description/",
			Content: mockDatasetV4JSONLD,
		},
		ldcontext.Document{
			URL:     "https://w3id.org/axone/ontology/v4/schema/credential/digital-resource/publication/",
			Content: mockPublicationV4JSONLD,
		},
	))
}

func createMockCtxProvider() *mockprovider.Provider {
	p := &mockprovider.Provider{
		ContextStoreValue:        mockldstore.NewMockContextStore(),
		RemoteProviderStoreValue: mockldstore.NewMockRemoteProviderStore(),
	}

	return p
}
