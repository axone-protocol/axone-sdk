package template

import (
	"testing"
	"time"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGovernanceDescriptor_Generate(t *testing.T) {
	tests := []struct {
		name    string
		vc      credential.Descriptor
		wantErr error
		check   func(*verifiable.Credential)
	}{
		{
			name: "Valid governance VC",
			vc: NewGovernance(
				"datasetID",
				"addr",
				WithID[*GovernanceDescriptor]("id"),
				WithIssuanceDate[*GovernanceDescriptor](time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldEqual, "https://w3id.org/axone/ontology/v4/schema/credential/governance/text/id")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["isGovernedBy"].(map[string]interface{})["fromGovernance"], ShouldEqual, "addr")
				So(vc.Issued.Time, ShouldEqual, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		{
			name: "Valid governance VC with default value",
			vc: NewGovernance(
				"datasetID",
				"addr",
			),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldStartWith, "https://w3id.org/axone/ontology/v4/schema/credential/governance/text/")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["isGovernedBy"].(map[string]interface{})["fromGovernance"], ShouldEqual, "addr")
				So(vc.Issued.Time, ShouldHappenWithin, time.Second, time.Now().UTC())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a credential generator", t, func() {
				docLoader, err := testutil.MockDocumentLoader()
				So(err, ShouldBeNil)

				parser := credential.NewDefaultParser(docLoader)
				generator := credential.New(
					test.vc,
					credential.WithParser(parser))

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
