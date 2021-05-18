package hostinfo

import (
	"os"
)

type HostInfo struct {
	Hostname string `json:"hostname"`
}

func Fetch() (*HostInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &HostInfo{Hostname: hostname}, nil
}
