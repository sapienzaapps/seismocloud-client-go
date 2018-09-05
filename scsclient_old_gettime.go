package scsclient

import "time"

func (c *scsClientOldProtoImpl) GetTime() *time.Time {
	if c.mqttc == nil || !c.mqttc.IsConnected() {
		return nil
	}
	if c.logger != nil {
		c.logger.Debugf("[%s] Time request\n", c.deviceId)
	}

	c.timechan = make(chan bool, 1)

	timepayload := make([]byte, 1+1+len(c.deviceId))
	timepayload[0] = API_TIMEREQ
	timepayload[1] = byte(len(c.deviceId))
	copy(timepayload[2:2+timepayload[1]], []byte(c.deviceId))

	c.mqttc.Publish("server", 2, false, timepayload).Wait()

	<-c.timechan

	return &c.lasttime
}
