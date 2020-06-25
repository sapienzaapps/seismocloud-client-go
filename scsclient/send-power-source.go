package scsclient

import "fmt"

func (c *_clientimpl) SendPowerSource(source PowerSource) error {
	if err := source.IsValid(); err != nil {
		return err
	}

	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/powersource", c.opts.DeviceId), 0, false, source)
	token.Wait()
	return token.Error()
}
