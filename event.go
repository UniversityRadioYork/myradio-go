package myradio

import (
	"strconv"
	"time"

	"github.com/UniversityRadioYork/myradio-go/api"
)

type Event struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	HTMLDescription string `json:"description_html"`
	StartTime       time.Time
	StartTimeRaw    string `json:"start"`
	EndTime         time.Time
	EndTimeRaw      string `json:"end"`
	Host            User   `json:"host"`
}

func (e *Event) populateEventTimes() (err error) {
	e.StartTime, err = time.Parse(time.RFC3339, e.StartTimeRaw)
	if err != nil {
		return
	}

	e.EndTime, err = time.Parse(time.RFC3339, e.EndTimeRaw)
	return
}

func (s *Session) GetEventsNext(n int) ([]Event, error) {
	req := api.NewRequest("/event/next")
	req.Params["n"] = []string{strconv.Itoa(n)}
	var events []Event
	err := s.do(req).Into(&events)
	if err != nil {
		return nil, err
	}

	for k := range events {
		err = events[k].populateEventTimes()
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func (s *Session) GetEventsInRange(start time.Time, end time.Time) ([]Event, error) {
	req := api.NewRequest("/event/inrange")
	req.Params["start"] = []string{start.String()}
	req.Params["end"] = []string{end.String()}
	var events []Event
	err := s.do(req).Into(&events)
	if err != nil {
		return nil, err
	}

	for k := range events {
		err = events[k].populateEventTimes()
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}
