package scsclient

import (
	"errors"
	"fmt"
)

func (c *_clientimpl) SendLocation(latitude float64, longitude float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/location", c.opts.DeviceID), 0, false, fmt.Sprintf("%f;%f", latitude, longitude))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
