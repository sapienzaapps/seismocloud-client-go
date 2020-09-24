package scsclient

import "fmt"

func (c *_clientimpl) SendThreshold(threshold float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/treshold", c.opts.DeviceID), 0, false, fmt.Sprintf("%f", threshold))
	token.WaitTimeout(clientTimeout)
	return token.Error()
}
