package myradio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// apiRequester is the type of anything that can handle an API request.
type apiRequester interface {
	// request sends out a full API request.
	request(endpoint string, mixins []string, params map[string][]string) (*json.RawMessage, error)
}

// authedRequester answers API requests by making an authed API call.
type authedRequester struct {
	apikey  string
	baseurl url.URL
}

// apiResponse provides the base structure of MyRadio API responses.
type apiResponse struct {
	Status  string
	Payload *json.RawMessage
}

// request fulfils an API request by making an authed API call.
func (s *authedRequester) request(endpoint string, mixins []string, params map[string][]string) (*json.RawMessage, error) {
	urlParams := url.Values{
		"api_key": []string{s.apikey},
	}
	if len(mixins) > 0 {
		urlParams.Add("mixins", strings.Join(mixins, ","))
	}
	for k, vs := range params {
		for _, v := range vs {
			urlParams.Add(k, v)
		}
	}

	theurl := s.baseurl
	theurl.Path += endpoint
	theurl.RawQuery = urlParams.Encode()
	req, err := http.NewRequest("GET", theurl.String(), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(endpoint + fmt.Sprintf(" Not ok: HTTP %d", res.StatusCode))
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response apiResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf(endpoint + fmt.Sprintf(" Response not OK: %v", response))
	}
	return response.Payload, nil
}

// mockRequester answers API requests by returning some stock response.
type mockRequester struct {
	message *json.RawMessage
}

// request pretends to fulfil an API request, but actually returns the mockRequester's stock response.
func (s *mockRequester) request(endpoint string, mixins []string, params map[string][]string) (*json.RawMessage, error) {
	return s.message, nil
}

// Session represents an open API session.
type Session struct {
	requester apiRequester
}

// NewSession constructs a new Session with the given API key.
func NewSession(apikey string) (*Session, error) {
	url, err := url.Parse(`https://ury.york.ac.uk/api/v2`)
	if err != nil {
		return nil, err
	}
	return &Session{
		requester: &authedRequester{
			apikey:  apikey,
			baseurl: *url,
		},
	}, nil
}

// MockSession creates a new mocked API session returning the JSON message stored in message.
func MockSession(message []byte) (*Session, error) {
	rm := json.RawMessage{}
	err := rm.UnmarshalJSON(message)
	if err != nil {
		return nil, err
	}
	return &Session{
		requester: &mockRequester{message: &rm},
	}, nil
}

// Struct request represents an API request being built.
type request struct {
	requester apiRequester
	endpoint  string
	mixins    []string
	params    map[string][]string
}

// mixin adds one or more mixins to an API request.
// It returns a pointer to the original request.
func (r *request) mixin(ms ...string) *request {
	r.mixins = append(r.mixins, ms...)
	return r
}

// param adds a parameter with key k and values v to an API request.
// It returns a pointer to the original request.
func (r *request) param(k string, vs ...string) *request {
	r.params[k] = vs
	return r
}

// do runs a request and returns the raw JSON and error.
func (r *request) do() (*json.RawMessage, error) {
	return r.requester.request(r.endpoint, r.mixins, r.params)
}

// into runs a request and tries to unmarshal the result into in.
func (r *request) into(in interface{}) error {
	data, aerr := r.do()
	if aerr != nil {
		return aerr
	}

	return json.Unmarshal(*data, in)
}

// get constructs a new request for the given endpoint.
func (s *Session) get(endpoint string) *request {
	r := request{}
	r.requester = s.requester
	r.endpoint = endpoint
	return &r
}

// getf constructs a request whose endpoint is defined by a formatted string.
func (s *Session) getf(format string, params ...interface{}) *request {
	return s.get(fmt.Sprintf(format, params...))
}
