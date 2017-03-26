package myradio

import (
	"encoding/json"
	"fmt"
	"time"
)

// Season represents a season in the MyRadio schedule.
// A MyRadio season contains timeslots.
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

// isScheduled returns whether the Season has been scheduled.
// This consumes no API requests.
func (s *Season) isScheduled() bool {
	return s.FirstTimeRaw != "Not Scheduled"
}

// populateSeasonTimes sets the times for the given Season given their raw values.
func (s *Season) populateSeasonTimes() (err error) {
	if s.isScheduled() {
		s.FirstTime, err = parseShortTime(s.FirstTimeRaw)
		if err != nil {
			return
		}
	}
	s.Submitted, err = parseShortTime(s.SubmittedRaw)
	return
}

// GetSeason retrieves the season with the given ID.
// This consumes one API request.
func (s *Session) GetSeason(id int) (season Season, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/season/%d/", id), []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &season)
	if err != nil {
		return
	}
	err = season.populateSeasonTimes()
	return
}

// GetTimeslotsForSeason retrieves all timeslots for the season with the given ID.
// This consumes one API request.
func (s *Session) GetTimeslotsForSeason(id int) (timeslots []Timeslot, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/season/%d/alltimeslots/", id), []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &timeslots)
	if err != nil {
		return
	}
	for k := range timeslots {
		err = timeslots[k].populateTimeslotTimes()
		if err != nil {
			return
		}
	}
	return
}

// GetAllSeasonsInLatestTerm gets all seasons in the most recent term.
// This consumes one API request.
func (s *Session) GetAllSeasonsInLatestTerm() (seasons []Season, err error) {
	data, err := s.apiRequest("/season/allseasonsinlatestterm/", []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &seasons)
	if err != nil {
		return
	}
	for k := range seasons {
		err = seasons[k].populateSeasonTimes()
		if err != nil {
			return
		}
	}
	return
}
