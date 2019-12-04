package myradio

import (
	"fmt"
	"github.com/UniversityRadioYork/myradio-go/api"
	"strings"
)

// Album contains information about an album in the URY track database.
type Album struct {
	// ID is the unique database ID of the album.
	ID uint64 `json:"recordid"`

	// Title is the title of the track.
	Title string `json:"title"`
	// Artist is the primary credited artist of the track.
	Artist string `json:"artist"`

	// DateAdded is the date on which the album entered the MyRadio library.
	DateAdded string `json:"date_added"`
	// DateReleased is the date on which the album was released.
	DateReleased string `json:"date_released"`
	// LastModified is the date on which the album was last modified.
	LastModified string `json:"last_modified"`

	// CDID is the ID of the CD, if this track comes from one.
	CDID string `json:"cdid"`

	// Location is the location of the physical copy of this album, if any.
	Location string `json:"location"`
	// ShelfLetter is the shelf on which the physical copy resides, if any.
	ShelfLetter string `json:"shelf_letter"`
	// ShelfNumber is the position on the shelf on which the physical copy resides, if any.
	ShelfNumber uint64 `json:"shelf_number"`

	// Format is a single-character code identifying the physical format.
	Format string `json:"format"`
	// Medium is a single-character code identifying the physical medium.
	Medium string `json:"media"`

	// AddingMember is the ID of the member who added this album.
	AddingMember uint64 `json:"member_add"`
	// EditingMember is the ID of the member who last modified this album.
	EditingMember uint64 `json:"member_edit"`

	// RecordLabel is the record label responsible for this album.
	RecordLabel string `json:"record_label"`

	// Status is the digitisation status code for this album.
	Status string `json:"status"`
}

// Track contains information about a track in the URY track database.
type Track struct {
	// ID is the unique database ID of the track.
	ID uint64 `json:"trackid"`

	// Title is the title of the track.
	Title string `json:"title"`
	// Artist is the primary credited artist of the track.
	Artist string `json:"artist"`
	// Type is the type ('central' etc.) of the track.
	Type string `json:"type"`
	// Length is the length of the track, in hours:minutes:seconds.
	Length string `json:"length"`
	// Intro is length of the track's intro, in seconds.
	Intro uint64 `json:"intro"`
	// IsClean is true if this track is clean (no expletives).
	IsClean bool `json:"clean"`
	// IsDigitised is true if this track is available in the playout system.
	IsDigitised bool `json:"digitised"`
}

// GetAlbum tries to get the Album for the given Track.
// This consumes one API request.
func (t *Track) GetAlbum(s *Session) (*Album, error) {
	return s.GetTrackAlbum(t.ID)
}

// LengthSec returns the track's length in seconds.
// Returns an error if the track's length is ill-formed.
// This consumes no API requests.
func (t *Track) LengthSec() (uint64, error) {
	var hours, minutes, seconds uint64

	_, err := fmt.Sscan(strings.Replace(t.Length, ":", " ", -1), &hours, &minutes, &seconds)
	if err != nil {
		return 0, err
	}

	return (hours * 60 * 60) + (minutes * 60) + seconds, nil
}

// LengthUsec returns the track's length in microseconds.
// This is not precise, as it is derived from the length in seconds.
// Consider estimating the correct length from the track file itself.
// Returns an error if the track's length is ill-formed.
// This consumes no API requests.
func (t *Track) LengthUsec() (uint64, error) {
	secs, err := t.LengthSec()
	if err != nil {
		return 0, err
	}

	return secs * 1000000, nil
}

// IntroUsec returns the track's intro in microseconds.
// This consumes no API requests.
func (t *Track) IntroUsec() uint64 {
	return t.Intro * 1000000
}

// GetTrack tries to get the Track with the given ID.
// Track IDs are unique, so we do not need the record ID.
// This consumes one API request.
func (s *Session) GetTrack(trackid uint64) (track *Track, err error) {
	err = s.getf("/track/%d", trackid).Into(&track)
	return
}

// GetTrackTitle tries to get the title of the track with the given ID.
// This consumes one API request.
func (s *Session) GetTrackTitle(trackid uint64) (title string, err error) {
	err = s.getf("/track/%d/title", trackid).Into(&title)
	return
}

// GetTrackAlbum tries to get the Album of the track with the given ID.
// This consumes one API request.
func (s *Session) GetTrackAlbum(trackid uint64) (album *Album, err error) {
	err = s.getf("/track/%d/album", trackid).Into(&album)
	return
}

// The parameters for searching for tracks.
// All of these are optional.
type TrackSearchParams struct {
	Title    string
	Artist   string
	RecordID uint64
	// The cleanliness of the track - either "y", "n", or "u"
	Digitised bool
	Clean     string
	Precise   bool
	Limit     uint64
	// Either "id" (default), "title", or "random"
	Sort             string
	ITonesPlaylistID string
}

func (s *Session) SearchTracks(p TrackSearchParams) (tracks []Track, err error) {
	args := make(map[string][]string)
	if p.Title != "" {
		args["title"] = []string{p.Title}
	}
	if p.Artist != "" {
		args["artist"] = []string{p.Artist}
	}
	if p.RecordID != 0 {
		args["recordid"] = []string{fmt.Sprint(p.RecordID)}
	}
	if p.Digitised != true {
		args["digitised"] = []string{fmt.Sprint(p.Digitised)}
	}
	if p.Clean != "" {
		args["clean"] = []string{p.Clean}
	}
	if p.Limit != 0 {
		args["limit"] = []string{fmt.Sprint(p.Limit)}
	}
	if p.Sort != "" {
		args["sort"] = []string{p.Sort}
	}
	if p.ITonesPlaylistID != "" {
		args["itonesplaylistid"] = []string{p.ITonesPlaylistID}
	}

	req := api.NewRequest("/track/search")
	req.Params = args
	err = s.do(req).Into(&tracks)

	return
}
