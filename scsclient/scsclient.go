package scsclient

import (
	"crypto/tls"
	"net"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

// Client is the SeismoCloud client interface
type Client interface {
	// Connect to the SeismoCloud network
	Connect() error

	// Check if it's connected to the SeismoCloud network
	IsConnected() bool

	// Returns the configured Device ID
	GetDeviceID() uuid.UUID

	// Send alive manually (IMPORTANT: this is periodically called by an internal ticker, so normally there is no need
	// to call alive manually)
	SendAlive() error

	// This function sends the new value for temperature sensor (if available)
	SendTemperature(temp float64) error

	// This function sends the level of battery (if available)
	SendBattery(batteryLevel float64) error

	// This function sends the current power source (if available)
	SendPowerSource(source PowerSource) error

	// This function sends the current location (if available)
	SendLocation(latitude decimal.Decimal, longitude decimal.Decimal) error

	// This function send the QUAKE message when a new vibration is detected
	Quake(quaketime time.Time, x float64, y float64, z float64) error

	// Retrieve the current time from SCS network
	RequestTime() error

	// Send stream data (if enabled)
	SendStreamData(datatime time.Time, x float64, y float64, z float64) error

	// Setnd local IP address information
	SendLocalIP(localAddr net.IP) error

	// Send public IP address information
	SendPublicIP(publicAddr net.IP) error

	// Send WiFi information (if applicable)
	SendWiFiInfo(rssi float64, bssid net.HardwareAddr, essid string) error

	// Send current threshold
	SendThreshold(threshold float64) error

	// Close the connection gracefully
	Close() error
}

// ClientOptions represent all options for the SeismoCloud Client
type ClientOptions struct {
	// Device ID
	DeviceID uuid.UUID

	// Model of this sensor. For example: esp8266, uno, etc. Check the documentation
	Model string

	// Version of the software
	Version string

	// Function to execute when a new sigma value is received
	OnNewSigma func(Client, float64)

	// Function to execute when a reboot command is received
	OnReboot func(Client)

	// Function to execute when the stream command ("on" or "off") is received
	// The second parameter indicates whether the stream should be started (true)
	// or not (false)
	OnStreamCommand func(Client, bool)

	// Function to execute when a new probe speed is received. The second
	// parameter is the new frequency of probing (in Hz)
	OnProbeSpeedSet func(Client, int64)

	// Function to execute when a new time is received
	OnTimeReceived func(Client, int64, int64, int64, int64)

	// SeismoCloud broker URL
	// For tests/dev, use: tls://mqtt-seismocloud.test.sapienzaapps.it
	SeismoCloudBroker string

	// SeismoCloud broker Username
	// For tests/dev, use: embedded
	Username string

	// SeismoCloud broker Password
	// For tests/dev, use: embedded
	Password string

	// Custom CA root store for connection
	TLSConfig *tls.Config
}
