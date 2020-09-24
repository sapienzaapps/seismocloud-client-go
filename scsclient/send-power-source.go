package scsclient

import (
	"errors"
	"fmt"
)

func (c *_clientimpl) SendPowerSource(source PowerSource) error {
	if err := source.IsValid(); err != nil {
		return err
	}

	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/powersource", c.opts.DeviceID), 0, false, source)
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
