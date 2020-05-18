package scsclient

import "net"

func (c *_clientimpl) SendWiFiInfo(rssi float64, bssid net.HardwareAddr, essid string) error {

	// TODO: send mqtt WiFi update
	return nil
}
