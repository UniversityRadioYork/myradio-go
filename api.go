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

// apiRequestWithParams conducts a GET request with custom parameters.
func (s *Session) apiRequestWithParams(endpoint string, mixins []string, params map[string][]string) (*json.RawMessage, error) {
	return s.requester.request(endpoint, mixins, params)
}

// apiRequest conducts a GET request without custom parameters.
func (s *Session) apiRequest(endpoint string, mixins []string) (*json.RawMessage, error) {
	return s.apiRequestWithParams(endpoint, mixins, map[string][]string{})
}
