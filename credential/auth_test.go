//nolint:lll
package credential_test

import (
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/axone-protocol/axone-sdk/credential"
	jld "github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ldcontext"
	mockldstore "github.com/hyperledger/aries-framework-go/pkg/mock/ld"
	mockprovider "github.com/hyperledger/aries-framework-go/pkg/mock/provider"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	//go:embed testdata/contexts/credentials-v1.jsonld
	mockCredentialsV1JSONLD []byte
	//go:embed testdata/contexts/authentication-v4.jsonld
	mockAuthenticationV4JSONLD []byte
	//go:embed testdata/contexts/security-v2.jsonld
	mockSecurityV2JSONLD []byte
)

func TestAuthParser_ParseSigned(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr error
		result  *credential.AuthClaim
	}{
		{
			name:    "valid credential",
			file:    "testdata/valid.jsonld",
			wantErr: nil,
			result: &credential.AuthClaim{
				ID:        "did:key:zQ3shhCAzQcroi4RqZ48eNudKWf75Fvv9ryJsxbaWCCPsfnFj",
				ToService: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			},
		},
		{
			name:    "credential not signed",
			file:    "testdata/invalid_not-signed.jsonld",
			wantErr: credential.NewVCError(credential.ErrInvalidProof, credential.NewVCError(credential.ErrMissingProof, nil)),
			result:  nil,
		},
		{
			name:    "credential with invalid signature",
			file:    "testdata/invalid_wrong-signature.jsonld",
			wantErr: credential.NewVCError(credential.ErrParse, fmt.Errorf("decode new credential: check embedded proof: check linked data proof: ecdsa: invalid signature")),
			result:  nil,
		},
		{
			name:    "credential without subject",
			file:    "testdata/invalid_malformated-subject.jsonld",
			wantErr: credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrMalformedSubject, nil)),
			result:  nil,
		},
		{
			name:    "credential with multiple subject",
			file:    "testdata/invalid_multiple-subject.jsonld",
			wantErr: credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrExpectSingleClaim, nil)),
			result:  nil,
		},
		{
			name:    "credential with issuer different from subject",
			file:    "testdata/invalid_issuer-differs-subject.jsonld",
			wantErr: credential.NewVCError(credential.ErrAuthClaim, fmt.Errorf("subject differs from issuer (subject: `did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr`, issuer: `did:key:zQ3shhCAzQcroi4RqZ48eNudKWf75Fvv9ryJsxbaWCCPsfnFj`)")),
			result:  nil,
		},
		{
			name:    "toService claims is not a string",
			file:    "testdata/invalid_service-not-string.jsonld",
			wantErr: credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrExtractClaim, fmt.Errorf("key 'toService' is not a string"))),
			result:  nil,
		},
		{
			name:    "toService claims key missing",
			file:    "testdata/invalid_service-key-missing.jsonld",
			wantErr: credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrExtractClaim, fmt.Errorf("key 'toService' not found"))),
			result:  nil,
		},
		{
			name:    "credential expired",
			file:    "testdata/invalid_expired.jsonld",
			wantErr: credential.NewVCError(credential.ErrExpired, fmt.Errorf("2023-01-01 00:00:00 +0000 UTC")),
			result:  nil,
		},
		{
			name:    "credential not issued now",
			file:    "testdata/invalid_futur-issued.jsonld",
			wantErr: credential.NewVCError(credential.ErrIssued, fmt.Errorf("2200-01-01 20:30:59.627706 +0200 +0200")),
			result:  nil,
		},
		{
			name:    "credential not issued now",
			file:    "testdata/invalid_futur-issued.jsonld",
			wantErr: credential.NewVCError(credential.ErrIssued, fmt.Errorf("2200-01-01 20:30:59.627706 +0200 +0200")),
			result:  nil,
		},
		{
			name:    "credential with not authentication proof purpose",
			file:    "testdata/invalid_not-authentication-proof.jsonld",
			wantErr: credential.NewVCError(credential.ErrAuthClaim, fmt.Errorf("proof purpose not targeting `authentication` (proof purpose: `assertionMethod`)")),
			result:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a credential", t, func() {
				raw, err := os.ReadFile(test.file)
				So(err, ShouldBeNil)

				docLoader, err := jld.NewDocumentLoader(createMockCtxProvider(), jld.WithExtraContexts(
					ldcontext.Document{
						URL:     "https://w3id.org/axone/ontology/v4/schema/credential/digital-service/authentication/",
						Content: mockAuthenticationV4JSONLD,
					},
					ldcontext.Document{
						URL:     "https://www.w3.org/2018/credentials/v1",
						Content: mockCredentialsV1JSONLD,
					},
					ldcontext.Document{
						URL:     "https://w3id.org/security/v2",
						Content: mockSecurityV2JSONLD,
					},
				))
				So(err, ShouldBeNil)

				parser := credential.NewAuthParser(docLoader)

				Convey("When parsing the credential", func() {
					authClaim, err := parser.ParseSigned(raw)

					Convey("Then the result should be as expected", func() {
						if test.wantErr != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldResemble, test.wantErr.Error())
						} else {
							So(err, ShouldBeNil)
						}
						So(authClaim, ShouldResemble, test.result)
					})
				})
			})
		})
	}
}

func createMockCtxProvider() *mockprovider.Provider {
	p := &mockprovider.Provider{
		ContextStoreValue:        mockldstore.NewMockContextStore(),
		RemoteProviderStoreValue: mockldstore.NewMockRemoteProviderStore(),
	}

	return p
}
