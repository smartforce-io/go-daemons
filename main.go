package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"github.com/smartforce-io/go-daemons/hostinfo"
	"log"
	"net/http"
	"os"
)

var (
	url = ""

	errWrongResponse = errors.New("a request returned wrong response")
)

func init() {
	flag.StringVar(&url, "url", "", "a target url")
}

func sendHostInfo(info *hostinfo.HostInfo) error {
	b, err := json.Marshal(info)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errWrongResponse
	}
	return nil
}

func main() {
	log.Println("Go Daemons")

	flag.Parse()

	if len(url) == 0 {
		log.Print("missing required argument `--url`\nRun --help")
		os.Exit(2)
	}

	info, err := hostinfo.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if err := sendHostInfo(info); err != nil {
		log.Fatal(err)
	}
	log.Println("Host Info was sent correctly!")
}
