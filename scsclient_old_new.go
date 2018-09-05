package scsclient

import (
	"github.com/op/go-logging"
	"time"
)

func NewSCSClientOldProtocol(deviceid string, server string, clientID string, user string, pass string, logger *logging.Logger,
	cfgcb SCSOldConfigCallback, rebootcb SCSOldRebootCallback, updatecb SCSOldUpdateCallback) SCSClientOldProtocol {

	return &scsClientOldProtoImpl{
		server:         server,
		clientID:       clientID,
		user:           user,
		pass:           pass,
		deviceId:       deviceid,
		logger:         logger,
		mqttc:          nil,
		lastalive:      time.Unix(0, 0),
		aliveticker:    nil,
		sigma:          3,
		lasttime:       time.Unix(0, 0),
		timechan:       nil,
		location:       SCSLocation{0, 0},
		cfgcallback:    cfgcb,
		rebootcallback: rebootcb,
		updatecallback: updatecb,
	}
}
