package scsclient

import (
	"crypto/tls"
	"github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"time"
)

func (c *clientV1impl) Connect() error {
	if c.mqttc != nil && c.mqttc.IsConnected() {
		return nil
	}

	if c.opts.Logger != nil {
		c.opts.Logger.Debugf("[%s] Connecting\n", c.opts.DeviceId)
	}

	c.opts.DeviceId = c.opts.DeviceId
	willpayload := make([]byte, 1+1+len(c.opts.DeviceId))
	willpayload[0] = API_DISCONNECT
	willpayload[1] = byte(len(c.opts.DeviceId))
	copy(willpayload[2:2+willpayload[1]], []byte(c.opts.DeviceId))

	// MQTT connection
	serverlist := make([]*url.URL, 1)
	serverlist[0], _ = url.Parse(c.opts.Server)
	clientOptions := mqtt.ClientOptions{
		AutoReconnect:  false,
		ClientID:       c.opts.ClientId,
		Servers:        serverlist,
		Username:       c.opts.User,
		Password:       c.opts.Pass,
		WillEnabled:    true,
		WillTopic:      "server",
		WillPayload:    willpayload,
		ConnectTimeout: 5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	if c.opts.SkipTLS {
		clientOptions.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	c.mqttc = mqtt.NewClient(&clientOptions)
	conntoken := c.mqttc.Connect()
	if conntoken.Wait() && conntoken.Error() != nil {
		return conntoken.Error()
	}

	c.mqttc.Subscribe("device-"+c.opts.DeviceId, 0, c.recvMessage)
	if c.opts.Logger != nil {
		c.opts.Logger.Debugf("[%s] Connected\n", c.opts.DeviceId)
	}

	c.Alive()
	c.aliveticker = time.NewTicker(15 * time.Minute)
	go func() {
		for range c.aliveticker.C {
			if !c.mqttc.IsConnected() {
				c.aliveticker.Stop()
				return
			}
			c.Alive()
		}
	}()
	return nil
}
