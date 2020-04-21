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
}

// Get retrieves the data for a single podcast from MyRadio given it's ID.
// This consumes one API request.
func (s *Session) Get(id int) (podcast *Podcast, err error) {
	err = s.getf("/podcast/%d", id).Into(&podcast)
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
	rq.Params["include_pending"] = []string{"0"}
	rs := s.do(rq)

	if err := rs.Into(&podcasts); err != nil {
		return nil, err
	}
	return

}
