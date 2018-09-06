package scsclient

import "time"

func (c *clientV1impl) GetTime() *time.Time {
	if c.mqttc == nil || !c.mqttc.IsConnected() {
		return nil
	}
	if c.opts.Logger != nil {
		c.opts.Logger.Debugf("[%s] Time request\n", c.opts.DeviceId)
	}

	c.timechan = make(chan bool, 1)

	timepayload := make([]byte, 1+1+len(c.opts.DeviceId))
	timepayload[0] = API_TIMEREQ
	timepayload[1] = byte(len(c.opts.DeviceId))
	copy(timepayload[2:2+timepayload[1]], []byte(c.opts.DeviceId))

	c.mqttc.Publish("server", 2, false, timepayload).Wait()

	<-c.timechan

	return &c.lasttime
}
