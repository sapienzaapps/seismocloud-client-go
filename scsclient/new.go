package scsclient

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
)

// TODO
var MQTT_BROKER = "tls://mqtt-seismocloud.test.sapienzaapps.it"

func New(options ClientOptions) (Client, error) {
	if options.DeviceId == uuid.Nil {
		return nil, errors.New("Device ID is missing!")
	}

	mqttoptions := mqtt.NewClientOptions().AddBroker(MQTT_BROKER)
	mqttoptions.SetClientID(options.DeviceId.String())
	mqttoptions.SetAutoReconnect(false)

	// TODO: fix this
	mqttoptions.SetUsername("embedded")
	mqttoptions.SetPassword("embedded")

	mqttoptions.SetOrderMatters(true)
	mqttoptions.SetKeepAlive(10 * time.Second)
	mqttoptions.SetPingTimeout(10 * time.Second)
	mqttoptions.SetConnectTimeout(15 * time.Second)

	mqttc := mqtt.NewClient(mqttoptions)

	return &_clientimpl{
		opts:  options,
		mqttc: mqttc,
	}, nil
}
