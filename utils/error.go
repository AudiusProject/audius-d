package utils

import "fmt"

type AudiusError struct {
	Message string
}

func NewAudiusError(msg string) *AudiusError {
	return &AudiusError{Message: msg}
}

func (ce AudiusError) Error() string {
	return fmt.Sprintf("audius-d error: %s", ce.Message)
}
