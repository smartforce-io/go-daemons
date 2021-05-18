package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sevlyar/go-daemon"
	"github.com/smartforce-io/go-daemons/hostinfo"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	ENV_URL = "GO_DAEMON_URL"
)

var (
	url              = ""
	errWrongResponse = errors.New("a request returned wrong response")
)

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

func worker() {
	info, err := hostinfo.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	if err := sendHostInfo(info); err != nil {
		log.Fatal(err)
	}
	log.Println("Host Info was sent correctly!")
}

func runner() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			worker()
		}
	}
}

func runLinuxDaemon() {
	cntxt := &daemon.Context{
		PidFileName: "go-daemon.pid",
		PidFilePerm: 0644,
		LogFileName: "go-daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	runner()
}

func main() {
	log.Println("Go Daemons")

	url = os.Getenv(ENV_URL)

	if len(url) == 0 {
		log.Printf("missing required environment variable %q", ENV_URL)
		os.Exit(2)
	}

	switch runtime.GOOS {
	case "linux":
		log.Print("Running Linux Daemon")
		runLinuxDaemon()
	default:
		log.Printf("This OS %q doesn't support yet.", runtime.GOOS)
	}
}
