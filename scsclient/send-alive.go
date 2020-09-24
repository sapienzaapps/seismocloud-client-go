package scsclient

import (
	"errors"
	"fmt"
)

func (c *_clientimpl) SendAlive() error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/alive", c.opts.DeviceID), 0, false, fmt.Sprintf("%s;%s", c.opts.Model, c.opts.Version))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
