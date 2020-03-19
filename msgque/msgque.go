package msgque

import (
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
	Type() interface{}
	Id() interface{}
	Callback(cbMsg *CallbackMsg, timeout time.Duration) bool
}

const (
	def_max_buffer  int           = 6
	def_max_threads int           = 6
	def_timeout     time.Duration = 3
)

type Fanout func(msg Message)

type MessageQueue struct {
	qBuf     int
	msgQ     chan Message
	maxThrds int
	thrds    chan int
	timeout  time.Duration
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetQueueBuffer(qBuf int) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.qBuf = qBuf
	}
}

func SetMaxThreads(maxThrds int) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.maxThrds = maxThrds
	}
}

func SetTimeout(timeout time.Duration) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.timeout = timeout
	}
}

func NewMessageQueue(opt ...MessageQueueSetting) *MessageQueue {
	msgQ := &MessageQueue{
		qBuf:     def_max_buffer,
		maxThrds: def_max_threads,
		timeout:  def_timeout,
	}
	for _, v := range opt {
		v(msgQ)
	}

	msgQ.msgQ = make(chan Message, msgQ.qBuf)
	msgQ.thrds = make(chan int, msgQ.maxThrds)

	return msgQ
}

func (this *MessageQueue) StartMsgQueue(out Fanout) {
	go func() {
		for v := range this.msgQ {
			this.thrds <- 1
			go func(msg Message) {
				out(msg)
				<-this.thrds
			}(v)
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	select {
	case <-time.After(this.timeout):
		return comerr.Overtime
	case this.msgQ <- msg:
		return nil
	}
}
