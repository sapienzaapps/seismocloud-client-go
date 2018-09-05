package scsclient

import "time"

type SCSClientOldProtocol interface {
	Alive()
	SetLocation(location SCSLocation)
	Connect(deviceid string, server string, clientID string, user string, pass string) error
	Quake()
	//Move()
	Disconnect()
	GetTime() time.Time
}

type SCSOldConfigCallback func(sigma float32)
type SCSOldRebootCallback func()
type SCSOldUpdateCallback func(hostname string, path string)
