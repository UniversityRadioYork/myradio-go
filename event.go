package myradio

import (
	"github.com/UniversityRadioYork/myradio-go/api"
	"strconv"
)

type Event struct {
	ID          int `json:"id"`
	Title       string
	Description string
	Start       string
	End         string
	Host        User
}

func (s *Session) GetEventsNext(n int) ([]Event, error) {
	req := api.NewRequest("/event/next")
	req.Params["n"] = []string{strconv.Itoa(n)}
	var events []Event
	err := s.do(req).Into(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Session) GetEventsInRange(start string, end string) ([]Event, error) {
	req := api.NewRequest("/event/inrange")
	req.Params["start"] = []string{start}
	req.Params["end"] = []string{end}
	var events []Event
	err := s.do(req).Into(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
