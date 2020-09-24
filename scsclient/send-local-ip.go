package scsclient

import (
	"fmt"
	"net"
)

func (c *_clientimpl) SendLocalIP(localAddr net.IP) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/localip", c.opts.DeviceID), 0, false, localAddr.String())
	token.WaitTimeout(clientTimeout)
	return token.Error()
}
