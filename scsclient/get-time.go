package scsclient

import "time"

func (c *_clientimpl) GetTime() (time.Time, error) {
	return time.Now(), nil
}
