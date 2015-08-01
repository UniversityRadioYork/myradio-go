package myradio

import (
	"encoding/json"
	"fmt"
)

type List struct {
	Listid     int
	Name       string
	Address    string
	Recipients int `json:"recipient_count"`
}

func (s *Session) GetAllLists() ([]List, error) {
	data, err := s.apiRequest("/list/alllists", nil)
	if err != nil {
		return nil, err
	}
	var lists []List
	err = json.Unmarshal(*data, &lists)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (s *Session) GetMembers(l *List) ([]Member, error) {
	data, err := s.apiRequest(fmt.Sprintf("/list/%d/members", l.Listid), []string{"personal_data"})
	if err != nil {
		return nil, err
	}
	var members []Member
	err = json.Unmarshal(*data, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}
