package api

import (
	"os"
	"testing"
)

func TestGetAPIKeyFromFile(t *testing.T) {
	var tests = []struct{ Path, Expected string }{
		{"testdata/.myradio.key", "THIS-IS-A-TEST-KEY-THAT-WILL-NOT-WORK"},
		{"testdata/.linebreaks.key", "THIS-KEY-HAS-SOME-LINE-BREAKS"},
		{"testdata/.shouldntexist.key", ""},
		{"testdata/.hasspaceinit.key", "this has spaces in it"},
	}

	for _, test := range tests {
		k := getAPIKeyFromFile(test.Path)
		if k != test.Expected {
			t.Fatal("expected:", test.Expected, "got:", k)
		}
	}
}

func TestGetAPIKeyEnv(t *testing.T) {
	os.Setenv("MYRADIOKEY", "foobar")
	k, err := getAPIKeyEnv()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if k != "foobar" {
		t.Fatal("expected foobar, got:", k)
	}

	os.Setenv("MYRADIOKEY", "")
	_, nerr := getAPIKeyEnv()
	if nerr != ErrNoMYRADIOKEY {
		t.Fatal("expected error:", ErrNoMYRADIOKEY, "got error:", nerr)
	}
}
