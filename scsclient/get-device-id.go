package scsclient

import "github.com/gofrs/uuid"

func (c *_clientimpl) GetDeviceID() uuid.UUID {
	return c.opts.DeviceID
}
