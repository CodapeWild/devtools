package clock

import (
	"time"
)

func ParseUnixMillsec(msec int64) time.Time {
	return time.Unix(0, msec*int64(time.Millisecond))
}

func ParseUnixSec(sec int64) time.Time {
	return ParseUnixMillsec(sec * 1000)
}
