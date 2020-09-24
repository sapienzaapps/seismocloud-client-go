package scsclient

import (
	"fmt"
	"time"
)

func (c *_clientimpl) RequestTime() error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/timereq", c.opts.DeviceID), 0, false,
		fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	token.WaitTimeout(clientTimeout)
	return token.Error()
}
