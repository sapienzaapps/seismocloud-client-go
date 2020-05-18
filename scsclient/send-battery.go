package scsclient

import "fmt"

func (c *_clientimpl) SendBattery(batteryLevel float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/battery", c.opts.DeviceId), 0, false, fmt.Sprintf("%f", batteryLevel))
	token.Wait()
	return token.Error()
}
