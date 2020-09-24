package scsclient

import "fmt"

func (c *_clientimpl) SendTemperature(temp float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/temperature", c.opts.DeviceID), 0, false,
		fmt.Sprintf("%f", temp))
	token.WaitTimeout(clientTimeout)
	return token.Error()
}
