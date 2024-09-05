package credential_test

import (
	"errors"
	"testing"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/credential/template"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestGenerator_Generate(t *testing.T) {
	loader, _ := testutil.MockDocumentLoader()
	controller := gomock.NewController(t)

	mockSigner := testutil.NewMockSigner(controller)
	mockSigner.EXPECT().Sign(gomock.Any()).Return([]byte("signature"), nil).AnyTimes()
	mockSigner.EXPECT().Alg().AnyTimes()

	tests := []struct {
		name      string
		generator *credential.Generator
		wantErr   error
		check     func(*verifiable.Credential)
	}{
		{
			name:      "without parser",
			generator: credential.New(template.NewGovernance().WithDatasetDID("datasetID").WithGovAddr("addr")),
			wantErr:   errors.New("no parser provided"),
		},
		{
			name:      "with descriptor error",
			generator: credential.New(template.NewGovernance().WithDatasetDID("datasetID")),
			wantErr:   credential.NewVCError(credential.ErrGenerate, errors.New("governance address is required")),
		},
		{
			name: "without signature",
			generator: credential.New(template.NewGovernance().WithDatasetDID("datasetID").WithGovAddr("addr"),
				credential.WithParser(credential.NewDefaultParser(loader))),
			wantErr: nil,
			check: func(vc *verifiable.Credential) {
				So(len(vc.Proofs), ShouldEqual, 0)
			},
		},
		{
			name: "with signature",
			generator: credential.New(template.NewGovernance().WithDatasetDID("datasetID").WithGovAddr("addr"),
				credential.WithParser(credential.NewDefaultParser(loader)), credential.WithSigner(mockSigner, "did:example:123")),
			wantErr: nil,
			check: func(vc *verifiable.Credential) {
				So(len(vc.Proofs), ShouldEqual, 1)
				So(vc.Proofs[0]["verificationMethod"], ShouldEqual, "did:example:123")
				So(vc.Proofs[0]["proofPurpose"], ShouldEqual, "authentication")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a credential generator", t, func() {
				Convey("When generating a credential", func() {
					vc, err := test.generator.Generate()

					Convey("Then an error or vc should be returned", func() {
						if test.wantErr == nil {
							So(err, ShouldBeNil)
							So(vc, ShouldNotBeNil)
							test.check(vc)
						} else {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
						}
					})
				})
			})
		})
	}
}
