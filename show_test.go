package myradio_test

import (
	"reflect"
	"testing"

	myradio "github.com/UniversityRadioYork/myradio-go"
)

const getSearchMetaJson = `
[{
	"show_id": 8675309,
	"title": "Jenny I've Got Your Number",
	"credits_string": "Tommy Tutone",
	"credits": [
		{
			"type": 1,
			"memberid": 666,
			"user": {
				"memberid": 666,
				"fname": "Tommy",
				"sname": "Tutone",
				"public_email": "tt500@example.com",
				"receive_email": false,
				"photo": "/media/image_meta/MyRadioImageMetadata/1.jpeg",
          		"bio": "generic bio"
			}
		}
	],
	"description": "Tommy Tutone's got your number, and he's gotta make you his.",
	"show_type_id": 1,
	"seasons": {
		"display": "season display",
		"value": "https://myradio.example.com/seasons/512",
		"title": "Seasons",
		"url": "https://myradio.example.com/seasons/512"
	},
	"editlink": {
		"display": "edit display",
		"value": "https://myradio.example.com/edit/8675309",
		"title": "Edit",
		"url": "https://myradio.example.com/edit/8675309"
	},
	"applylink": {
		"display": "apply display",
		"value": "https://myradio.example.com/apply/8675309",
		"title": "Apply",
		"url": "https://myradio.example.com/apply/8675309"
	},
	"micrositelink": {
		"display": "microsite display",
		"value": "https://myradio.example.com/microsites/8675309",
		"title": "Microsites",
		"url": "https://myradio.example.com/microsites/8675309"
	},
	"photo": "https://myradio.example.com/photos/shows/8675309"
}]`

// TestGetSearchMetaUnmarshal tests the unmarshalling logic of GetSearchMeta.
// It does not test the API endpoint.
func TestGetSearchMetaUnmarshal(t *testing.T) {
	expected := []myradio.ShowMeta{{
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

	session, err := myradio.MockSession([]byte(getSearchMetaJson))
	if err != nil {
		t.Error(err)
	}

	showMeta, err := session.GetSearchMeta("tutone")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(showMeta, expected) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, showMeta)
	}
}
