package scsutils

import (
	"time"
)

func SyncSNTPTime(t0 int64, t1 int64, t2 int64, t3 int64) time.Time {
	if t0 == 0 {
		// NTP can't handle 34 or 68 years differences, so if we start from 0 then return t2
		return time.Unix(0, t2*int64(time.Millisecond))
	}
	a := t1 - t0
	b := t2 - t3
	theta := (a + b) / int64(2)
	sigma := (t3 - t0) - (t2 - t1)

	return time.Unix(0, t0*int64(time.Millisecond)).Add(time.Duration(theta+sigma/int64(2)) * time.Millisecond)
}
