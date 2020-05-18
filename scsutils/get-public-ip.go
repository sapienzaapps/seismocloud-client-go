package scsutils

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func GetPublicIP() (net.IP, error) {
	client := new(http.Client)

	client.Timeout = 10 * time.Second
	request, err := http.NewRequest("GET", "https://api.ipify.org", nil)
	if err != nil {
		return net.IP{}, nil
	}

	resp, err := client.Do(request)
	if err != nil {
		return net.IP{}, nil
	}

	ipaddr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return net.IP{}, nil
	}

	return net.ParseIP(string(ipaddr)), nil
}
