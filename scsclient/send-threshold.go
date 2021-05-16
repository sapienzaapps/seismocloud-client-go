package scsclient

import (
	"errors"
	"fmt"
)

func (c *_clientimpl) SendThreshold(threshold float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/threshold", c.opts.DeviceID), 0, false, fmt.Sprintf("%f", threshold))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
