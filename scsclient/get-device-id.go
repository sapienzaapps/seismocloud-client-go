package scsclient

import uuid "github.com/satori/go.uuid"

func (c *_clientimpl) GetDeviceID() uuid.UUID {
	return c.opts.DeviceID
}
