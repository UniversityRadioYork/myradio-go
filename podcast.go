package myradio

// Podcast represents a podcast media item.
type Podcast struct {
	PodcastID     int `json:"podcast_id"`
	Title         string
	Description   string
	Status        string
	Time          Time
	Photo         string
	File          string `json:"file"`
	WebFile       string
	EditLink      Link `json:"editlink"`
	MicrositeLink Link `json:"micrositelink"`
}

// Get retrieves the data for a single podcast from MyRadio given it's ID.
// This consumes one API request.
func (s *Session) Get(id int) (podcast *Podcast, err error) {
	err = s.getf("/podcast/%d", id).Into(&podcast)
	return
}

// GetAllPodcasts retrieves the latest podcasts from MyRadio.
// This consumes one API request.
func (s *Session) GetAllPodcasts() (podcasts []Podcast, err error) {
	err = s.get("/podcast/allpodcasts").Into(&podcasts)
	return
}
