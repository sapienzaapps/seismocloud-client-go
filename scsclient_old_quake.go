package scsclient

func (c *clientV1impl) Quake() {
	if c.mqttc == nil || !c.mqttc.IsConnected() {
		return
	}
	quakepayload := make([]byte, 1+1+len(c.opts.DeviceId))
	quakepayload[0] = API_QUAKE
	quakepayload[1] = byte(len(c.opts.DeviceId))
	copy(quakepayload[2:2+quakepayload[1]], []byte(c.opts.DeviceId))

	c.mqttc.Publish("server", 2, false, quakepayload).Wait()
}
