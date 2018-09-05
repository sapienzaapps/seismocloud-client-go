package scsclient

import (
	"github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"time"
)

func (c *scsClientOldProtoImpl) Connect() error {
	if c.mqttc != nil && c.mqttc.IsConnected() {
		return nil
	}

	if c.logger != nil {
		c.logger.Debugf("[%s] Connecting\n", c.deviceId)
	}

	c.deviceId = c.deviceId
	willpayload := make([]byte, 1+1+len(c.deviceId))
	willpayload[0] = API_DISCONNECT
	willpayload[1] = byte(len(c.deviceId))
	copy(willpayload[2:2+willpayload[1]], []byte(c.deviceId))

	// MQTT connection
	serverlist := make([]*url.URL, 1)
	serverlist[0], _ = url.Parse(c.server)
	clientOptions := mqtt.ClientOptions{
		AutoReconnect: false,
		ClientID:      c.clientID,
		Servers:       serverlist,
		Username:      c.user,
		Password:      c.pass,
		WillEnabled:   true,
		WillTopic:     "server",
		WillPayload:   willpayload,
	}

	c.mqttc = mqtt.NewClient(&clientOptions)
	conntoken := c.mqttc.Connect()
	if conntoken.Wait() && conntoken.Error() != nil {
		return conntoken.Error()
	}

	c.mqttc.Subscribe("device-"+c.deviceId, 0, c.recvMessage)
	if c.logger != nil {
		c.logger.Debugf("[%s] Connected\n", c.deviceId)
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
