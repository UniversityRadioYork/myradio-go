package myradio

import (
	"fmt"
	"net/url"
)

// Credit represents a show credit associating a user with a show.
type Credit struct {
	Type     int  `json:"type"`
	MemberID int  `json:"memberid"`
	User     User `json:"User"`
}

// ShowMeta represents a show in the MyRadio schedule.
// A MyRadio show contains seasons, each containing timeslots.
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

// Link represents a MyRadio action link.
type Link struct {
	Display string      `json:"display"`
	Value   interface{} `json:"value"`
	Title   string      `json:"title,omitempty"`
	URL     string      `json:"url"`
}

// GetSearchMeta retrieves all shows whose metadata matches a given search term.
// This consumes one API request.
func (s *Session) GetSearchMeta(term string) (shows []ShowMeta, err error) {
	q := url.QueryEscape(term)

	err = s.apiRequestInto(&shows, fmt.Sprintf("/show/searchmeta/%s", q), []string{})
	return
}

// GetShow retrieves the show with the given ID.
// This consumes one API request.
func (s *Session) GetShow(id int) (show *ShowMeta, err error) {
	err = s.apiRequestInto(&show, fmt.Sprintf("/show/%d", id), []string{})
	return
}

// GetSeasons retrieves the seasons of the show with the given ID.
// This consumes one API request.
func (s *Session) GetSeasons(id int) (seasons []Season, err error) {
	if err = s.apiRequestInto(&seasons, fmt.Sprintf("/show/%d/allseasons", id), []string{}); err != nil {
		return
	}

	for i := range seasons {
		err = seasons[i].populateSeasonTimes()
		if err != nil {
			return
		}
	}

	return
}
