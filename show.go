package myradio

import (
	"net/url"
)

// Credit represents a show credit associating a user with a show.
type Credit struct {
	Type     int  `json:"type"`
	MemberID int  `json:"memberid"`
	User     User `json:"User"`
}

// ShowMeta represents a show in the MyRadio schedule.
// A MyRadio show contains seasons, each containing timeslots.
// @TODO: Refactor this to something better named
type ShowMeta struct {
	ShowID        int      `json:"show_id"`
	Title         string   `json:"title"`
	CreditsString string   `json:"credits_string"`
	Credits       []Credit `json:"credits"`
	Description   string   `json:"description"`
	ShowTypeID    int      `json:"show_type_id"`
	Season        Link     `json:"seasons"`
	EditLink      Link     `json:"editlink"`
	ApplyLink     Link     `json:"applylink"`
	MicroSiteLink Link     `json:"micrositelink"`
	Photo         string   `json:"photo"`
}

// Link represents a MyRadio action link.
type Link struct {
	Display string      `json:"display"`
	Value   interface{} `json:"value"`
	Title   string      `json:"title,omitempty"`
	URL     string      `json:"url"`
}

// GetSearchMeta retrieves all shows whose metadata matches a given search term.
// This consumes one API request.
func (s *Session) GetSearchMeta(term string) (shows []ShowMeta, err error) {
	err = s.getf("/show/searchmeta/%s", url.QueryEscape(term)).Into(&shows)
	return
}

// GetShow retrieves the show with the given ID.
// This consumes one API request.
func (s *Session) GetShow(id int) (show *ShowMeta, err error) {
	err = s.getf("/show/%d", id).Into(&show)
	return
}

// GetSeasons retrieves the seasons of the show with the given ID.
// This consumes one API request.
func (s *Session) GetSeasons(id int) (seasons []Season, err error) {
	if err = s.getf("/show/%d/allseasons", id).Into(&seasons); err != nil {
		return
	}

	for i := range seasons {
		err = seasons[i].populateSeasonTimes()
		if err != nil {
			return
		}
	}

	return
}

// GetCreditsToUsers retrieves a map of credit names to users.
// This consumes two API request.
func (s *Session) GetCreditsToUsers(id int, isTimeslot bool) (creditsToUsers map[string][]User, err error) {

	type creditType struct {
		Type int    `json:"value,string"`
		Name string `json:"text"`
	}

	// First get the credit type to name
	var creditTypes []creditType
	if err = s.get("/scheduler/credittypes").Into(&creditTypes); err != nil {
		return
	}

	// Get the credits of a show
	var credits []Credit
	var requestPath string

	if isTimeslot {
		requestPath = "/timeslot/%d/credits"
	} else {
		requestPath = "/show/%d/credits"
	}

	if err = s.getf(requestPath, id).Into(&credits); err != nil {
		return
	}

	// A map of credit type (int value) to the name of the credit (string)
	var creditToName = make(map[int]string)
	for _, cT := range creditTypes {
		creditToName[cT.Type] = cT.Name
	}

	// A map of credit names (strings) to a list of users
	creditsToUsers = make(map[string][]User)
	for _, credit := range credits {
		var creditName = creditToName[credit.Type]
		creditsToUsers[creditName] = append(creditsToUsers[creditName], credit.User)
	}

	return
}
