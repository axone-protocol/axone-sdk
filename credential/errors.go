package credential

import "fmt"

type MessageError string

const (
	ErrExpired           MessageError = "verifiable credential expired"
	ErrIssued            MessageError = "verifiable credential issued in the future"
	ErrMissingProof      MessageError = "missing verifiable credential proof"
	ErrInvalidProof      MessageError = "invalid verifiable credential proof"
	ErrMalformedSubject  MessageError = "malformed verifiable credential subject"
	ErrExpectSingleClaim MessageError = "expect a single verifiable credential claim"
	ErrExtractClaim      MessageError = "failed to extract claim"
	ErrMalformed         MessageError = "malformed verifiable credential"
)

type VCError struct {
	message MessageError
	detail  error
}

func (e *VCError) Error() string {
	if e.detail == nil {
		return fmt.Sprintf("%v", e.message)
	}
	return fmt.Sprintf("%v: %v", e.message, e.detail)
}

func NewVCError(message MessageError, detail error) error {
	return &VCError{
		message: message,
		detail:  detail,
	}
}
