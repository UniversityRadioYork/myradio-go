package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/UniversityRadioYork/myradio-go/api"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-url=URL] ENDPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	baseurl := flag.String("url", `https://ury.org.uk/api/v2`, "The MyRadio API base URL.")
	flag.Parse()

	endpoint := flag.Arg(0)
	if endpoint == "" {
		usage()
		return
	}

	k, err := api.GetAPIKey()
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(*baseurl)
	if err != nil {
		log.Fatal(err)
	}

	r := api.NewRequester(k, *u)
	j, err := api.Get(r, endpoint).Do()
	if err != nil {
		log.Fatal(err)
	}

	out, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(out)
	fmt.Println("")
}
