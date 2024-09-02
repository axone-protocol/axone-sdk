package credential

import "github.com/piprate/json-gold/ld"

//nolint:revive
func NewCredentialParser(documentLoader ld.DocumentLoader) *credentialParser {
	return &credentialParser{documentLoader: documentLoader}
}
