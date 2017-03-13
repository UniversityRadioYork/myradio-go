package myradio

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// CurrentAndNext stores a pair of current and next show.
type CurrentAndNext struct {
	Next    Show `json:"next"`
	Current Show `json:"current"`
}

// Show contains a summary of information about a URY schedule timeslot.
type Show struct {
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Photo        string `json:"photo"`
	StartTimeRaw int64  `json:"start_time"`
	StartTime    time.Time
	EndTimeRaw   string `json:"end_time"`
	EndTime      time.Time
	Presenters   string `json:"presenters,omitempty"`
	Url          string `json:"url,omitempty"`
	Id           uint64 `json:"id,omitempty"`
}

// Timeslot contains information about a single timeslot in the URY schedule.
// A timeslot is a single slice of time on the schedule, typically one hour long.
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

// TracklistItem represents a single item in a show tracklist.
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

// GetCurrentAndNext gets the current and next shows at the time of the call.
// This consumes one API request.
func (s *Session) GetCurrentAndNext() (*CurrentAndNext, error) {
	data, err := s.apiRequest("/timeslot/currentandnext", []string{})
	if err != nil {
		return nil, err
	}
	var currentAndNext CurrentAndNext
	var currentEndTime int64
	var nextEndTime int64
	err = json.Unmarshal(*data, &currentAndNext)
	if err != nil {
		return nil, err
	}
	if currentAndNext.Current.EndTimeRaw != "The End of Time" && (currentAndNext.Current.EndTimeRaw != "") {
		currentEndTime, err = strconv.ParseInt(currentAndNext.Current.EndTimeRaw, 10, 64)
		if err != nil {
			return nil, err
		}
		currentAndNext.Current.EndTime = time.Unix(currentEndTime, 0)
	}
	if (currentAndNext.Next.EndTimeRaw != "The End of Time") && (currentAndNext.Next.EndTimeRaw != "") {
		nextEndTime, err = strconv.ParseInt(currentAndNext.Next.EndTimeRaw, 10, 64)
		if err != nil {
			return nil, err
		}
		currentAndNext.Next.EndTime = time.Unix(nextEndTime, 0)
	}

	currentAndNext.Current.StartTime = time.Unix(currentAndNext.Current.StartTimeRaw, 0)
	currentAndNext.Next.StartTime = time.Unix(currentAndNext.Next.StartTimeRaw, 0)
	return &currentAndNext, nil
}

// populateTimes fills in the cooked times in a MyRadio timeslot from their raw equivalents.
func populateTimes(timeslot *Timeslot) error {
	var err error = nil
	timeslot.Time = time.Unix(timeslot.TimeRaw, 0)

	// MyRadio returns local timestamps, not UTC.
	timeslot.FirstTime, err = time.ParseInLocation("02/01/2006 15:04", timeslot.FirstTimeRaw, time.Local)
	if err != nil {
		return err
	}
	timeslot.Submitted, err = time.ParseInLocation("02/01/2006 15:04", timeslot.SubmittedRaw, time.Local)
	if err != nil {
		return err
	}
	timeslot.StartTime, err = time.ParseInLocation("02/01/2006 15:04", timeslot.StartTimeRaw, time.Local)
	if err != nil {
		return err
	}
	timeslot.Duration, err = parseDuration(timeslot.DurationRaw)
	if err != nil {
		return err
	}

	return nil
}

// GetWeekSchedule gets the weekly schedule for ISO 8601 week week of year year.
// It returns the result as an map from ISO 8601 weekdays to timeslot slices.
// Thus, 1 maps to Monday's timeslots; 2 to Tuesday; and so on.
// Each slice progresses chronologically from start of URY day to finish of URY day.
// This consumes one API request.
func (s *Session) GetWeekSchedule(year, week int) (map[int][]Timeslot, error) {
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
	timeslots := make(map[int][]Timeslot)
	for sday, ts := range stringyTimeslots {
		day, err := strconv.Atoi(sday)
		if err != nil {
			return nil, err
		}
		for i := range ts {
			err = populateTimes(&ts[i])
			if err != nil {
				return nil, err
			}
		}
		timeslots[day] = ts
	}

	return timeslots, nil
}

// GetTimeslot retrieves the timeslot with the given ID.
// This consumes one API request.
func (s *Session) GetTimeslot(id int) (timeslot Timeslot, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/timeslot/%d", id), []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &timeslot)
	if err != nil {
		return
	}
	err = populateTimes(&timeslot)
	if err != nil {
		return
	}
	return
}

// GetTrackListForTimeslot retrieves the tracklist for the timeslot with the given ID.
// This consumes one API request.
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
