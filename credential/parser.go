package credential

import (
    "fmt"
    "github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"

    "github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
    "github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
    "github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
    "github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2020"
    "github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
    "github.com/hyperledger/aries-framework-go/pkg/vdr"
    "github.com/hyperledger/aries-framework-go/pkg/vdr/key"
)

type Claim interface {
    From(vc *verifiable.Credential) error
}

type Parser[T Claim] interface {
    ParseSigned(raw []byte) (T, error)
}

func parseWithSignVerification(raw []byte) (*verifiable.Credential, error) {
    publicKeyFetcher := verifiable.NewVDRKeyResolver(vdr.New(vdr.WithVDR(key.New()))).PublicKeyFetcher()

    vc, err := verifiable.ParseCredential(
        raw,
        verifiable.WithJSONLDValidation(),
        verifiable.WithPublicKeyFetcher(func(issuerID, keyID string) (*verifier.PublicKey, error) {
            // HACK: as the publicKeyFetcher doesn't manage `EcdsaSecp256k1VerificationKey2019` as verification method
            // we got to manage it ourselves.
            pubKey, err := mayResolveSecp256k1PubKey(issuerID, keyID)
            if err != nil {
                return nil, err
            }

            if pubKey != nil {
                return pubKey, nil
            }

            return publicKeyFetcher(issuerID, keyID)
        }),
        verifiable.WithEmbeddedSignatureSuites(
            ed25519signature2018.New(suite.WithVerifier(ed25519signature2018.NewPublicKeyVerifier())),
            ed25519signature2020.New(suite.WithVerifier(ed25519signature2020.NewPublicKeyVerifier())),
            ecdsasecp256k1signature2019.New(suite.WithVerifier(ecdsasecp256k1signature2019.NewPublicKeyVerifier())),
        ),
        verifiable.WithJSONLDDocumentLoader(a.documentLoader),
    )
    if err != nil {
        return nil, err
    }
}

func extractCustomStringClaim(claim *verifiable.Subject, key string) (string, error) {
    field, ok := claim.CustomFields[key]
    if !ok {
        return "", fmt.Errorf("key '%s' not found", key)
    }

    strField, ok := field.(string)
    if !ok {
        return "", fmt.Errorf("key '%s' is not a string", key)
    }
    return strField, nil
}
