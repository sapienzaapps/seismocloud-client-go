package scsclient

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func (c *_clientimpl) SendLocation(latitude decimal.Decimal, longitude decimal.Decimal) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/location", c.opts.DeviceId), 0, false, fmt.Sprintf("%s;%s", latitude, longitude))
	token.Wait()
	return token.Error()
}
