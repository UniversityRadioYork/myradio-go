package myradio

import (
	"errors"
	"time"

	"github.com/UniversityRadioYork/myradio-go/api"
)

// User represents a MyRadio user.
type User struct {
	MemberID     int
	Fname, Sname string
	Email        string `json:"public_email"`
	Receiveemail bool   `json:"receive_email"`
	//@TODO: fix the api and make it return a photo object
	Photo string
	Bio   string
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

// College represents a college.
type College struct {
	CollegeId   int    `json:"value,string"`
	CollegeName string `json:"text"`
}

// GetUser retrieves the User with the given ID.
// This consumes one API request.
func (s *Session) GetUser(id int) (user *User, err error) {
	rq := api.NewRequestf("/user/%d", id)
	rq.Mixins = []string{"personal_data"}
	err = s.do(rq).Into(&user)
	return
}

// GetUserBio retrieves the biography of the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserBio(id int) (bio string, err error) {
	rs := s.getf("/user/%d/bio/", id)
	if rs.IsEmpty() {
		err = errors.New("No bio set")
		return
	}
	err = rs.Into(&bio)
	return
}

// GetUserName retrieves the name of the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserName(id int) (name string, err error) {
	err = s.getf("/user/%d/name/", id).Into(&name)
	return
}

// GetUserProfilePhoto retrieves the profile photo of the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserProfilePhoto(id int) (profilephoto Photo, err error) {
	rs := s.getf("/user/%d/profilephoto/", id)
	if rs.IsEmpty() {
		err = errors.New("No profile picture set")
		return
	}
	err = rs.Into(&profilephoto)
	if err != nil {
		return
	}
	profilephoto.DateAdded, err = time.Parse("02/01/2006 15:04", profilephoto.DateAddedRaw)
	return
}

// GetUserOfficerships retrieves all officerships held by the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserOfficerships(id int) (officerships []Officership, err error) {
	err = s.getf("/user/%d/officerships/", id).Into(&officerships)
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
	err = s.getf("/user/%d/shows/", id).Into(&shows)
	return
}

// GetUserAliases retrieves all aliases associated with the user with the given ID.
// This consumes one API request.
func (s *Session) GetUserAliases() ([]UserAlias, error) {
	raw := [][]string{}
	err := s.get("/user/allaliases/").Into(&raw)
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

// CreateOrActivateUser creates oir activates a new myradio user with the given parameters
// This consumes one API request.
func (s *Session) CreateOrActivateUser(formParams map[string][]string) (user *User, err error) {
	rs := s.post("/user/createoractivate", formParams)
	err = rs.Into(&user)
	return
}

// GetColleges retrieves a list of all current colleges
// This consumes one API request.
func (s *Session) GetColleges() (colleges []College, err error) {
	err = s.get("/user/colleges").Into(&colleges)
	return
}
