package myradio

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type ShowMeta struct {
	ShowID        int    `json:"show_id"`
	Title         string `json:"title"`
	Credits       string `json:"credits"`
	Description   string `json:"description"`
	ShowTypeID    int    `json:"show_type_id"`
	Season        Link   `json:"seasons"`
	EditLink      Link   `json:"editlink"`
	ApplyLink     Link   `json:"applylink"`
	MicroSiteLink Link   `json:"micrositelink"`
	Photo         string `json:"photo"`
}

type Link struct {
	Display string      `json:"display"`
	Value   interface{} `json:"value"`
	Title   string      `json:"title"`
	URL     string      `json:"url"`
}

func (s *Session) GetSearchMeta(term string) (*[]ShowMeta, error) {

	q := url.QueryEscape(term)

	data, err := s.apiRequest(fmt.Sprintf("/show/searchmeta/%s", q), []string{})

	if err != nil {
		return nil, err
	}

	var shows []ShowMeta

	err = json.Unmarshal(*data, &shows)

	if err != nil {
		return nil, err
	}

	return &shows, nil

}
