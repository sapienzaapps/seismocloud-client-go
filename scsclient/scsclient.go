// Package scsclient implements a client for the SeismoCloud IoT network.
package scsclient

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"net"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// Client is the SeismoCloud client interface
type Client interface {
	// Connect to the SeismoCloud network
	Connect() error

	// IsConnected check if it's connected to the SeismoCloud network
	IsConnected() bool

	// GetDeviceID returns the configured Device ID
	GetDeviceID() uuid.UUID

	// SendAlive sends ALIVE command manually (IMPORTANT: this is periodically called by an internal ticker, so normally
	// there is no need to call alive manually)
	SendAlive() error

	// SendTemperature sends the new value for temperature sensor (if available)
	SendTemperature(temp float64) error

	// SendBattery sends the level of battery (if available)
	SendBattery(batteryLevel float64) error

	// SendPowerSource sends the current power source (if available)
	SendPowerSource(source PowerSource) error

	// SendLocation sends the current location (if available)
	SendLocation(latitude decimal.Decimal, longitude decimal.Decimal) error

	// Quake send the QUAKE message when a new vibration is detected
	Quake(quaketime time.Time, x float64, y float64, z float64) error

	// RequestTime requests the current time from the SCS network
	RequestTime() error

	// SendStreamData sends stream data (if enabled) to the server
	SendStreamData(datatime time.Time, x float64, y float64, z float64) error

	// SendLocalIP sends local IP address information to the server
	SendLocalIP(localAddr net.IP) error

	// SendPublicIP sends public IP address information to the server
	SendPublicIP(publicAddr net.IP) error

	// SendWiFiInfo sends WiFi information (if applicable) to the server
	SendWiFiInfo(rssi float64, bssid net.HardwareAddr, essid string) error

	// SendThreshold sends current threshold to the server
	SendThreshold(threshold float64) error

	// Close the connection gracefully
	Close() error
}

// ClientOptions represent all options for the SeismoCloud Client
type ClientOptions struct {
	// DeviceID holds the device ID
	DeviceID uuid.UUID

	// Model of this sensor. For example: esp8266, uno, etc. Check the documentation for possible values
	Model string

	// Version of the software
	Version string

	// OnNewSigma is executed when a new sigma value is received
	OnNewSigma func(Client, float64)

	// OnReboot is executed when a reboot command is received
	OnReboot func(Client)

	// OnStreamCommand is executed when the stream command ("on" or "off") is received. The second parameter indicates
	// whether the stream should be started (true) or stopped (false)
	OnStreamCommand func(Client, bool)

	// OnProbeSpeedSet s executed when a new probe speed is received. The second parameter is the new frequency of
	// probing (in Hz)
	OnProbeSpeedSet func(Client, int64)

	// OnTimeReceived is executed when a new time is received
	OnTimeReceived func(Client, int64, int64, int64, int64)

	// SeismoCloudBroker is the broker URL (in the form: tcp://hostname:port or tls://hostname:port)
	SeismoCloudBroker string

	// Username is the seismoCloud broker username
	Username string

	// Password is the seismoCloud broker password
	Password string

	// TLSConfig holds the custom configuration for TLS. Use this for inject the certificate
	TLSConfig *tls.Config

	// LocalDiscovery indicates whether the local (LAN) discovery part is enabled. When enabled, the sensor replies to
	// scan probes from apps
	// Not implemented: LocalDiscovery bool

	// Logger is a logger for the library
	Logger logrus.FieldLogger
}
