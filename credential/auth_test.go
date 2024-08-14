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
		name      string
		serviceID string
		file      string
		wantErr   error
		result    *credential.AuthClaim
	}{
		{
			name:      "valid credential",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/valid.jsonld",
			wantErr:   nil,
			result: &credential.AuthClaim{
				ID:        "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
				ToService: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			},
		},
		{
			name:      "credential not signed",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_not-signed.jsonld",
			wantErr:   credential.NewVCError(credential.ErrMissingProof, nil),
			result:    nil,
		},
		{
			name:      "credential with invalid signature",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_wrong-signature.jsonld",
			wantErr:   fmt.Errorf("decode new credential: check embedded proof: check linked data proof: ecdsa: invalid signature"),
			result:    nil,
		},
		{
			name:      "valid credential but wrong service id",
			serviceID: "did:key:foo",
			file:      "testdata/valid.jsonld",
			wantErr:   credential.NewVCError(credential.ErrAuthClaim, fmt.Errorf("target doesn't match current service id: did:key:foo (target: did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz)")),
			result:    nil,
		},
		{
			name:      "credential without subject",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_malformated-subject.jsonld",
			wantErr:   credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrMalformedSubject, nil)),
			result:    nil,
		},
		{
			name:      "credential with multiple subject",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_multiple-subject.jsonld",
			wantErr:   credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrExpectSingleClaim, nil)),
			result:    nil,
		},
		{
			name:      "credential with issuer different from subject",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_issuer-differs-subject.jsonld",
			wantErr:   credential.NewVCError(credential.ErrAuthClaim, fmt.Errorf("subject differs from issuer (subject: did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr, issuer: did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz)")),
			result:    nil,
		},
		{
			name:      "toService claims is not a string",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_service-not-string.jsonld",
			wantErr:   credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrExtractClaim, fmt.Errorf("key 'toService' is not a string"))),
			result:    nil,
		},
		{
			name:      "toService claims key missing",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_service-key-missing.jsonld",
			wantErr:   credential.NewVCError(credential.ErrMalformed, credential.NewVCError(credential.ErrExtractClaim, fmt.Errorf("key 'toService' not found"))),
			result:    nil,
		},
		{
			name:      "credential expired",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_expired.jsonld",
			wantErr:   credential.NewVCError(credential.ErrExpired, fmt.Errorf("2023-01-01 00:00:00 +0000 UTC")),
			result:    nil,
		},
		{
			name:      "credential not issued now",
			serviceID: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_futur-issued.jsonld",
			wantErr:   credential.NewVCError(credential.ErrIssued, fmt.Errorf("2200-01-01 20:30:59.627706 +0200 +0200")),
			result:    nil,
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

				parser := credential.NewAuthParser(test.serviceID, docLoader)

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
