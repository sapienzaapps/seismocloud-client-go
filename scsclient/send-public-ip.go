package scsclient

import (
	"errors"
	"fmt"
	"net"
)

func (c *_clientimpl) SendPublicIP(localAddr net.IP) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/publicip", c.opts.DeviceID), 0, false, localAddr.String())
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
