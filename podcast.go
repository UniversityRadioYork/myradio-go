package myradio

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

// GetAllPodcasts retrieves the latest podcasts from MyRadio.
// This consumes one API request.
func (s *Session) GetAllPodcasts() (podcasts []Podcast, err error) {
	err = s.get("/podcast/allpodcasts").Into(&podcasts)
	return
}
