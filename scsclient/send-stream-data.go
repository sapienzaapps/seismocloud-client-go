package scsclient

import "time"

func (c *_clientimpl) SendStreamData(datatime time.Time, x float64, y float64, z float64) error {
	// TODO: send mqtt temperature update
	return nil
}
