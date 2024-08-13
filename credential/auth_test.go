package credential_test

import (
	_ "embed"
	"fmt"
	"github.com/axone-protocol/axone-sdk/credential"
	jld "github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ldcontext"
	mockldstore "github.com/hyperledger/aries-framework-go/pkg/mock/ld"
	mockprovider "github.com/hyperledger/aries-framework-go/pkg/mock/provider"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
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
		serviceId string
		file      string
		wantErr   error
		result    *credential.AuthClaim
	}{
		{
			name:      "valid credential",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/valid.jsonld",
			wantErr:   nil,
			result: &credential.AuthClaim{
				ID:        "did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr",
				ToService: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			},
		},
		{
			name:      "credential not signed",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_not-signed.jsonld",
			wantErr:   fmt.Errorf("missing verifiable credential proof"),
			result:    nil,
		},
		{
			name:      "credential not signed",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_wrong-signature.jsonld",
			wantErr:   fmt.Errorf("decode new credential: check embedded proof: check linked data proof: ecdsa: invalid signature"),
			result:    nil,
		},
		{
			name:      "valid credential but wrong service id",
			serviceId: "did:key:foo",
			file:      "testdata/valid.jsonld",
			wantErr:   fmt.Errorf("auth claim target doesn't match current service id: did:key:foo (target: did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz)"),
			result:    nil,
		},
		{
			name:      "credential without subject",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_malformated-subject.jsonld",
			wantErr:   fmt.Errorf("malformed auth claim: malformed vc subject"),
			result:    nil,
		},
		{
			name:      "credential with multiple subject",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_multiple-subject.jsonld",
			wantErr:   fmt.Errorf("malformed auth claim: expected a single vc claim"),
			result:    nil,
		},
		{
			name:      "credential with issuer different from subject",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_issuer-differs-subject.jsonld",
			wantErr:   fmt.Errorf("auth claim subject differs from issuer (subject: did:key:zQ3shpoUHzwcgdt2gxjqHHnJnNkBVd4uX3ZBhmPiM7J93yqCr, issuer: did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz)"),
			result:    nil,
		},
		{
			name:      "toService claims is not a string",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_service-not-string.jsonld",
			wantErr:   fmt.Errorf("malformed auth claim: key 'toService' is not a string"),
			result:    nil,
		},
		{
			name:      "toService claims key missing",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_service-key-missing.jsonld",
			wantErr:   fmt.Errorf("malformed auth claim: key 'toService' not found"),
			result:    nil,
		},
		{
			name:      "credential expired",
			serviceId: "did:key:zQ3shZxyDoD3QorxHJrFS68EjzDgQZSqZcj3wQqc1ngbF1vgz",
			file:      "testdata/invalid_expired.jsonld",
			wantErr:   fmt.Errorf("verifiable credential expired"),
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

				parser := credential.NewAuthParser(test.serviceId, docLoader)

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