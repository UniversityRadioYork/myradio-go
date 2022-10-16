package api

import "fmt"

// Error is the error type returned when API requests result in an error.
type Error struct {
	Endpoint string
	Code     int
	Payload  string
}

func (a Error) Error() string {
	return fmt.Sprintf("%s NOT OK: %d: %s", a.Endpoint, a.Code, a.Payload)
}
