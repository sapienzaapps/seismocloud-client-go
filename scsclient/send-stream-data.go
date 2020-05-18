package scsclient

import (
	"fmt"
	"time"
)

func (c *_clientimpl) SendStreamData(datatime time.Time, x float64, y float64, z float64) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/streamdata", c.opts.DeviceId), 0, false,
		fmt.Sprintf("%d;%f;%f;%f", datatime.UnixNano()/1000000, x, y, z))
	token.Wait()
	return token.Error()
}
