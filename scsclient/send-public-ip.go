package scsclient

import (
	"fmt"
	"net"
)

func (c *_clientimpl) SendPublicIP(localAddr net.IP) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/publicip", c.opts.DeviceID), 0, false, localAddr.String())
	token.WaitTimeout(clientTimeout)
	return token.Error()
}
