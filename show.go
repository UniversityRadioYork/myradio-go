package myradio

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"log"
)

type Credit struct {
	Type     int    `json:"type"`
	MemberID int    `json:"memberid"`
	User     Member `json:"User"`
}

// @TODO: Refactor this to something better named
type ShowMeta struct {
	ShowID        int      `json:"show_id"`
	Title         string   `json:"title"`
	CreditsString string   `json:"credits_string"`
	Credits       []Credit `json:"credits"`
	Description   string   `json:"description"`
	ShowTypeID    int      `json:"show_type_id"`
	Season        Link     `json:"seasons"`
	EditLink      Link     `json:"editlink"`
	ApplyLink     Link     `json:"applylink"`
	MicroSiteLink Link     `json:"micrositelink"`
	Photo         string   `json:"photo"`
}

type Link struct {
	Display string      `json:"display"`
	Value   interface{} `json:"value"`
	Title   string      `json:"title,omitempty"`
	URL     string      `json:"url"`
}

type Season struct {
	ShowMeta
	SeasonID      int    `json:"season_id"`
	SeasonNum     int    `json:"season_num"`
	SubmittedRaw  string `json:"submitted"`
	Submitted     time.Time
	RequestedTime string `json:"requested_time"`
	FirstTimeRaw  string `json:"first_time"`
	FirstTime     time.Time
	NumEpisodes   Link `json:"num_episodes"`
	AllocateLink  Link `json:"allocatelink"`
	RejectLink    Link `json:"rejectlink"`
}

func (s *Session) GetSearchMeta(term string) ([]ShowMeta, error) {

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

	return shows, nil

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

func (s *Session) GetSeasons(id int) ([]Season, error) {

	data, err := s.apiRequest(fmt.Sprintf("/show/%d/allseasons", id), []string{})

	if err != nil {
		return nil, err
	}

	var seasons []Season

	err = json.Unmarshal(*data, &seasons)

	if err != nil {
		return nil, err
	}

	for k, v := range seasons {

		seasons[k].FirstTime, err = time.Parse("02/01/2006 15:04", v.FirstTimeRaw)

		if err != nil {
			log.Print(err)
		}

		seasons[k].Submitted, err = time.Parse("02/01/2006 15:04", v.SubmittedRaw)

		if err != nil {
			log.Print(err)
		}

	}

	return seasons, nil

}
