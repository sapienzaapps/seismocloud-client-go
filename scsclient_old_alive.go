package scsclient

import (
	"encoding/binary"
	"math"
	"time"
)

func (c *scsClientOldProtoImpl) Alive() {
	if c.mqttc == nil || !c.mqttc.IsConnected() {
		return
	}

	if c.location.IsValid() {
		if c.logger != nil {
			c.logger.Debugf("[%s] Alive with position: %f %f\n", c.deviceId, c.location.Lat, c.location.Lng)
		}
	} else {
		if c.logger != nil {
			c.logger.Debugf("[%s] Alive\n", c.deviceId)
		}
	}
	model := "linux-x86"
	modellen := len(model)

	if c.mqttc == nil || !c.mqttc.IsConnected() {
		if c == nil || c.aliveticker == nil {
			c.aliveticker.Stop()
			c.aliveticker = nil
		}
		return
	}
	alivepayloadlen := 1 + 1 + len(c.deviceId) + 1 + modellen + 1 + 4
	if c.location.IsValid() {
		alivepayloadlen += 8
	}

	alivepayload := make([]byte, alivepayloadlen)
	j := 0
	if c.location.IsValid() {
		alivepayload[j] = API_KEEPALIVE_POSITION
	} else {
		alivepayload[j] = API_KEEPALIVE
	}
	j++

	// Device ID
	alivepayload[j] = byte(len(c.deviceId))
	j++
	j += copy(alivepayload[j:j+len(c.deviceId)], []byte(c.deviceId))

	// Model
	alivepayload[j] = byte(modellen)
	j++
	j += copy(alivepayload[j:j+modellen], []byte(model))

	// Version
	alivepayload[j] = 4
	j++
	j += copy(alivepayload[j:j+4], []byte("0.00"))

	if c.location.IsValid() {
		binary.LittleEndian.PutUint32(alivepayload[j:j+4], math.Float32bits(float32(c.location.Lat)))
		j += 4

		binary.LittleEndian.PutUint32(alivepayload[j:j+4], math.Float32bits(float32(c.location.Lng)))
		j += 4
	}

	c.mqttc.Publish("server", 2, false, alivepayload).Wait()
	c.lastalive = time.Now()
}
