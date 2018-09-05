package scsclient

import (
	"encoding/binary"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"math"
	"net/url"
	"time"
)

const (
	API_KEEPALIVE          = 1
	API_QUAKE              = 2
	API_TIMEREQ            = 3
	API_TIMERESP           = 4
	API_CFG                = 5
	API_DISCONNECT         = 6
	API_TEMPERATURE        = 7
	API_REBOOT             = 8
	API_KEEPALIVE_POSITION = 9
)

type scsClientOldProtoImpl struct {
	mqttc          mqtt.Client
	deviceId       string
	lastalive      time.Time
	aliveticker    *time.Ticker
	sigma          float32
	lasttime       time.Time
	timechan       chan bool
	noiseticker    *time.Ticker
	location       SCSLocation
	logger         *logging.Logger
	cfgcallback    SCSOldConfigCallback
	rebootcallback SCSOldRebootCallback
	updatecallback SCSOldUpdateCallback
}

func NewSCSClientOldProtocol(logger *logging.Logger, cfgcb SCSOldConfigCallback, rebootcb SCSOldRebootCallback, updatecb SCSOldUpdateCallback) SCSClientOldProtocol {
	return &scsClientOldProtoImpl{
		mqttc:          nil,
		deviceId:       "",
		lastalive:      time.Unix(0, 0),
		aliveticker:    nil,
		sigma:          3,
		lasttime:       time.Unix(0, 0),
		timechan:       nil,
		location:       SCSLocation{0, 0},
		cfgcallback:    cfgcb,
		rebootcallback: rebootcb,
		updatecallback: updatecb,
	}
}

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

func (c *scsClientOldProtoImpl) Alive() {
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

func (c *scsClientOldProtoImpl) SetLocation(location SCSLocation) {
	c.location = location
}

func (c *scsClientOldProtoImpl) Connect(deviceid string, server string, clientID string, user string, pass string) error {
	if c.logger != nil {
		c.logger.Debugf("[%s] Connecting\n", deviceid)
	}

	c.deviceId = deviceid
	willpayload := make([]byte, 1+1+len(deviceid))
	willpayload[0] = API_DISCONNECT
	willpayload[1] = byte(len(deviceid))
	copy(willpayload[2:2+willpayload[1]], []byte(deviceid))

	// MQTT connection
	serverlist := make([]*url.URL, 1)
	serverlist[0], _ = url.Parse(server)
	clientOptions := mqtt.ClientOptions{
		AutoReconnect: false,
		ClientID:      clientID,
		Servers:       serverlist,
		Username:      user,
		Password:      pass,
		WillEnabled:   true,
		WillTopic:     "server",
		WillPayload:   willpayload,
	}

	c.mqttc = mqtt.NewClient(&clientOptions)
	conntoken := c.mqttc.Connect()
	if conntoken.Wait() && conntoken.Error() != nil {
		return conntoken.Error()
	}

	c.mqttc.Subscribe("device-"+deviceid, 0, c.recvMessage)
	if c.logger != nil {
		c.logger.Debugf("[%s] Connected\n", deviceid)
	}

	c.Alive()
	c.aliveticker = time.NewTicker(15 * time.Minute)
	go func() {
		for range c.aliveticker.C {
			c.Alive()
		}
	}()
	return nil
}

func (c *scsClientOldProtoImpl) Quake() {
	quakepayload := make([]byte, 1+1+len(c.deviceId))
	quakepayload[0] = API_QUAKE
	quakepayload[1] = byte(len(c.deviceId))
	copy(quakepayload[2:2+quakepayload[1]], []byte(c.deviceId))

	if c.mqttc != nil {
		c.mqttc.Publish("server", 2, false, quakepayload).Wait()
	}
}

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

func (c *scsClientOldProtoImpl) GetTime() time.Time {
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

	return c.lasttime
}
