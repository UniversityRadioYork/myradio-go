package myradio

import (
	"encoding/json"
	"fmt"
)

// List represents a mailing list.
type List struct {
	Listid     int
	Name       string
	Address    string
	Recipients int `json:"recipient_count"`
}

// GetAllLists retrieves all mailing lists in the MyRadio system.
// This consumes one API request.
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

// GetUsers retrieves all users subscribed to a given mailing list.
// This consumes one API request.
func (s *Session) GetUsers(l *List) ([]User, error) {
	data, err := s.apiRequest(fmt.Sprintf("/list/%d/members", l.Listid), []string{"personal_data"})
	if err != nil {
		return nil, err
	}
	var user []User
	err = json.Unmarshal(*data, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
