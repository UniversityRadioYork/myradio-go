package myradio

import (
	"encoding/json"
	"fmt"
	"time"
)

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
		timeslots[k].FirstTime, err = time.Parse("02/01/2006 15:04", v.FirstTimeRaw)
		if err != nil {
			return
		}
		timeslots[k].Submitted, err = time.Parse("02/01/2006 15:04", v.SubmittedRaw)
		if err != nil {
			return
		}
		timeslots[k].StartTime, err = time.Parse("02/01/2006 15:04", v.StartTimeRaw)
		if err != nil {
			return
		}
		timeslots[k].Duration, err = parseDuration("15:04:05", v.DurationRaw)
		if err != nil {
			return
		}
	}
	return
}
