package legacy

import (
	"errors"
	"fmt"
)

func (c *_clientimpl) SendBattery(batteryLevel float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/battery", c.opts.DeviceID), 0, false, fmt.Sprintf("%f", batteryLevel))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
