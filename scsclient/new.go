package scsclient

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
)

func New(options ClientOptions) (Client, error) {
	if options.DeviceId == uuid.Nil {
		return nil, errors.New("device ID is missing")
	}

	mqttoptions := mqtt.NewClientOptions()

	mqttoptions.AddBroker(options.SeismoCloudBroker)
	if options.Username != "" {
		mqttoptions.SetUsername(options.Username)
		mqttoptions.SetPassword(options.Password)
	}
	if options.TLSConfig != nil {
		mqttoptions.SetTLSConfig(options.TLSConfig)
	}

	mqttoptions.SetClientID(options.DeviceId.String())
	mqttoptions.SetAutoReconnect(false)

	mqttoptions.SetOrderMatters(true)
	mqttoptions.SetKeepAlive(10 * time.Second)
	mqttoptions.SetPingTimeout(10 * time.Second)
	mqttoptions.SetConnectTimeout(15 * time.Second)

	mqttc := mqtt.NewClient(mqttoptions)

	return &_clientimpl{
		opts:            options,
		mqttc:           mqttc,
		aliveTickerStop: make(chan int),
	}, nil
}
