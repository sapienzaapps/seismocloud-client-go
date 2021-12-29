package legacy

import (
	"errors"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var deviceIDRx = regexp.MustCompile("^[A-Za-z0-9]+$")

// New creates a new instance of the SeismoCloud client - it does NOT connect the clien to the network yet
func New(options ClientOptions) (Client, error) {
	if !deviceIDRx.MatchString(options.DeviceID) {
		return nil, errors.New("device ID is not valid")
	}

	if options.Logger == nil {
		options.Logger = logrus.New()
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

	mqttoptions.SetClientID(options.DeviceID)
	mqttoptions.SetAutoReconnect(false)

	mqttoptions.SetOrderMatters(true)
	mqttoptions.SetKeepAlive(10 * time.Second)
	mqttoptions.SetPingTimeout(10 * time.Second)
	mqttoptions.SetConnectTimeout(15 * time.Second)

	mqttc := mqtt.NewClient(mqttoptions)

	return &_clientimpl{
		opts:            options,
		mqttc:           mqttc,
		aliveTickerStop: make(chan int, 1),
	}, nil
}
