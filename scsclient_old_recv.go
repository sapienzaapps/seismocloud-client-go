package scsclient

import (
	"encoding/binary"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

func (c *clientV1impl) recvMessage(sock mqtt.Client, m mqtt.Message) {
	payload := m.Payload()
	switch payload[0] {
	case API_CFG:
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] CFG\n", c.opts.DeviceId)
		}
		if c.opts.ConfigCallback != nil {
			c.opts.ConfigCallback(float32frombytes(payload[1:5]))
		}

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

		if hostname != "" && path != "" && c.opts.UpdateCallback != nil {
			// Do update
			c.opts.UpdateCallback(hostname, path)
		}
	case API_REBOOT:
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] Reboot\n", c.opts.DeviceId)
		}
		if c.opts.RebootCallback != nil {
			c.opts.RebootCallback()
		}
	case API_TIMERESP:
		c.lasttime = time.Unix(int64(binary.LittleEndian.Uint32(payload[1:])), 0)
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] Time response\n", c.opts.DeviceId)
		}
		if c.timechan != nil {
			c.timechan <- true
		}
	default:
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] Unknown message\n", c.opts.DeviceId)
		}
	}
}
