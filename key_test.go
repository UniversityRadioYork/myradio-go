package myradio

import (
	"testing"
)

func TestGetApiKeyFromFile(t *testing.T) {

	var tests = []struct{Path, Expected string}{
		{"testdata/.myradio.key", "THIS-IS-A-TEST-KEY-THAT-WILL-NOT-WORK"},
		{"testdata/.linebreaks.key", "THIS-KEY-HAS-SOME-LINE-BREAKS"},
		{"testdata/.shouldntexist.key", ""},
	}

	for _, test := range tests {
		k := getApiKeyFromFile(test.Path)
		if k != test.Expected {
			t.Fail()
		}
	}

}
