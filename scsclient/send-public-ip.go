package scsclient

import (
	"fmt"
	"net"
)

func (c *_clientimpl) SendPublicIP(localAddr net.IP) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/publicip", c.opts.DeviceId), 0, false, localAddr.String())
	token.Wait()
	return token.Error()
}