package scsclient

import (
	"errors"
	"time"
)

func NewClientV1(options ClientV1Options) (ClientV1, error) {
	if options.Server == "" || options.Version == "" || options.Model == "" || options.DeviceId == "" {
		return nil, errors.New("missing required parameters")
	}

	return &clientV1impl{
		opts:        options,
		mqttc:       nil,
		lastalive:   time.Unix(0, 0),
		aliveticker: nil,
		lasttime:    time.Unix(0, 0),
		timechan:    nil,
	}, nil
}
