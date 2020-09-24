package scsclient

import (
	"errors"
	"fmt"
	"net"
)

func (c *_clientimpl) SendWiFiInfo(rssi float64, bssid net.HardwareAddr, essid string) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/essid", c.opts.DeviceID), 0, false, essid)
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	if token.Error() != nil {
		return token.Error()
	}

	token = c.mqttc.Publish(fmt.Sprintf("sensor/%s/rssi", c.opts.DeviceID), 0, false, fmt.Sprint(rssi))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	if token.Error() != nil {
		return token.Error()
	}

	token = c.mqttc.Publish(fmt.Sprintf("sensor/%s/bssid", c.opts.DeviceID), 0, false, bssid.String())
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
