package myradio

import "fmt"

// List represents a mailing list.
type List struct {
	Listid     int
	Name       string
	Address    string
	Recipients int `json:"recipient_count"`
}

// GetAllLists retrieves all mailing lists in the MyRadio system.
// This consumes one API request.
func (s *Session) GetAllLists() (lists []List, err error) {
	err = s.apiRequestInto(&lists, "/list/alllists", nil)
	return
}

// GetUsers retrieves all users subscribed to a given mailing list.
// This consumes one API request.
func (s *Session) GetUsers(l *List) (users []User, err error) {
	err = s.apiRequestInto(&users, fmt.Sprintf("/list/%d/members", l.Listid), []string{"personal_data"})
	return
}
