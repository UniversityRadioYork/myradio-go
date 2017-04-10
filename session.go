package myradio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"github.com/UniversityRadioYork/myradio-go/api"
)

// Session represents an open API session.
type Session struct {
	requester api.Requester
}

// NewSession constructs a new Session with the given API key.
func NewSession(apikey string) (*Session, error) {
	url, err := url.Parse(`https://ury.york.ac.uk/api/v2`)
	if err != nil {
		return nil, err
	}
	return &Session{requester: api.NewRequester(apikey, *url)}, nil
}

// MockSession creates a new mocked API session returning the JSON message stored in message.
func MockSession(message []byte) (*Session, error) {
	rm := json.RawMessage{}
	err := rm.UnmarshalJSON(message)
	if err != nil {
		return nil, err
	}
	return &Session{requester: api.MockRequester(&rm)}, nil
}

// get constructs a new request for the given endpoint.
func (s *Session) get(endpoint string) *api.Request {
	return api.Get(s.requester, endpoint)
}

// getf constructs a request whose endpoint is defined by a formatted string.
func (s *Session) getf(format string, params ...interface{}) *api.Request {
	return s.get(fmt.Sprintf(format, params...))
}
