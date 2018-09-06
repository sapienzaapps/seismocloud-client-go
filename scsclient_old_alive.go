package scsclient

import (
	"encoding/binary"
	"math"
	"time"
)

func (c *clientV1impl) Alive() {
	if c.mqttc == nil || !c.mqttc.IsConnected() {
		return
	}

	if c.opts.Location.IsValid() {
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] Alive with position: %f %f\n", c.opts.DeviceId, c.opts.Location.Lat, c.opts.Location.Lng)
		}
	} else {
		if c.opts.Logger != nil {
			c.opts.Logger.Debugf("[%s] Alive\n", c.opts.DeviceId)
		}
	}
	modellen := len(c.opts.Model)
	versionlen := len(c.opts.Version)

	if c.mqttc == nil || !c.mqttc.IsConnected() {
		if c == nil || c.aliveticker == nil {
			c.aliveticker.Stop()
			c.aliveticker = nil
		}
		return
	}
	alivepayloadlen := 1 + 1 + len(c.opts.DeviceId) + 1 + modellen + 1 + versionlen
	if c.opts.Location.IsValid() {
		alivepayloadlen += 8
	}

	alivepayload := make([]byte, alivepayloadlen)
	j := 0
	if c.opts.Location.IsValid() {
		alivepayload[j] = API_KEEPALIVE_POSITION
	} else {
		alivepayload[j] = API_KEEPALIVE
	}
	j++

	// Device ID
	alivepayload[j] = byte(len(c.opts.DeviceId))
	j++
	j += copy(alivepayload[j:j+len(c.opts.DeviceId)], []byte(c.opts.DeviceId))

	// Model
	alivepayload[j] = byte(modellen)
	j++
	j += copy(alivepayload[j:j+modellen], []byte(c.opts.Model))

	// Version
	alivepayload[j] = byte(versionlen)
	j++
	j += copy(alivepayload[j:j+4], []byte(c.opts.Version))

	if c.opts.Location.IsValid() {
		binary.LittleEndian.PutUint32(alivepayload[j:j+4], math.Float32bits(float32(c.opts.Location.Lat)))
		j += 4

		binary.LittleEndian.PutUint32(alivepayload[j:j+4], math.Float32bits(float32(c.opts.Location.Lng)))
		j += 4
	}

	c.mqttc.Publish("server", 2, false, alivepayload).Wait()
	c.lastalive = time.Now()
}
