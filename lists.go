package myradio

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/UniversityRadioYork/myradio-go/api"
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

// OptIn subscribes the given user to the given list
// This consumes one API request.
func (s *Session) OptIn(UserID int, ListID int) (err error) {
	body := bytes.NewBufferString(fmt.Sprintf("userid=%d", UserID))
	var ok *bool
	if err = s.putf("/list/%d/optin", *body, ListID).Into(&ok); err != nil {
		return
	}
	if !*ok {
		err = errors.New("API responded with false")
	}
	return
}
