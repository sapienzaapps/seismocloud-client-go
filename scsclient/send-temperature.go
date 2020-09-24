package scsclient

import (
	"errors"
	"fmt"
)

func (c *_clientimpl) SendTemperature(temp float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/temperature", c.opts.DeviceID), 0, false,
		fmt.Sprintf("%f", temp))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
