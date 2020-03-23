package msgque

import (
	"devtools/code"
	"devtools/comerr"
	"time"
)

type Message interface {
	Id() interface{}
	Type() interface{}
	Callback(cbMsg interface{}, timeout time.Duration) bool
}

type Ticket interface {
	Threads() int
	Generate() interface{}
	Retrieve() interface{}
	Recede(ticket interface{})
}

type TicketQueue struct {
	maxThrds int
	tickets  chan interface{}
}

func NewTicketQueue(maxThrds int) *TicketQueue {
	return &TicketQueue{
		maxThrds: maxThrds,
		tickets:  make(chan interface{}, maxThrds),
	}
}

func (this *TicketQueue) Threads() int {
	return this.maxThrds
}

func (this *TicketQueue) Generate() interface{} {
	return code.RandBase64(16)
}

func (this *TicketQueue) Retrieve() interface{} {
	return <-this.tickets
}

func (this *TicketQueue) Recede(ticket interface{}) {
	this.tickets <- ticket
}

const (
	def_max_buffer int           = 6
	def_timeout    time.Duration = time.Second
)

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
	for i := 0; i < this.Threads(); i++ {
		this.Recede(this.Generate())
	}

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
