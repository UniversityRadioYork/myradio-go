package myradio

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type CurrentAndNext struct {
	Next    Show `json:"next"`
	Current Show `json:"current"`
}

type Show struct {
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Photo        string `json:"photo"`
	StartTimeRaw int64  `json:"start_time"`
	StartTime    time.Time
	EndTimeRaw   int64 `json:"end_time"`
	EndTime      time.Time
	Presenters   string `json:"presenters,omitempty"`
	Url          string `json:"url,omitempty"`
	Id           uint64 `json:"id,omitempty"`
}

type Timeslot struct {
	Season
	TimeslotID     uint64   `json:"timeslot_id"`
	TimeslotNum    int      `json:"timeslot_num"`
	Tags           []string `json:"tags"`
	Time           time.Time
	TimeRaw        int64 `json:"time"`
	StartTime      time.Time
	StartTimeRaw   string `json:"start_time"`
	Duration       time.Duration
	DurationRaw    string `json:"duration"`
	MixcloudStatus string `json:"mixcloud_status"`
}

type TracklistItem struct {
	Track
	Album        Album `json:"album"`
	EditLink     Link  `json:"editlink"`
	DeleteLink   Link  `json:"deletelink"`
	Time         time.Time
	TimeRaw      int64 `json:"time"`
	StartTime    time.Time
	StartTimeRaw string `json:"starttime"`
	AudioLogID   uint   `json:"audiologid"`
}

type Schedule [][]Timeslot

type Day []Timeslot

func (s *Session) GetCurrentAndNext() (*CurrentAndNext, error) {
	data, err := s.apiRequest("/timeslot/currentandnext", []string{})
	if err != nil {
		return nil, err
	}
	var currentAndNext CurrentAndNext
	err = json.Unmarshal(*data, &currentAndNext)
	if err != nil {
		return nil, err
	}
	currentAndNext.Current.StartTime = time.Unix(currentAndNext.Current.StartTimeRaw, 0)
	currentAndNext.Current.EndTime = time.Unix(currentAndNext.Current.EndTimeRaw, 0)
	currentAndNext.Next.StartTime = time.Unix(currentAndNext.Next.StartTimeRaw, 0)
	currentAndNext.Next.EndTime = time.Unix(currentAndNext.Next.EndTimeRaw, 0)
	return &currentAndNext, nil
}

// GetWeekSchedule gets the weekly schedule for ISO 8601 week week of year year.
// It returns the result as an map from ISO 8601 weekdays to timeslot slices.
// Thus, 1 maps to Monday's timeslots; 2 to Tuesday; and so on.
// Each slice progresses chronologically from start of URY day to finish of URY day.
func (s *Session) GetWeekSchedule(year, week int) (Schedule, error) {
	// TODO(CaptainHayashi): proper errors
	if year < 0 {
		return nil, fmt.Errorf("year %d is too low", year)
	}
	if week < 1 || 53 < week {
		return nil, fmt.Errorf("week %d is not within the ISO range 1..53", week)
	}

	data, err := s.apiRequestWithParams(fmt.Sprintf("/timeslot/weekschedule/%d", week), []string{}, map[string][]string{"year": {strconv.Itoa(year)}})
	if err != nil {
		return nil, err
	}

	// The timeslots come to us with string keys labelled with the weekday.
	// These timeslots start from "1" (Monday) and go up to "7" (Sunday).
	// Note that this is different from Go's view of the week (0 = Sunday, 1 = Monday)!
	stringyTimeslots := make(map[string][]Timeslot)
	err = json.Unmarshal(*data, &stringyTimeslots)
	if err != nil {
		return nil, err
	}

	// Now convert the string keys into proper indices.
	timeslots := make(Schedule, len(stringyTimeslots))
	for sday, ts := range stringyTimeslots {
		day, err := strconv.Atoi(sday)
		if err != nil {
			return nil, err
		}
		timeslots[day - 1] = ts
	}

	return timeslots, nil
}

func (s *Session) GetTimeslot(id int) (timeslot Timeslot, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/timeslot/%d", id), []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &timeslot)
	timeslot.Time = time.Unix(timeslot.TimeRaw, 0)
	timeslot.FirstTime, err = time.Parse("02/01/2006 15:04", timeslot.FirstTimeRaw)
	if err != nil {
		return
	}
	timeslot.Submitted, err = time.Parse("02/01/2006 15:04", timeslot.SubmittedRaw)
	if err != nil {
		return
	}
	timeslot.StartTime, err = time.Parse("02/01/2006 15:04", timeslot.StartTimeRaw)
	if err != nil {
		return
	}
	timeslot.Duration, err = parseDuration("15:04:05", timeslot.DurationRaw)
	if err != nil {
		return
	}
	return
}

func (s *Session) GetTrackListForTimeslot(id int) (tracklist []TracklistItem, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/tracklistItem/tracklistfortimeslot/%d", id), []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &tracklist)
	for k, v := range tracklist {
		tracklist[k].Time = time.Unix(tracklist[k].TimeRaw, 0)
		tracklist[k].StartTime, err = time.Parse("02/01/2006 15:04:05", v.StartTimeRaw)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (t *Timeslot) EndTime() time.Time {
	return t.StartTime.Add(t.Duration)
}
