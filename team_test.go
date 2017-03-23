package myradio_test

import (
	"reflect"
	"testing"

	myradio "github.com/UniversityRadioYork/myradio-go"
)

const assistantHeadPositionJSON = `
[{
	"User": {
		"memberid": 10,
		"fname": "John",
		"sname": "Smith",
		"sex": "m",
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

// TestGetSearchMetaUnmarshal tests the unmarshalling logic of GetSearchMeta.
// It does not test the API endpoint.
func TestGetTeamHeadPositions(t *testing.T) {
	expected := []myradio.Officer{
		{
			User: myradio.User{
				Memberid:     10,
				Fname:        "John",
				Sname:        "Smith",
				Sex:          "m",
				Email:        "john.smith@example.org.uk",
				Receiveemail: true,
			},
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

	session, err := myradio.MockSession([]byte(assistantHeadPositionJSON))
	if err != nil {
		t.Error(err)
	}

	heads, err := session.GetTeamHeadPositions(1, nil)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(heads, expected) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, heads)
	}
}
