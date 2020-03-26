package msgque

import (
	"devtools/comerr"
	"time"
)

const (
	def_max_buffer int           = 6
	def_timeout    time.Duration = time.Second
)

type Callback interface {
}

type Message interface {
	Id() interface{}
	Type() interface{}
	Callback(cbMsg interface{}, timeout time.Duration) bool
}

type FanoutHandler func(ticket interface{}, msg Message)

type MessageQueue struct {
	Ticket
	msgChan chan Message
	qBuf    int
	timeout time.Duration
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetTicket(tick Ticket) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.Ticket = tick
	}
}

func SetQueueBuffer(qBuf int) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.qBuf = qBuf
	}
}

func SetQueueTimeout(timeout time.Duration) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.timeout = timeout
	}
}

func NewMessageQueue(opt ...MessageQueueSetting) *MessageQueue {
	msgQ := &MessageQueue{
		qBuf:    def_max_buffer,
		timeout: def_timeout,
	}
	for _, v := range opt {
		v(msgQ)
	}

	msgQ.msgChan = make(chan Message, msgQ.qBuf)

	return msgQ
}

func (this *MessageQueue) StartUp(fanout FanoutHandler) {
	this.Fill()

	go func() {
		for v := range this.msgChan {
			tick := this.Retrieve()
			go func(tick interface{}, msg Message) {
				fanout(tick, msg)
				this.Recede(tick)
			}(tick, v)
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	select {
	case <-time.After(this.timeout):
		return comerr.Overtime
	case this.msgChan <- msg:
		return nil
	}
}
