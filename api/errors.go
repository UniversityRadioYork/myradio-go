package api

import "fmt"

// Error is the error type returned when API requests result in an error.
// Note: callers should use [errors.As] to check if an error is an API error,
// as the error may be wrapped in another error type.
type Error struct {
	// Endpoint was the MyRadio API endpoint that was called.
	Endpoint string
	// Code is the HTTP status code of the response from MyRadio.
	Code int
	// Payload is the raw error message from MyRadio.
	Payload string
}

func (a Error) Error() string {
	return fmt.Sprintf("%s NOT OK: %d: %s", a.Endpoint, a.Code, a.Payload)
}
