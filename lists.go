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
	var body bytes.Buffer
	var ok *bool
	body.WriteString(fmt.Sprintf("userid=%d", UserID))
	err = s.putf("/list/%d/optin", body, ListID).Into(&ok)
	if err != nil {
		return
	}
	if *ok != true {
		err = errors.New("API responded with false")
	}
	return
}
