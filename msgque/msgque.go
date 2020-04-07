package msgque

import (
	"devtools/comerr"
	"log"
	"time"
)

const (
	def_que_buffer   int           = 6
	def_send_timeout time.Duration = time.Second
)

type Message interface {
	Id() interface{}
	Type() interface{}
	MustFetch() bool
	Callback
}

type FanoutHandler func(ticket interface{}, msg Message)

type MessageQueue struct {
	TicketQueue
	msgChan     chan Message
	queBuf      int
	sendTimeout time.Duration
	suspending  bool
	resume      chan int
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetTicket(tick TicketQueue) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.TicketQueue = tick
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
		suspending:  false,
		resume:      make(chan int),
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
			if this.suspending {
				log.Println(comerr.Suspending.Error())
				<-this.resume
				log.Println("process resumed")
			}

			if v.MustFetch() {
				go func(ticket interface{}, msg Message) {
					fanout(ticket, msg)
					this.Restore(ticket)
				}(this.Fetch(), v)
			} else {
				go fanout(nil, v)
			}
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	select {
	case <-time.After(this.sendTimeout):
		if this.suspending {
			select {
			case this.msgChan <- msg:
				return nil
			}
		}

		return comerr.Overtime
	case this.msgChan <- msg:
		return nil
	}
}

func (this *MessageQueue) Suspend() {
	this.suspending = true
}

func (this *MessageQueue) Resume() {
	this.suspending = false
	this.resume <- 1
}
