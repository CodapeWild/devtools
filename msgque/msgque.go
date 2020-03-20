package msgque

import (
	"devtools/code"
	"devtools/comerr"
	"time"
)

type MsgState int

const (
	Msg_ProcSuccess MsgState = iota + 1
	Msg_ProcFailed
)

type CallbackMsg struct {
	Id      interface{}
	Type    interface{}
	State   int
	Err     error
	Payload interface{}
}

type Message interface {
	Id() interface{}
	Type() interface{}
	Callback(cbMsg *CallbackMsg, timeout time.Duration) bool
}

type TicketManager interface {
	Threads() int
	Generate() interface{}
	Retrieve() interface{}
	Recede(ticket interface{}, msg Message)
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

func (this *TicketQueue) Recede(ticket interface{}, msg Message) {
	this.tickets <- ticket
}

const (
	def_max_buffer int           = 6
	def_timeout    time.Duration = 3
)

type FanoutHandler func(ticket interface{}, msg Message)

type MessageQueue struct {
	TicketManager
	msgC    chan Message
	qBuf    int
	timeout time.Duration
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetTicketManager(mngr TicketManager) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.TicketManager = mngr
	}
}

func SetQueueBuffer(qBuf int) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.qBuf = qBuf
	}
}

func SetTimeout(timeout time.Duration) MessageQueueSetting {
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

	msgQ.msgC = make(chan Message, msgQ.qBuf)

	return msgQ
}

func (this *MessageQueue) StartMsgQueue(out FanoutHandler) {
	for i := 0; i < this.Threads(); i++ {
		this.Recede(this.Generate(), nil)
	}

	go func() {
		for v := range this.msgC {
			tick := this.Retrieve()
			go func(tick interface{}, msg Message) {
				out(tick, msg)
				this.Recede(tick, msg)
			}(tick, v)
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	select {
	case <-time.After(this.timeout):
		return comerr.Overtime
	case this.msgC <- msg:
		return nil
	}
}
