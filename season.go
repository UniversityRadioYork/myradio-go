package myradio

import (
	"encoding/json"
	"fmt"
)

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
