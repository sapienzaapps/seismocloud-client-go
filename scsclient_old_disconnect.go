package scsclient

func (c *scsClientOldProtoImpl) Disconnect() {
	if c.mqttc != nil && c.mqttc.IsConnected() {
		if c.logger != nil {
			c.logger.Debugf("[%s] Disconnecting\n", c.deviceId)
		}

		alivemovepayload := make([]byte, 1+1+len(c.deviceId))
		alivemovepayload[0] = API_DISCONNECT
		alivemovepayload[1] = byte(len(c.deviceId))
		copy(alivemovepayload[2:2+alivemovepayload[1]], []byte(c.deviceId))

		c.mqttc.Publish("server", 2, false, alivemovepayload).Wait()

		c.aliveticker.Stop()
		c.aliveticker = nil
		c.mqttc.Disconnect(0)
		c.mqttc = nil
	}
}
