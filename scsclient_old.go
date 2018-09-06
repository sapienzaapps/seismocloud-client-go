package scsclient

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"time"
)

type ClientV1 interface {
	SetLocation(location SCSLocation)
	Connect() error
	Quake()
	//Move()
	Disconnect()
	GetTime() *time.Time
	IsConnected() bool

	// CAUTION: can lead to MITM - use only for testing
	SetSkipTLS(bool)
}

type V1ConfigCallback func(sigma float32)
type V1RebootCallback func()
type V1UpdateCallback func(hostname string, path string)

type ClientV1Options struct {
	Server   string
	ClientId string
	User     string
	Pass     string
	DeviceId string
	SkipTLS  bool

	Model    string
	Version  string
	Location SCSLocation

	Logger         *logging.Logger
	ConfigCallback V1ConfigCallback
	UpdateCallback V1UpdateCallback
	RebootCallback V1RebootCallback
}

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

type clientV1impl struct {
	opts        ClientV1Options
	mqttc       mqtt.Client
	lastalive   time.Time
	aliveticker *time.Ticker
	lasttime    time.Time
	timechan    chan bool
	noiseticker *time.Ticker
}

func (c *clientV1impl) IsConnected() bool {
	return c.mqttc.IsConnected()
}

func (c *clientV1impl) SetLocation(location SCSLocation) {
	c.opts.Location = location
}

func (c *clientV1impl) SetSkipTLS(b bool) {
	c.opts.SkipTLS = b
}
