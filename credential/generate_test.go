package credential_test

import (
	"errors"
	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestGenerateGovernanceVC(t *testing.T) {
	tests := []struct {
		name    string
		vc      credential.Descriptor
		wantErr error
		check   func(*verifiable.Credential)
	}{
		{
			name: "Valid governance VC",
			vc: credential.NewGovernanceVC().
				WithID("id").
				WithDatasetDID("datasetID").
				WithGovAddr("addr").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldEqual, "https://w3id.org/axone/ontology/v4/schema/credential/governance/text/id")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["isGovernedBy"].(map[string]interface{})["fromGovernance"], ShouldEqual, "addr")
				So(vc.Issued.Time, ShouldEqual, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		{
			name: "Valid governance VC without id",
			vc: credential.NewGovernanceVC().
				WithDatasetDID("datasetID").
				WithGovAddr("addr").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldStartWith, "https://w3id.org/axone/ontology/v4/schema/credential/governance/text/")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["isGovernedBy"].(map[string]interface{})["fromGovernance"], ShouldEqual, "addr")
				So(vc.Issued.Time, ShouldEqual, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		{
			name: "Valid governance VC without issuance date",
			vc: credential.NewGovernanceVC().
				WithID("id").
				WithDatasetDID("datasetID").
				WithGovAddr("addr"),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldEqual, "https://w3id.org/axone/ontology/v4/schema/credential/governance/text/id")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["isGovernedBy"].(map[string]interface{})["fromGovernance"], ShouldEqual, "addr")
				So(vc.Issued.Time, ShouldHappenWithin, time.Second, time.Now().UTC())
			},
		},
		{
			name: "Missing dataset id",
			vc: credential.NewGovernanceVC().
				WithID("id").
				WithGovAddr("addr").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: errors.New("dataset DID is required"),
		},
		{
			name: "Missing gov addr",
			vc: credential.NewGovernanceVC().
				WithID("id").
				WithDatasetDID("datasetID").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: errors.New("governance address is required"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a credential generator", t, func() {
				docLoader, err := mockDocumentLoader()
				So(err, ShouldBeNil)

				parser := credential.NewCredentialParser(docLoader)
				generator := credential.NewGenerator(test.vc).
					WithParser(parser)

				Convey("When a governance VC is generated", func() {
					vc, err := generator.Generate()

					Convey("Then the governance VC should be generated", func() {
						if test.wantErr != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(vc, ShouldBeNil)
						} else {
							So(vc, ShouldNotBeNil)
							test.check(vc)
							So(err, ShouldBeNil)
						}
					})
				})
			})
		})
	}
}

func vcSubject(vc *verifiable.Credential) verifiable.Subject {
	subjects, ok := vc.Subject.([]verifiable.Subject)
	if !ok {
		panic("invalid subject type")
	}

	return subjects[0]
}

func TestGenerator_Generate(t *testing.T) {
	loader, _ := mockDocumentLoader()
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
			generator: credential.NewGenerator(credential.NewGovernanceVC().WithDatasetDID("datasetID").WithGovAddr("addr")),
			wantErr:   errors.New("no parser provided"),
		},
		{
			name:      "with descriptor error",
			generator: credential.NewGenerator(credential.NewGovernanceVC().WithDatasetDID("datasetID")),
			wantErr:   errors.New("governance address is required"),
		},
		{
			name: "without signature",
			generator: credential.NewGenerator(credential.NewGovernanceVC().WithDatasetDID("datasetID").WithGovAddr("addr")).
				WithParser(credential.NewCredentialParser(loader)),
			wantErr: nil,
			check: func(vc *verifiable.Credential) {
				So(len(vc.Proofs), ShouldEqual, 0)
			},
		},
		{
			name: "with signature",
			generator: credential.NewGenerator(credential.NewGovernanceVC().WithDatasetDID("datasetID").WithGovAddr("addr")).
				WithParser(credential.NewCredentialParser(loader)).
				WithSignature(mockSigner, "did:example:123"),
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
