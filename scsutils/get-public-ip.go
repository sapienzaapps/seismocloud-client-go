package scsutils

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// GetPublicIP retrieve the public IP address
func GetPublicIP() (net.IP, error) {
	client := new(http.Client)

	client.Timeout = 10 * time.Second
	request, err := http.NewRequest("GET", "https://api.ipify.org", nil)
	if err != nil {
		return net.IPv4zero, nil
	}

	resp, err := client.Do(request)
	if err != nil {
		return net.IPv4zero, nil
	}

	ipaddr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return net.IPv4zero, nil
	}

	return net.ParseIP(string(ipaddr)), nil
}
