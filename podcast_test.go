package myradio_test

import (
	"reflect"
	"testing"
	"time"

	myradio "github.com/UniversityRadioYork/myradio-go"
)

const getPodcastSearchMetaJson = `
[{
	"podcast_id": 1234,
	"title": "test podcast",
	"description": "an amazing podcast",
	"status": "Published",
	"time": 1234512345,
	"uri": "https://myradio.example.com/media/podcasts/podcast1234.mp3",
	"editlink": {
		"display": "icon",
		"value": "pencil",
		"title": "Edit Podcast",
		"url": "https://myradio.example.com/podcast/editPodcast/1234"
	},
	"micrositelink": {
		"display": "icon",
		"value": "link",
		"title": "Microsites",
		"url": "https://myradio.example.com/microsites/1234"
	},
	"photo": "https://myradio.example.com/media/image_meta/MyRadioImageMetadata/podcast1234.jpeg"
}]`

// TestGetSearchMetaUnmarshal tests the unmarshalling logic of GetSearchMeta.
// It does not test the API endpoint.
func TestGetPodcastSearchMetaUnmarshal(t *testing.T) {
	expected := []myradio.Podcast{{
		PodcastID:   1234,
		Title:       "test podcast",
		Description: "an amazing podcast",
		Status:      "Published",
		Time:        myradio.Time{time.Unix(1234512345, 0)},
		File:        "https://myradio.example.com/media/podcasts/podcast1234.mp3",
		EditLink: myradio.Link{
			Display: "icon",
			Value:   "pencil",
			Title:   "Edit Podcast",
			URL:     "https://myradio.example.com/podcast/editPodcast/1234",
		},
		MicrositeLink: myradio.Link{
			Display: "icon",
			Value:   "link",
			Title:   "Microsites",
			URL:     "https://myradio.example.com/microsites/1234",
		},
		Photo: "https://myradio.example.com/media/image_meta/MyRadioImageMetadata/podcast1234.jpeg",
	}}

	session, err := myradio.MockSession([]byte(getPodcastSearchMetaJson))
	if err != nil {
		t.Error(err)
	}

	podcast, err := session.GetPodcastSearchMeta("test")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(podcast, expected) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, podcast)
	}
}
