package msgque

import "time"

type Callback interface {
	Put(msg interface{}) bool
	Wait() (msg interface{})
}

type NoCallback struct{}

func (this NoCallback) Put(msg interface{}) bool { return false }

func (this NoCallback) Wait() interface{} { return nil }

type CallbackQueue struct {
	cbChan  chan interface{}
	timeout time.Duration
}

func NewCallbackQueue(timeout time.Duration) *CallbackQueue {
	return &CallbackQueue{
		cbChan:  make(chan interface{}),
		timeout: timeout,
	}
}

func (this *CallbackQueue) Put(msg interface{}) bool {
	if this.cbChan != nil {
		select {
		case <-time.After(this.timeout):
			return false
		case this.cbChan <- msg:
			return true
		}
	}

	return false
}

func (this *CallbackQueue) Wait() (msg interface{}) {
	if this.cbChan != nil {
		select {
		case <-time.After(this.timeout):
		case msg = <-this.cbChan:
		}
	}

	return
}
