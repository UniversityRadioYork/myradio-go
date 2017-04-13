package myradio

import (
	"encoding/json"
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

// do fulfils, a request for the given endpoint.
func (s *Session) do(r *api.Request) *api.Response {
	return s.requester.Do(r)
}

// get creates, and fulfils, a GET request for the given endpoint.
func (s *Session) get(endpoint string) *api.Response {
	return s.do(api.Get(endpoint))
}

// get creates, and fulfils, a GET request for the endpoint created by
// the given format string and parameters.
func (s *Session) getf(format string, params ...interface{}) *api.Response {
	return s.do(api.Getf(format, params...))
}

// NewSessionFromKeyFile tries to open a Session with the key from an API key file.
func NewSessionFromKeyFile() (*Session, error) {
	apikey, err := api.GetAPIKey()
	if err != nil {
		return nil, err
	}

	return NewSession(apikey)
}
