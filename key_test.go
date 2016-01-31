package myradio

import (
	"testing"
)

func TestGetApiKeyFromFileWorks(t *testing.T){

	k := getApiKeyFromFile("testdata/.myradio.key")

	if k != "THIS-IS-A-TEST-KEY-THAT-WILL-NOT-WORK" {
		t.Fail()
	}

}

func TestGetApiKeyFromFileWorksWithLineBreaks(t *testing.T) {

	k := getApiKeyFromFile("testdata/.linebreaks.key")

	if k != "THIS-KEY-HAS-SOME-LINE-BREAKS" {
		t.Fail()
	}

}

func TestGetApiKeyFromFileReturnsEmptyString(t *testing.T) {

	k := getApiKeyFromFile("testdata/.shouldntexist.key")

	if k != "" {
		t.Fail()
	}

}