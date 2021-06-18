package myradio_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	myradio "github.com/UniversityRadioYork/myradio-go"
	"github.com/UniversityRadioYork/myradio-go/api"
)

func TestMockSession(t *testing.T) {
	session := myradio.MockSession(func(r *api.Request) []byte {
		endpoint := r.Endpoint
		if strings.HasPrefix(endpoint, "/show") {
			return []byte(getSearchMetaJson)
		}
		if strings.HasPrefix(endpoint, "/team") {
			return []byte(positionJSON)
		}
		return []byte("")
	})

	expectedShow := []myradio.ShowMeta{{
		ShowID:        8675309,
		Title:         "Jenny I've Got Your Number",
		CreditsString: "Tommy Tutone",
		Credits: []myradio.Credit{
			{
				Type:     1,
				MemberID: 666,
				User: myradio.User{
					MemberID:     666,
					Fname:        "Tommy",
					Sname:        "Tutone",
					Email:        "tt500@example.com",
					Receiveemail: false,
					Photo:        "/media/image_meta/MyRadioImageMetadata/1.jpeg",
					Bio:          "generic bio",
				},
			},
		},
		Description: "Tommy Tutone's got your number, and he's gotta make you his.",
		ShowTypeID:  1,
		Season: myradio.Link{
			Display: "season display",
			Value:   "https://myradio.example.com/seasons/512",
			Title:   "Seasons",
			URL:     "https://myradio.example.com/seasons/512",
		},
		EditLink: myradio.Link{
			Display: "edit display",
			Value:   "https://myradio.example.com/edit/8675309",
			Title:   "Edit",
			URL:     "https://myradio.example.com/edit/8675309",
		},
		ApplyLink: myradio.Link{
			Display: "apply display",
			Value:   "https://myradio.example.com/apply/8675309",
			Title:   "Apply",
			URL:     "https://myradio.example.com/apply/8675309",
		},
		MicroSiteLink: myradio.Link{
			Display: "microsite display",
			Value:   "https://myradio.example.com/microsites/8675309",
			Title:   "Microsites",
			URL:     "https://myradio.example.com/microsites/8675309",
		},
		Photo: "https://myradio.example.com/photos/shows/8675309",
	}}

	expectedTeams := []myradio.Officer{
		{
			User: myradio.User{
				MemberID:     10,
				Fname:        "John",
				Sname:        "Smith",
				Email:        "john.smith@example.org.uk",
				Receiveemail: true,
				Photo:        "/media/image_meta/MyRadioImageMetadata/1.jpeg",
				Bio:          "generic bio",
			},
			From:            time.Unix(1479081600, 0),
			FromRaw:         1479081600,
			MemberOfficerID: 1,
			Position: myradio.OfficerPosition{
				OfficerID: 2,
				Name:      "Station Manager",
				Alias:     "station.manager",
				Team: myradio.Team{
					TeamID:      1,
					Name:        "Station Management",
					Alias:       "management",
					Ordering:    10,
					Description: "",
					Status:      "c",
				},
				Ordering:    2,
				Description: "",
				Status:      "c",
				Type:        "a",
			},
		},
	}

	showMeta, err := session.GetSearchMeta("tutone")
	if err != nil {
		t.Error(err)
	}

	heads, err := session.GetTeamHeadPositions(1, nil)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(showMeta, expectedShow) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expectedShow, showMeta)
	}

	if !reflect.DeepEqual(heads, expectedTeams) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expectedTeams, heads)
	}
}
