# MyRadio.go

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