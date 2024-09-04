package template

import (
	"errors"
	"testing"
	"time"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/testutil"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPublicationDescriptor_Generate(t *testing.T) {
	tests := []struct {
		name    string
		vc      credential.Descriptor
		wantErr error
		check   func(*verifiable.Credential)
	}{
		{
			name: "Valid publication VC",
			vc: NewPublication().
				WithID("id").
				WithDatasetDID("datasetID").
				WithDatasetURI("datasetURI").
				WithStorageDID("storageID").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldEqual, "https://w3id.org/axone/ontology/v4/schema/credential/digital-resource/publication/id")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "storageID")
				So(vcSubject(vc).CustomFields["hasIdentifier"], ShouldEqual, "datasetURI")
				So(vcSubject(vc).CustomFields["servedBy"], ShouldEqual, "storageID")
				So(vc.Issued.Time, ShouldEqual, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		{
			name: "Valid publication VC with default issuance date and default id",
			vc: NewPublication().
				WithDatasetDID("datasetID").
				WithDatasetURI("datasetURI").
				WithStorageDID("storageID"),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldStartWith, "https://w3id.org/axone/ontology/v4/schema/credential/digital-resource/publication/")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "storageID")
				So(vcSubject(vc).CustomFields["hasIdentifier"], ShouldEqual, "datasetURI")
				So(vcSubject(vc).CustomFields["servedBy"], ShouldEqual, "storageID")
				So(vc.Issued.Time, ShouldHappenWithin, time.Second, time.Now().UTC())
			},
		},
		{
			name: "Missing dataset DID",
			vc: NewPublication().
				WithID("id").
				WithDatasetURI("datasetURI").
				WithStorageDID("storageID").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: credential.NewVCError(credential.ErrGenerate, errors.New("dataset DID is required")),
		},
		{
			name: "Missing storage DID",
			vc: NewPublication().
				WithID("id").
				WithDatasetDID("datasetID").
				WithDatasetURI("datasetURI").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: credential.NewVCError(credential.ErrGenerate, errors.New("storage DID is required")),
		},
		{
			name: "Missing dataset URI",
			vc: NewPublication().
				WithID("id").
				WithDatasetDID("datasetID").
				WithStorageDID("storageID").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: credential.NewVCError(credential.ErrGenerate, errors.New("dataset URI is required")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Convey("Given a credential generator", t, func() {
				docLoader, err := testutil.MockDocumentLoader()
				So(err, ShouldBeNil)

				parser := credential.NewDefaultParser(docLoader)
				generator := credential.New(test.vc).
					WithParser(parser)

				Convey("When a publication VC is generated", func() {
					vc, err := generator.Generate()

					Convey("Then the publication VC should be generated", func() {
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
