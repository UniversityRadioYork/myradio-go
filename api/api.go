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

// Requester is the type of anything that can handle an API request.
type Requester interface {
	// Do fulfils an API request.
	Do(r *Request) (*json.RawMessage, error)
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
func (s *authedRequester) Do(r *Request) (*json.RawMessage, error) {
	urlParams := url.Values{
		"api_key": []string{s.apikey},
	}
	if len(r.mixins) > 0 {
		urlParams.Add("mixins", strings.Join(r.mixins, ","))
	}
	for k, vs := range r.params {
		for _, v := range vs {
			urlParams.Add(k, v)
		}
	}

	theurl := s.baseurl
	theurl.Path += r.endpoint
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
		return nil, fmt.Errorf(r.endpoint + fmt.Sprintf(" Not ok: HTTP %d", res.StatusCode))
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response struct {
		Status  string
		Payload *json.RawMessage
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "OK" {
		return nil, fmt.Errorf(r.endpoint + fmt.Sprintf(" Response not OK: %v", response))
	}
	return response.Payload, nil
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
func (s *mockRequester) Do(r *Request) (*json.RawMessage, error) {
	return s.message, nil
}


// Request represents an API request being built.
type Request struct {
	requester Requester
	endpoint  string
	mixins    []string
	params    map[string][]string
}

// Mixin adds one or more mixins to an API request.
// It returns a pointer to the original request.
func (r *Request) Mixin(ms ...string) *Request {
	r.mixins = append(r.mixins, ms...)
	return r
}

// param adds a parameter with key k and values v to an API request.
// It returns a pointer to the original request.
func (r *Request) Param(k string, vs ...string) *Request {
	r.params[k] = vs
	return r
}

// Do runs a request and returns the raw JSON and error.
func (r *Request) Do() (*json.RawMessage, error) {
	return r.requester.Do(r)
}

// into runs a request and tries to unmarshal the result into in.
func (r *Request) Into(in interface{}) error {
	data, aerr := r.Do()
	if aerr != nil {
		return aerr
	}

	return json.Unmarshal(*data, in)
}

// Get constructs a new request for the given endpoint.
func Get(rq Requester, endpoint string) *Request {
	r := Request{}
	r.requester = rq
	r.endpoint = endpoint
	return &r
}
