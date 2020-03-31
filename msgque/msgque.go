package msgque

import (
	"devtools/comerr"
	"time"
)

type Callback interface {
	Put(msg interface{}) bool
	Wait() (msg interface{})
}

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

const (
	def_que_buffer   int           = 6
	def_send_timeout time.Duration = time.Second
)

type Message interface {
	Id() interface{}
	Type() interface{}
	Callback
}

type FanoutHandler func(ticket interface{}, msg Message)

type MessageQueue struct {
	Ticket
	msgChan     chan Message
	queBuf      int
	sendTimeout time.Duration
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetTicket(tick Ticket) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.Ticket = tick
	}
}

func SetQueueBuffer(queBuf int) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.queBuf = queBuf
	}
}

func SetSendTimeout(timeout time.Duration) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.sendTimeout = timeout
	}
}

func NewMessageQueue(opt ...MessageQueueSetting) *MessageQueue {
	msgQ := &MessageQueue{
		queBuf:      def_que_buffer,
		sendTimeout: def_send_timeout,
	}
	for _, v := range opt {
		v(msgQ)
	}

	msgQ.msgChan = make(chan Message, msgQ.queBuf)

	return msgQ
}

func (this *MessageQueue) StartUp(fanout FanoutHandler) {
	this.Fill()

	go func() {
		for v := range this.msgChan {
			go func(tick interface{}, msg Message) {
				fanout(tick, msg)
				this.Restore(tick)
			}(this.Fetch(), v)
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	select {
	case <-time.After(this.sendTimeout):
		return comerr.Overtime
	case this.msgChan <- msg:
		return nil
	}
}
