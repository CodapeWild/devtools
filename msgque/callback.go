package msgque

import "time"

type Callback interface {
	Put(msg interface{}) bool
	Wait() (msg interface{})
}

type NoCallback struct{}

func (this NoCallback) Put(msg interface{}) bool { return false }

func (this NoCallback) Wait() interface{} { return nil }

type SimpleCallback struct {
	cbch    chan interface{}
	timeout time.Duration
}

func NewSimpleCallback(timeout time.Duration) *SimpleCallback {
	return &SimpleCallback{
		cbch:    make(chan interface{}),
		timeout: timeout,
	}
}

func (this *SimpleCallback) Put(msg interface{}) bool {
	if this.cbch != nil {
		select {
		case <-time.After(this.timeout):
			return false
		case this.cbch <- msg:
			return true
		}
	}

	return false
}

func (this *SimpleCallback) Wait() (msg interface{}) {
	if this.cbch != nil {
		select {
		case <-time.After(this.timeout):
		case msg = <-this.cbch:
		}
	}

	return
}
