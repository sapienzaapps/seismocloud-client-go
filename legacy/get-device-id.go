package legacy

func (c *_clientimpl) GetDeviceID() string {
	return c.opts.DeviceID
}
