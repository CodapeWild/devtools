package msgque

import (
	"time"
)

type Callback interface {
	Call(rslt interface{})
	CallWithTimeout(rslt interface{}, timeout time.Duration) error
	Wait(timeout time.Duration) (rslt interface{}, err error)
}
type Message interface {
	Id() interface{}   // used as fanout identify
	Type() interface{} // use as fanout identity
	MustInvoice() bool // message can be put into message queue wihtout fetching a ticket
	Callback           // message processing result callback
}

type NoCallback struct{}

func NewNoCallback() *NoCallback {
	return &NoCallback{}
}

func (this NoCallback) Call(interface{}) bool { return false }

func (this NoCallback) CallWithTimeout(interface{}, time.Duration) error { return nil }

func (this NoCallback) Wait(time.Duration) (error, interface{}) { return nil, nil }

type SimpleCallback chan interface{}

func NewSimpleCallback() SimpleCallback {
	return make(chan interface{})
}

func (this SimpleCallback) Call(rslt interface{}) {
	if this != nil {
		this <- rslt
	}
}

func (this SimpleCallback) CallWithTimeout(rslt interface{}, timeout time.Duration) (err error) {
	if this != nil {
		select {
		case this <- rslt:
		case <-time.After(timeout):
			err = ErrCallbackSendTimeout
		}
	}

	return
}

func (this SimpleCallback) Wait(timeout time.Duration) (rslt interface{}, err error) {
	if this != nil {
		select {
		case rslt = <-this:
		case <-time.After(timeout):
			err = ErrCallbackReceiveTimeout
		}
	}

	return
}
