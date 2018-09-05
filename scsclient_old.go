package scsclient

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"time"
)

type SCSClientOldProtocol interface {
	SetLocation(location SCSLocation)
	Connect() error
	Quake()
	//Move()
	Disconnect()
	GetTime() *time.Time
	IsConnected() bool
}

type SCSOldConfigCallback func(sigma float32)
type SCSOldRebootCallback func()
type SCSOldUpdateCallback func(hostname string, path string)

type scsClientOldProtoImpl struct {
	server         string
	clientID       string
	user           string
	pass           string
	deviceId       string
	mqttc          mqtt.Client
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

func (c *scsClientOldProtoImpl) IsConnected() bool {
	return c.mqttc.IsConnected()
}

func (c *scsClientOldProtoImpl) SetLocation(location SCSLocation) {
	c.location = location
}
