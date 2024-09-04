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

func TestDatasetDescriptor_Generate(t *testing.T) {
	tests := []struct {
		name    string
		vc      credential.Descriptor
		wantErr error
		check   func(*verifiable.Credential)
	}{
		{
			name: "Valid dataset VC",
			vc: NewDataset().
				WithID("id").
				WithDatasetDID("datasetID").
				WithTitle("title").
				WithDescription("description").
				WithFormat("format").
				WithTags([]string{"tag1", "tag2"}).
				WithTopic("topic").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldEqual, "https://w3id.org/axone/ontology/v4/schema/credential/dataset/description/id")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["hasTitle"], ShouldEqual, "title")
				So(vcSubject(vc).CustomFields["hasDescription"], ShouldEqual, "description")
				So(vcSubject(vc).CustomFields["hasFormat"], ShouldEqual, "format")
				So(vcSubject(vc).CustomFields["hasTag"], ShouldResemble, []interface{}{"tag1", "tag2"})
				So(vcSubject(vc).CustomFields["hasTopic"], ShouldEqual, "topic")
				So(vc.Issued.Time, ShouldEqual, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		{
			name: "Valid dataset VC with default issuance date",
			vc: NewDataset().
				WithID("id").
				WithDatasetDID("datasetID").
				WithTitle("title").
				WithDescription("description").
				WithFormat("format").
				WithTags([]string{"tag1", "tag2"}).
				WithTopic("topic"),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldEqual, "https://w3id.org/axone/ontology/v4/schema/credential/dataset/description/id")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["hasTitle"], ShouldEqual, "title")
				So(vcSubject(vc).CustomFields["hasDescription"], ShouldEqual, "description")
				So(vcSubject(vc).CustomFields["hasFormat"], ShouldEqual, "format")
				So(vcSubject(vc).CustomFields["hasTag"], ShouldResemble, []interface{}{"tag1", "tag2"})
				So(vcSubject(vc).CustomFields["hasTopic"], ShouldEqual, "topic")
				So(vc.Issued.Time, ShouldHappenWithin, time.Second, time.Now().UTC())
			},
		},
		{
			name: "Valid dataset VC with generated id",
			vc: NewDataset().
				WithDatasetDID("datasetID").
				WithTitle("title").
				WithDescription("description").
				WithFormat("format").
				WithTags([]string{"tag1", "tag2"}).
				WithTopic("topic").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			check: func(vc *verifiable.Credential) {
				So(vc.ID, ShouldStartWith, "https://w3id.org/axone/ontology/v4/schema/credential/dataset/description/")
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["hasTitle"], ShouldEqual, "title")
				So(vcSubject(vc).CustomFields["hasDescription"], ShouldEqual, "description")
				So(vcSubject(vc).CustomFields["hasFormat"], ShouldEqual, "format")
				So(vcSubject(vc).CustomFields["hasTag"], ShouldResemble, []interface{}{"tag1", "tag2"})
				So(vcSubject(vc).CustomFields["hasTopic"], ShouldEqual, "topic")
				So(vc.Issued.Time, ShouldEqual, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
			},
		},
		{
			name: "Missing title",
			vc: NewDataset().
				WithID("id").
				WithDatasetDID("datasetID").
				WithDescription("description").
				WithFormat("format").
				WithTags([]string{"tag1", "tag2"}).
				WithTopic("topic").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: credential.NewVCError(credential.ErrGenerate, errors.New("title is required")),
		},
		{
			name: "Missing dataset DID",
			vc: NewDataset().
				WithID("id").
				WithTitle("title").
				WithDescription("description").
				WithFormat("format").
				WithTags([]string{"tag1", "tag2"}).
				WithTopic("topic").
				WithIssuanceDate(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			wantErr: credential.NewVCError(credential.ErrGenerate, errors.New("dataset DID is required")),
		},
		{
			name: "Valid dataset VC without optional value",
			vc: NewDataset().
				WithDatasetDID("datasetID").
				WithTitle("title"),
			check: func(vc *verifiable.Credential) {
				So(vcSubject(vc).ID, ShouldEqual, "datasetID")
				So(vc.Issuer.ID, ShouldEqual, "datasetID")
				So(vcSubject(vc).CustomFields["hasTitle"], ShouldEqual, "title")
				So(vcSubject(vc).CustomFields["hasDescription"], ShouldEqual, "")
				So(vcSubject(vc).CustomFields["hasFormat"], ShouldEqual, "")
				So(vcSubject(vc).CustomFields["hasTag"], ShouldResemble, []interface{}{})
				So(vcSubject(vc).CustomFields["hasTopic"], ShouldEqual, "")
			},
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

				Convey("When a dataset VC is generated", func() {
					vc, err := generator.Generate()

					Convey("Then the dataset VC should be generated", func() {
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
