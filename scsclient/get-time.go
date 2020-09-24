package scsclient

import (
	"errors"
	"fmt"
	"time"
)

func (c *_clientimpl) RequestTime() error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/timereq", c.opts.DeviceID), 0, false,
		fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
