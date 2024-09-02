package credential_test

import (
	"errors"
	"github.com/axone-protocol/axone-sdk/credential"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGenerateGovernanceVC(t *testing.T) {
	tests := []struct {
		name    string
		vc      credential.Descriptor
		wantErr error
	}{
		{
			name: "Valid governance VC",
			vc: credential.NewGovernanceVC().
				WithID("id").
				WithDatasetDID("datasetID").
				WithGovAddr("addr").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "Valid governance VC without id",
			vc: credential.NewGovernanceVC().
				WithDatasetDID("datasetID").
				WithGovAddr("addr").
				WithIssuanceDate(time.Now().UTC()),
		},
		{
			name: "Valid governance VC without issuance date",
			vc: credential.NewGovernanceVC().
				WithID("id").
				WithDatasetDID("datasetID").
				WithGovAddr("addr"),
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
					buf, err := generator.Generate()

					Convey("Then the governance VC should be generated", func() {
						if test.wantErr != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, test.wantErr.Error())
							So(buf, ShouldBeNil)
						} else {
							So(buf, ShouldNotBeNil)
							So(err, ShouldBeNil)
						}
					})
				})
			})
		})
	}
}
