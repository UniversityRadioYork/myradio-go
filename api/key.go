package api

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// The contents of this file are heavily based on
// https://github.com/UniversityRadioYork/urydb-go/blob/master/urydb.go

var (
	// ErrNoMYRADIOKEYFILE is the error thrown when MYRADIOKEYFILE is not present in
	// the environment.
	ErrNoMYRADIOKEYFILE = errors.New("MYRADIOKEYFILE not in environment")
	// ErrNoKeyFile is the error thrown when there
	// is no myradio.key file.
	ErrNoKeyFile = errors.New("couldn't find any API key file")
)

// KeyFiles is the list of possible places to
// search for a urydb file.
var KeyFiles = []string{
	".myradio.key",
	"${HOME}/.myradio.key",
	"/etc/myradio.key",
	"/usr/local/etc/myradio.key",
}

// GetAPIKey tries to get an API key from all possible sources.
// This tries the following paths for a file containing one line (the API key):
//   1) Whichever path is set in the environment variable `MYRADIOKEYFILE`;
//   2) `.myradio.key`, in the current directory;
//   3) `.myradio.key`, in the user's home directory;
//   4) `/etc/myradio.key`;
//   5) `/usr/local/etc/myradio.key`.`
func GetAPIKey() (apikey string, err error) {
	apikey, err = getAPIKeyEnv()
	if err != nil {
		apikey, err = getAPIKeyFile()
	}
	return
}

// getAPIKey tries to get an API key from the environment.
func getAPIKeyEnv() (apikey string, err error) {
	apikey, err = os.Getenv("MYRADIOKEYFILE"), nil
	if apikey == "" {
		err = ErrNoMYRADIOKEYFILE
	}
	return
}

// getAPIKeyFile tries to get an API key from a known file.
func getAPIKeyFile() (apikey string, err error) {
	for _, rawPath := range KeyFiles {
		apikey = getAPIKeyFromFile(rawPath)
		if apikey != "" {
			return
		}
	}
	if apikey == "" {
		err = ErrNoKeyFile
	}
	return
}

// getAPIKeyFromFile tries to get an apikey from a file.
// Returns an empty string if it fails
func getAPIKeyFromFile(path string) string {
	path = os.ExpandEnv(path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	s := string(b)
	return strings.TrimSpace(s)
}

