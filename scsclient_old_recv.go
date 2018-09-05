package scsclient

import (
	"encoding/binary"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

func (c *scsClientOldProtoImpl) recvMessage(sock mqtt.Client, m mqtt.Message) {
	payload := m.Payload()
	switch payload[0] {
	case API_CFG:
		if c.logger != nil {
			c.logger.Debugf("[%s] CFG\n", c.deviceId)
		}
		c.cfgcallback(float32frombytes(payload[1:5]))

		hostname := ""
		path := ""
		offset := uint(5)
		hlen := payload[offset]
		offset++
		if hlen > 0 {
			hostname = string(payload[offset : offset+uint(hlen)])
			offset += uint(hlen)
		}
		plen := payload[offset]
		plen++
		if plen > 0 {
			path = string(payload[offset : offset+uint(hlen)])
			offset += uint(plen)
		}

		if hostname != "" && path != "" {
			// Do update
			c.updatecallback(hostname, path)
		}
	case API_REBOOT:
		if c.logger != nil {
			c.logger.Debugf("[%s] Reboot\n", c.deviceId)
		}
		c.rebootcallback()
	case API_TIMERESP:
		c.lasttime = time.Unix(int64(binary.LittleEndian.Uint32(payload[1:])), 0)
		if c.logger != nil {
			c.logger.Debugf("[%s] Time response\n", c.deviceId)
		}
		if c.timechan != nil {
			c.timechan <- true
		}
	default:
		if c.logger != nil {
			c.logger.Debugf("[%s] Unknown message\n", c.deviceId)
		}
	}
}
