package myradio

import (
	"bufio"
	"errors"
	"os"
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

func getApiKey() (apikey string, err error) {
	apikey, err = getApiKeyEnv()
	if err != nil {
		apikey, err = getApiKeyFile()
	}
	return
}

func getApiKeyEnv() (apikey string, err error) {
	apikey, err = os.Getenv("MYRADIOKEYFILE"), nil
	if apikey == "" {
		err = ErrNoMYRADIOKEYFILE
	}
	return
}

func getApiKeyFile() (apikey string, err error) {
	for _, rawPath := range KeyFiles {
		path := os.ExpandEnv(rawPath)
		file, ferr := os.Open(path)
		if ferr != nil {
			apikey = ""
			continue
		}

		bufrd := bufio.NewReader(file)
		apikey, ferr = bufrd.ReadString('\n')

		if ferr != nil {
			apikey = ""
			continue
		}

		return
	}

	if apikey == "" {
		err = ErrNoKeyFile
	}
	return
}

// NewSessionFromKeyFile tries to open a Session with the key from an API key file.
//
// This tries the following paths for a file containing one line (the API key):
//   1) Whichever path is set in the environment variable `MYRADIOKEYFILE`;
//   2) `.myradio.key`, in the current directory;
//   3) `.myradio.key`, in the user's home directory;
//   4) `/etc/myradio.key`;
//   5) `/usr/local/etc/myradio.key`.
func NewSessionFromKeyFile() (*Session, error) {
	apikey, err := getApiKey()
	if err != nil {
		return nil, err
	}

	return NewSession(apikey)
}
