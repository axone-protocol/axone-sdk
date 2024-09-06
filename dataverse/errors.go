package dataverse

import "fmt"

type MessageError string

const (
	ErrNoResult    MessageError = "no result found in binding"
	ErrVarNotFound MessageError = "variable not found in binding result"
	ErrType        MessageError = "variable result type mismatch in binding result"
)

type DVError struct {
	message MessageError
	detail  error
}

func (e *DVError) Error() string {
	if e.detail == nil {
		return fmt.Sprintf("%v", e.message)
	}
	return fmt.Sprintf("%v: %v", e.message, e.detail)
}

func NewDVError(message MessageError, detail error) error {
	return &DVError{
		message: message,
		detail:  detail,
	}
}
