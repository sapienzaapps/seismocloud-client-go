package scsclient

import "fmt"

func (c *_clientimpl) SendThreshold(threshold float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/treshold", c.opts.DeviceId), 0, false, fmt.Sprintf("%f", threshold))
	token.Wait()
	return token.Error()
}
