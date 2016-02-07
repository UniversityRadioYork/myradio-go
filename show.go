package myradio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
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
	Title   string      `json:"title,omitempty"`
	URL     string      `json:"url"`
}

type Season struct {
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
	SeasonID      int    `json:"season_id"`
	SeasonNum     int    `json:"season_num"`
	Submitted     time.Time   `json:"submitted"`
	RequestedTime string      `json:"requested_time"`
	FirstTime     time.Time   `json:"first_time"`
	NumEpisodes   Link        `json:"num_episodes"`
	AllocateLink  Link        `json:"allocatelink"`
	RejectLink    Link        `json:"rejectlink"`
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

func (s *Session) GetShow(id int) (*ShowMeta, error) {

	data, err := s.apiRequest(fmt.Sprintf("/show/%d", id), []string{})

	if err != nil {
		return nil, err
	}

	var show ShowMeta

	err = json.Unmarshal(*data, &show)

	if err != nil {
		return nil, err
	}

	return &show, nil

}

func (s *Session) GetSeasons(id int) (*[]Season, error) {

	data, err := s.apiRequest(fmt.Sprintf("/show/%d/allseasons", id), []string{})

	if err != nil {
		return nil, err
	}

	var seasons []Season

	err = json.Unmarshal(*data, &seasons)

	if err != nil {
		return nil, err
	}

	return &seasons, nil

}
