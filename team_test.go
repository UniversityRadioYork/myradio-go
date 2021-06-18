package myradio_test

import (
	"reflect"
	"testing"

	"time"

	myradio "github.com/UniversityRadioYork/myradio-go"
)

const positionJSON = `
[{
	"User": {
		"memberid": 10,
		"fname": "John",
		"sname": "Smith",
		"public_email": "john.smith@example.org.uk",
		"url": "//example.org.uk/myradio/Profile/view/?memberid=10",
		"receive_email": true,
		"photo": "/media/image_meta/MyRadioImageMetadata/1.jpeg",
		"bio": "generic bio"
	},
	"from": 1479081600,
	"memberofficerid": 1,
	"position": {
		"officerid": 2,
		"name": "Station Manager",
		"alias": "station.manager",
		"team": {
			"teamid": 1,
			"name": "Station Management",
			"alias": "management",
			"ordering": 10,
			"description": "",
			"status": "c"
		},
		"ordering": 2,
		"description": "",
		"status": "c",
		"type": "a"
	}
}]`

// TestGetTeamHeadPositions tests the getter for head positions of a team.
// It does not test the API endpoint.
func TestGetTeamHeadPositions(t *testing.T) {
	expected := []myradio.Officer{
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

	session := myradio.MockSession([]byte(positionJSON))

	heads, err := session.GetTeamHeadPositions(1, nil)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(heads, expected) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, heads)
	}
}

// TestGetTeamAssistantHeadPositions tests the getter for assistant head positions of a team.
// It does not test the API endpoint.
func TestGetTeamAssistantHeadPositions(t *testing.T) {
	expected := []myradio.Officer{
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

	session := myradio.MockSession([]byte(positionJSON))

	heads, err := session.GetTeamAssistantHeadPositions(1, nil)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(heads, expected) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, heads)
	}
}
