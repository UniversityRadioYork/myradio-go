# MyRadio.go
[![Build Status](https://travis-ci.org/UniversityRadioYork/myradio-go.svg?branch=master)](https://travis-ci.org/UniversityRadioYork/myradio-go)
[![Coverage Status](https://coveralls.io/repos/github/UniversityRadioYork/myradio-go/badge.svg?branch=master)](https://coveralls.io/github/UniversityRadioYork/myradio-go?branch=master)

A go wrapper for the MyRadio API. Incomplete. Quick and dirty. Originally for the go rewrite of aliasgen.

## Usage

```go
import "github.com/UniversityRadioYork/myradio-go"

...

session, _ := myradio.NewSession("your_api_key")

lists := session.GetAllLists()
```


## Testing

```bash
$ go test
```