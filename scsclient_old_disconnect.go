package scsclient

func (c *clientV1impl) Disconnect() {
	if c.mqttc != nil && c.mqttc.IsConnected() {
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] Disconnecting\n", c.opts.DeviceId)
		}

		alivemovepayload := make([]byte, 1+1+len(c.opts.DeviceId))
		alivemovepayload[0] = API_DISCONNECT
		alivemovepayload[1] = byte(len(c.opts.DeviceId))
		copy(alivemovepayload[2:2+alivemovepayload[1]], []byte(c.opts.DeviceId))

		c.mqttc.Publish("server", 2, false, alivemovepayload).Wait()

		c.aliveticker.Stop()
		c.aliveticker = nil
		c.mqttc.Disconnect(0)
		c.mqttc = nil
	}
}
