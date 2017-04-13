// Package api exposes a low-level interface to the MyRadio API.
//
// It exposes the Requester interface for types that represent
// connections to the API, methods for constructing Requesters, and
// functions and methods for using Requesters to make requests.
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request represents an API request being built.
type Request struct {
	// The endpoint, as a suffix of the API root URL.
	Endpoint string
	// The set of mixins to use.
	Mixins []string
	// The map of parameters to use.
	Params map[string][]string
}

// Get constructs a new request for the given endpoint.
func Get(endpoint string) *Request {
	r := Request{}
	r.Endpoint = endpoint
	return &r
}

// Get constructs a new request for the endpoint constructed with the given format string and parameters.
func Getf(format string, params ...interface{}) *Request {
	return Get(fmt.Sprintf(format, params...))
}

// Response represents the result of an API request.
type Response struct {
	raw *json.RawMessage
	err error
}

// IsEmpty checks whether the response payload is present, but empty.
func (r *Response) IsEmpty() bool {
	if r.err == nil {
		return false
	}

	if r.raw == nil {
		return true
	}

	// Check for 'empty' JSON payloads.
	bs, err := r.raw.MarshalJSON()
	if err != nil {
		return false
	}

	if len(bs) != 2 {
		return false
	}

	if bs[0] == '[' && bs[1] == ']' {
		return true
	}

	if bs[0] == '{' && bs[1] == '}' {
		return true
	}

	return false
}

// JSON returns r as raw JSON.
func (r *Response) JSON() (*json.RawMessage, error) {
	return r.raw, r.err
}

// Into unmarshals the response r into in.
func (r *Response) Into(in interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.Unmarshal(*r.raw, in)
}

// Requester is the type of anything that can handle an API request.
type Requester interface {
	// Do fulfils an API request.
	Do(r *Request) *Response
}

// authedRequester answers API requests by making an authed API call.
type authedRequester struct {
	apikey  string
	baseurl url.URL
}

// NewRequester creates a new 'live' requester.
func NewRequester(apikey string, url url.URL) Requester {
	return &authedRequester{
		apikey:  apikey,
		baseurl: url,
	}
}

// Do fulfils an API request.
func (s *authedRequester) Do(r *Request) *Response {
	urlParams := url.Values{
		"api_key": []string{s.apikey},
	}
	if len(r.Mixins) > 0 {
		urlParams.Add("mixins", strings.Join(r.Mixins, ","))
	}
	for k, vs := range r.Params {
		for _, v := range vs {
			urlParams.Add(k, v)
		}
	}

	theurl := s.baseurl
	theurl.Path += r.Endpoint
	theurl.RawQuery = urlParams.Encode()
	req, err := http.NewRequest("GET", theurl.String(), nil)
	if err != nil {
		return &Response{err: err}
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &Response{err: err}
	}
	if res.StatusCode != 200 {
		return &Response{err: fmt.Errorf(r.Endpoint + fmt.Sprintf(" Not ok: HTTP %d", res.StatusCode))}
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Response{err: err}
	}
	var response struct {
		Status  string
		Payload *json.RawMessage
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return &Response{err: err}
	}
	if response.Status != "OK" {
		return &Response{err: fmt.Errorf(r.Endpoint + fmt.Sprintf(" Response not OK: %v", response))}
	}
	return &Response{raw: response.Payload, err: nil}
}

// mockRequester answers API requests by returning some stock response.
type mockRequester struct {
	message *json.RawMessage
}

// MockRequester creates a new mocked requester.
func MockRequester(message *json.RawMessage) Requester {
	return &mockRequester{message: message}
}

// Do pretends to fulfil an API request, but actually returns the mockRequester's stock response.
func (s *mockRequester) Do(r *Request) *Response {
	return &Response{raw: s.message, err: nil}
}
