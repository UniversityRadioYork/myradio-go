package myradio

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// User represents a MyRadio user.
type User struct {
	Memberid     int
	Fname, Sname string
	Sex          string
	Email        string `json:"public_email"`
	Receiveemail bool   `json:"receive_email"`
}

// Officership represents an officership a user holds.
type Officership struct {
	OfficerId   uint   `json:"officerid,string"`
	OfficerName string `json:"officer_name"`
	TeamId      uint   `json:"teamid,string"`
	FromDateRaw string `json:"from_date,omitempty"`
	FromDate    time.Time
	TillDateRaw string `json:"till_date,omitempty"`
	TillDate    time.Time
}

// Photo represents a photo of a user.
type Photo struct {
	PhotoId      uint   `json:"photoid"`
	DateAddedRaw string `json:"date_added"`
	DateAdded    time.Time
	Format       string `json:"format"`
	Owner        uint   `json:"owner"`
	Url          string `json:"url"`
}

// UserAlias represents a user alias.
type UserAlias struct {
	Source      string
	Destination string
}

// GetUser retrieves the User with the given ID.
// This consumes one API request.
func (s *Session) GetUser(id int) (user *User, err error) {
	err = s.get(fmt.Sprintf("/user/%d", id)).mixin("personal_data").into(&user)
	return
}

// GetUserBio retrieves the biography of the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserBio(id int) (bio string, err error) {
	data, err := s.get(fmt.Sprintf("/user/%d/bio/", id)).do()
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("No bio set")
		return
	}
	err = json.Unmarshal(*data, &bio)
	return
}

// GetUserName retrieves the name of the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserName(id int) (name string, err error) {
	err = s.get(fmt.Sprintf("/user/%d/name/", id)).into(&name)
	return
}

// GetUserProfilePhoto retrieves the profile photo of the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserProfilePhoto(id int) (profilephoto Photo, err error) {
	data, err := s.get(fmt.Sprintf("/user/%d/profilephoto/", id)).do()
	if err != nil {
		return
	}
	if data == nil {
		err = errors.New("No profile picture set")
		return
	}
	err = json.Unmarshal(*data, &profilephoto)
	if err != nil {
		return
	}
	profilephoto.DateAdded, err = time.Parse("02/01/2006 15:04", profilephoto.DateAddedRaw)
	return
}

// GetUserOfficerships retrieves all officerships held by the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserOfficerships(id int) (officerships []Officership, err error) {
	err = s.get(fmt.Sprintf("/user/%d/officerships/", id)).into(&officerships)
	if err != nil {
		return
	}
	for k, v := range officerships {
		if officerships[k].FromDateRaw != "" {
			officerships[k].FromDate, err = time.Parse("2006-01-02", v.FromDateRaw)
			if err != nil {
				return
			}
		}
		if officerships[k].TillDateRaw != "" {
			officerships[k].TillDate, err = time.Parse("2006-01-02", v.TillDateRaw)
			if err != nil {
				return
			}
		}
	}
	return
}

// GetUserShowCredits retrieves all show credits associated with the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserShowCredits(id int) (shows []ShowMeta, err error) {
	err = s.get(fmt.Sprintf("/user/%d/shows/", id)).into(&shows)
	return
}

// GetUserAliases retrieves all aliases associated with the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserAliases() ([]UserAlias, error) {
	raw := [][]string{}
	err := s.get("/user/allaliases/").into(&raw)
	if err != nil {
		return nil, err
	}
	var aliases = make([]UserAlias, len(raw))
	for k, v := range raw {
		aliases[k].Source = v[0]
		aliases[k].Destination = v[1]
	}
	return aliases, nil
}
