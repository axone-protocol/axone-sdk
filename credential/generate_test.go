//nolint:lll
package credential_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/axone-protocol/axone-sdk/credential"
	"github.com/axone-protocol/axone-sdk/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

func TestGenerator_Generate(t *testing.T) {
	t.Run("without parser", func(t *testing.T) {
		Convey("Given a credential generator with mocked descriptor", t, func() {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockDescriptor := testutil.NewMockDescriptor(controller)
			mockDescriptor.EXPECT().Generate().Return(nil, nil).Times(1)
			generator := credential.New(mockDescriptor)

			Convey("When generating a credential", func() {
				vc, err := generator.Generate()
				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "no parser provided")
					So(vc, ShouldBeNil)
				})
			})
		})
	})

	t.Run("error generate", func(t *testing.T) {
		Convey("Given a credential generator with mocked descriptor", t, func() {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockDescriptor := testutil.NewMockDescriptor(controller)
			mockDescriptor.EXPECT().Generate().Return(nil, errors.New("failed")).Times(1)
			generator := credential.New(mockDescriptor)

			Convey("When generating a credential", func() {
				vc, err := generator.Generate()
				Convey("Then an error should be returned", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, credential.NewVCError(credential.ErrGenerate, errors.New("failed")).Error())
					So(vc, ShouldBeNil)
				})
			})
		})
	})

	t.Run("with signer", func(t *testing.T) {
		Convey("Given a credential generator with mocked descriptor", t, func() {
			controller := gomock.NewController(t)
			defer controller.Finish()

			buf := bytes.NewBufferString(`{"@context":["https://www.w3.org/2018/credentials/v1"],"type":["VerifiableCredential"],"credentialSubject":{"id":"did:example:123"}}`)
			mockDescriptor := testutil.NewMockDescriptor(controller)
			mockDescriptor.EXPECT().Generate().Return(buf, nil).Times(1)
			mockDescriptor.EXPECT().IssuedAt().Times(1)
			mockDescriptor.EXPECT().ProofPurpose().Return("proof").Times(1)

			mockSigner := testutil.NewMockKeyring(controller)
			mockSigner.EXPECT().Sign(gomock.Any()).Return([]byte("signature"), nil).Times(1)
			mockSigner.EXPECT().Alg().AnyTimes()
			mockSigner.EXPECT().DIDKeyID().Return("did:example:123#123").Times(1)

			loader, _ := testutil.MockDocumentLoader()

			generator := credential.New(mockDescriptor,
				credential.WithParser(credential.NewDefaultParser(loader)),
				credential.WithSigner(mockSigner),
			)

			Convey("When generating a credential", func() {
				vc, err := generator.Generate()
				Convey("Then VC should be returned with proofs", func() {
					So(err, ShouldBeNil)
					So(vc, ShouldNotBeNil)
					So(len(vc.Proofs), ShouldEqual, 1)
					So(vc.Proofs[0]["verificationMethod"], ShouldEqual, "did:example:123#123")
					So(vc.Proofs[0]["proofPurpose"], ShouldEqual, "proof")
				})
			})
		})
	})

	t.Run("without signer", func(t *testing.T) {
		Convey("Given a credential generator with mocked descriptor", t, func() {
			controller := gomock.NewController(t)
			defer controller.Finish()

			buf := bytes.NewBufferString(`{"@context":["https://www.w3.org/2018/credentials/v1"],"type":["VerifiableCredential"],"credentialSubject":{"id":"did:example:123"}}`)
			mockDescriptor := testutil.NewMockDescriptor(controller)
			mockDescriptor.EXPECT().Generate().Return(buf, nil).Times(1)

			loader, _ := testutil.MockDocumentLoader()

			generator := credential.New(mockDescriptor,
				credential.WithParser(credential.NewDefaultParser(loader)),
			)

			Convey("When generating a credential", func() {
				vc, err := generator.Generate()
				Convey("Then VC should be returned without proofs", func() {
					So(err, ShouldBeNil)
					So(vc, ShouldNotBeNil)
					So(len(vc.Proofs), ShouldEqual, 0)
				})
			})
		})
	})
}
