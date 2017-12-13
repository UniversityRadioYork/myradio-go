package myradio

import "GitHub.com/UniversityRadioYork/myradio-go/api"

// Podcast represents a podcast media item.
type Podcast struct {
	PodcastID     int `json:"podcast_id"`
	Title         string
	Description   string
	Status        string
	Time          Time
	Photo         string
	EditLink      Link `json:"editlink"`
	MicrositeLink Link `json:"micrositelink"`
}

// Get retrieves the data for a single podcast from MyRadio given it's ID.
// This consumes one API request.
func (s *Session) Get(id int) (podcast *Podcast, err error) {
	rq := api.NewRequestf("/podcast/%s", id)
	if err = s.do(rq).Into(&podcast); err != nil {
		return
	}
	return
}

// GetAllPodcasts retrieves the latest podcasts from MyRadio.
// This consumes one API request.
func (s *Session) GetAllPodcasts() (podcasts []Podcast, err error) {
	err = s.get("/podcast/allpodcasts").Into(&podcasts)
	return
}
