package myradio

import "github.com/UniversityRadioYork/myradio-go/api"

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
	err = s.get("/list/alllists").Into(&lists)
	return
}

// GetUsers retrieves all users subscribed to a given mailing list.
// This consumes one API request.
func (s *Session) GetUsers(l *List) (users []User, err error) {
	rq := api.NewRequestf("/list/%d/members", l.Listid)
	rq.Mixins = []string{"personal_data"}
	err = s.do(rq).Into(&users)
	return
}
