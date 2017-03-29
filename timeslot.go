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
	EndTimeRaw   string `json:"end_time"` // Sometimes "The End of Time"
	EndTime      time.Time
	Presenters   string `json:"presenters,omitempty"`
	Url          string `json:"url,omitempty"`
	Id           uint64 `json:"id,omitempty"`
}

// populateShowTimes sets the times for the given Show given their raw values.
func (s *Show) populateShowTimes() error {
	s.StartTime = time.Unix(s.StartTimeRaw, 0)

	timeint, err := strconv.ParseInt(s.EndTimeRaw, 10, 64)
	if err != nil {
		return err
	}

	s.EndTime = time.Unix(timeint, 0)
	return nil
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

// populateTimeslotTimes sets the times for the given Timeslot given their raw values.
func (t *Timeslot) populateTimeslotTimes() (err error) {
	// Remember: a Timeslot is a supertype of Season.
	if err = t.populateSeasonTimes(); err != nil {
		return
	}

	t.Time = time.Unix(t.TimeRaw, 0)

	t.StartTime, err = parseShortTime(t.StartTimeRaw)
	if err != nil {
		return
	}

	t.Duration, err = parseDuration(t.DurationRaw)
	return
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
	data, aerr := s.apiRequest("/timeslot/currentandnext", []string{})
	if aerr != nil {
		return nil, aerr
	}

	var currentAndNext CurrentAndNext
	if err := json.Unmarshal(*data, &currentAndNext); err != nil {
		return nil, err
	}
	
	if err := currentAndNext.Current.populateShowTimes(); err != nil {
		return nil, err
	}

	if err := currentAndNext.Next.populateShowTimes(); err != nil {
		return nil, err
	}

	return &currentAndNext, nil
}

// GetWeekSchedule gets the weekly schedule for ISO 8601 week week of year year.
// If such a schedule exists, it returns the result as an map from ISO 8601 weekdays to timeslot slices.
// Thus, 1 maps to Monday's timeslots; 2 to Tuesday; and so on.
// Each slice progresses chronologically from start of URY day to finish of URY day.
// If no such schedule exists, it returns a map of empty slices.
// If an error occurred, this is returned in error, and the timeslot map is undefined.
// This consumes one API request.
func (s *Session) GetWeekSchedule(year, week int) (map[int][]Timeslot, error) {
	// TODO(CaptainHayashi): proper errors
	if year < 0 {
		return nil, fmt.Errorf("year %d is too low", year)
	}
	if week < 1 || 53 < week {
		return nil, fmt.Errorf("week %d is not within the ISO range 1..53", week)
	}

	data, aerr := s.apiRequestWithParams(fmt.Sprintf("/timeslot/weekschedule/%d", week), []string{}, map[string][]string{"year": {strconv.Itoa(year)}})
	if aerr != nil {
		return nil, aerr
	}

	// MyRadio responds in a different way when the schedule is empty, so we need to catch that.
	// See https://github.com/UniversityRadioYork/MyRadio/issues/665 for details.
	if isEmptySchedule(data) {
		return map[int][]Timeslot{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
			6: {},
			7: {},
		}, nil
	}

	// The timeslots come to us with string keys labelled with the weekday.
	// These timeslots start from "1" (Monday) and go up to "7" (Sunday).
	// Note that this is different from Go's view of the week (0 = Sunday, 1 = Monday)!
	stringyTimeslots := make(map[string][]Timeslot)
	if jerr := json.Unmarshal(*data, &stringyTimeslots); jerr != nil {
		return nil, jerr
	}

	return destringTimeslots(stringyTimeslots)
}

// isEmptySchedule tries to work out, from MyRadio schedule JSON, whether the schedule is empty.
func isEmptySchedule(data json.Marshaler) bool {
	bs, err := data.MarshalJSON()
	if err != nil {
		// The logic later on in GetWeekSchedule should hit this same error, so handle it there.
		return false
	}

	if len(bs) != 2 {
		return false
	}

	// Due to a quirk in MyRadio (well, PHP), the empty schedule can be returned as the empty array '[]',
	// instead of the empty object '{}'.
	if bs[0] == '[' && bs[1] == ']' {
		return true
	}

	if bs[0] == '{' && bs[1] == '}' {
		return true
	}

	return false
}

// destringTimeslots converts a week schedule from string indices to integer indices.
// It takes a map from strings "1"--"7" to day schedules, and returns a map from integers 1--7 to day schedules.
// It returns an error if any of the string indices cannot be converted.
func destringTimeslots(stringyTimeslots map[string][]Timeslot) (map[int][]Timeslot, error) {
	timeslots := make(map[int][]Timeslot)
	for sday, ts := range stringyTimeslots {
		day, derr := strconv.Atoi(sday)
		if derr != nil {
			return nil, derr
		}
		for i := range ts {
			if terr := ts[i].populateTimeslotTimes(); terr != nil {
				return nil, terr
			}
		}
		timeslots[day] = ts
	}

	return timeslots, nil
}

// GetTimeslot retrieves the timeslot with the given ID.
// This consumes one API request.
func (s *Session) GetTimeslot(id int) (timeslot Timeslot, err error) {
	var data *json.RawMessage
	if data, err = s.apiRequest(fmt.Sprintf("/timeslot/%d", id), []string{}); err != nil {
		return
	}

	if err = json.Unmarshal(*data, &timeslot); err != nil {
		return
	}

	err = timeslot.populateTimeslotTimes()
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
