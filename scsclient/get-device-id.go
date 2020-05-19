package scsclient

import uuid "github.com/satori/go.uuid"

func (c *_clientimpl) GetDeviceId() uuid.UUID {
	return c.opts.DeviceId
}
