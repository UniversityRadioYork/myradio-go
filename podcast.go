package myradio

import (
	"strconv"

	"github.com/UniversityRadioYork/myradio-go/api"
)

// Podcast represents a podcast media item.
type Podcast struct {
	PodcastID     int `json:"podcast_id"`
	Title         string
	Description   string
	Status        string
	Time          Time
	Photo         string
	File          string `json:"uri"`
	EditLink      Link   `json:"editlink"`
	MicrositeLink Link   `json:"micrositelink"`
	Show          *ShowMeta
}

// Get retrieves the data for a single podcast from MyRadio given it's ID.
// This consumes one API request.
func (s *Session) Get(id int) (podcast *Podcast, err error) {
	err = s.getf("/podcast/%d", id).Into(&podcast)
	return
}

// Get retrieves the data for a single podcast, and its associated show.
// This only consumes one API request.
func (s *Session) GetPodcastWithShow(id int) (podcast *Podcast, err error) {
	req := api.NewRequestf("/podcast/%d", id)
	req.Mixins = []string{"show"}
	err = s.do(req).Into(&podcast)
	return
}

// GetAllPodcasts retrieves the latest podcasts from MyRadio.
// This consumes one API request.
func (s *Session) GetAllPodcasts(numResults int, page int, includeSuspended bool) (podcasts []Podcast, err error) {

	rq := api.NewRequest("/podcast/allpodcasts")
	rq.Params["num_results"] = []string{strconv.Itoa(numResults)}
	rq.Params["page"] = []string{strconv.Itoa(page)}
	suspended := "0"
	if includeSuspended {
		suspended = "1"
	}
	rq.Params["include_suspended"] = []string{suspended}
	rs := s.do(rq)

	if err := rs.Into(&podcasts); err != nil {
		return nil, err
	}
	return

}

// GetAllShowPodcasts returns all podcasts linked to the given show.
func (s *Session) GetAllShowPodcasts(id int) (result []Podcast, err error) {
	err = s.getf("/show/%d/allpodcasts", id).Into(&result)
	return
}

// GetPodcastSearchMeta retrieves all Podcasts whose metadata matches a given search term.
// This consumes one API request.
func (s *Session) GetPodcastSearchMeta(term string) (podcasts []Podcast, err error) {
	err = s.getf("/podcast/searchmeta/%s", term).Into(&podcasts)
	return
}
