package scsclient

import (
	"errors"
	"fmt"
	"net"
)

func (c *_clientimpl) SendLocalIP(localAddr net.IP) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/localip", c.opts.DeviceID), 0, false, localAddr.String())
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
