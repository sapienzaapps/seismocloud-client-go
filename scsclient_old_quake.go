package scsclient

func (c *scsClientOldProtoImpl) Quake() {
	if c.mqttc == nil || !c.mqttc.IsConnected() {
		return
	}
	quakepayload := make([]byte, 1+1+len(c.deviceId))
	quakepayload[0] = API_QUAKE
	quakepayload[1] = byte(len(c.deviceId))
	copy(quakepayload[2:2+quakepayload[1]], []byte(c.deviceId))

	c.mqttc.Publish("server", 2, false, quakepayload).Wait()
}
