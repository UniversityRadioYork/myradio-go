package myradio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Session represents an open API session.
type Session struct {
	apikey  string
	baseurl url.URL
}

// NewSession constructs a new Session with the given API key.
func NewSession(apikey string) (*Session, error) {
	url, err := url.Parse(`https://ury.york.ac.uk/api/v2`)
	if err != nil {
		return nil, err
	}
	return &Session{
		apikey:  apikey,
		baseurl: *url,
	}, nil
}

// apiResponse provides the base structure of MyRadio API responses.
type apiResponse struct {
	Status  string
	Payload *json.RawMessage
}

// apiRequestWithParams conducts a GET request with custom parameters.
func (s *Session) apiRequestWithParams(endpoint string, mixins []string, params map[string][]string) (*json.RawMessage, error) {
	urlParams := url.Values{
		"api_key": []string{s.apikey},
		"mixins":  []string{strings.Join(mixins, ",")},
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

// apiRequest conducts a GET request without custom parameters.
func (s *Session) apiRequest(endpoint string, mixins []string) (*json.RawMessage, error) {
	return s.apiRequestWithParams(endpoint, mixins, map[string][]string{})
}
