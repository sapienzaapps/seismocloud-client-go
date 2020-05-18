package scsclient

import (
	"fmt"
	"net"
)

func (c *_clientimpl) SendWiFiInfo(rssi float64, bssid net.HardwareAddr, essid string) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/essid", c.opts.DeviceId), 0, false, essid)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	token = c.mqttc.Publish(fmt.Sprintf("sensor/%s/rssi", c.opts.DeviceId), 0, false, fmt.Sprint(rssi))
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}

	token = c.mqttc.Publish(fmt.Sprintf("sensor/%s/bssid", c.opts.DeviceId), 0, false, bssid.String())
	token.Wait()
	return token.Error()
}
