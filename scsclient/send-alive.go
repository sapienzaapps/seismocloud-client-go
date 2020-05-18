package scsclient

import "fmt"

func (c *_clientimpl) SendAlive() error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/alive", c.opts.DeviceId), 0, false, fmt.Sprintf("%s;%s", c.opts.Model, c.opts.Version))
	token.Wait()
	return token.Error()
}
