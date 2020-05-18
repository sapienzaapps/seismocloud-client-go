package scsclient

import "fmt"

func (c *_clientimpl) SendTemperature(temp float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/temperature", c.opts.DeviceId), 0, false,
		fmt.Sprintf("%f", temp))
	token.Wait()
	return token.Error()
}
