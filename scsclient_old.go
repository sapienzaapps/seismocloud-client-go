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
