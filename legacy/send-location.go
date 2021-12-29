package legacy

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

func (c *_clientimpl) SendLocation(latitude decimal.Decimal, longitude decimal.Decimal) error {
	token := c.mqttc.Publish(fmt.Sprintf("sensor/%s/location", c.opts.DeviceID), 0, false, fmt.Sprintf("%s;%s", latitude, longitude))
	if !token.WaitTimeout(clientTimeout) {
		return errors.New("command timeout")
	}
	return token.Error()
}
