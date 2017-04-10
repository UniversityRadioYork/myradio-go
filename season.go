package myradio

import "time"

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
	if err = s.getf("/season/%d/", id).Into(&season); err != nil {
		return
	}

	err = season.populateSeasonTimes()

	return
}

// GetTimeslotsForSeason retrieves all timeslots for the season with the given ID.
// This consumes one API request.
func (s *Session) GetTimeslotsForSeason(id int) (timeslots []Timeslot, err error) {
	if err = s.getf("/season/%d/alltimeslots/", id).Into(&timeslots); err != nil {
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
	if err = s.get("/season/allseasonsinlatestterm/").Into(&seasons); err != nil {
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
