package render

import "time"

type Thread interface {
	BeforeRender()
	Render()
	AfterRender()
	Interval() time.Duration
	Timeout() time.Duration
}
