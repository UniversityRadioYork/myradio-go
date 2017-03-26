package myradio

import (
	"encoding/json"
	"fmt"
	"time"
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
	season.FirstTime, err = time.Parse("02/01/2006 15:04", season.FirstTimeRaw)
	if err != nil {
		return
	}
	season.Submitted, err = time.Parse("02/01/2006 15:04", season.SubmittedRaw)
	return
}

// GetTimeslotsForSeason retrieves all timeslots for the season with the given ID.
// This consumes one API request.
func (s *Session) GetTimeslotsForSeason(id int) (timeslots []Timeslot, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/season/%d/alltimeslots/", id), []string{})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(*data, &timeslots)
	if err != nil {
		return nil, err
	}
	for k, v := range timeslots {
		timeslots[k].Time = time.Unix(v.TimeRaw, 0)
		if v.FirstTimeRaw != "Not Scheduled" {
			timeslots[k].FirstTime, err = time.Parse("02/01/2006 15:04", v.FirstTimeRaw)
			if err != nil {
				return
			}
		}
		timeslots[k].Submitted, err = time.Parse("02/01/2006 15:04", v.SubmittedRaw)
		if err != nil {
			return
		}
		timeslots[k].StartTime, err = time.Parse("02/01/2006 15:04", v.StartTimeRaw)
		if err != nil {
			return
		}
		timeslots[k].Duration, err = parseDuration(v.DurationRaw)
		if err != nil {
			return
		}
	}
	return
}

// GetAllSeasonsInLatestTerm gets all seasons in the most recent term.
// It consumes one API request.
func (s *Session) GetAllSeasonsInLatestTerm() (seasons []Season, err error) {
	data, err := s.apiRequest("/season/allseasonsinlatestterm/", []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &seasons)
	if err != nil {
		return
	}
	for k, season := range seasons {
		if season.FirstTimeRaw != "Not Scheduled" {
			seasons[k].FirstTime, err = time.Parse("02/01/2006 15:04", season.FirstTimeRaw)
			if err != nil {
				return
			}
		}
		if season.FirstTimeRaw != "Not Scheduled" {
			seasons[k].Submitted, err = time.Parse("02/01/2006 15:04", season.SubmittedRaw)
			if err != nil {
				return
			}
		}
	}
	return
}
