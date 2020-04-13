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
	msgChan     chan Message
	queBuf      int
	sendTimeout time.Duration
	TicketQueue
	suspending bool
	resume     chan int
	closer     chan int
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetQueueBuffer(queBuf int) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.queBuf = queBuf
	}
}

func SetTimeout(timeout time.Duration) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.sendTimeout = timeout
	}
}

func SetTicket(tick TicketQueue) MessageQueueSetting {
	return func(msgQ *MessageQueue) {
		msgQ.TicketQueue = tick
	}
}

func NewMessageQueue(opt ...MessageQueueSetting) *MessageQueue {
	msgQ := &MessageQueue{
		queBuf:      def_que_buffer,
		sendTimeout: def_send_timeout,
		suspending:  false,
		resume:      make(chan int),
		closer:      make(chan int),
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
			if v == nil {
				close(this.msgChan)

				return
			}

			if this.suspending {
				log.Println("message queue suspending")
				<-this.resume
				log.Println("message queue resumed")
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
	if msg == nil {
		return comerr.ParamInvalid
	}

	select {
	case <-this.closer:
		return comerr.ChannelClosed
	default:
	}

	select {
	case <-this.closer:
		return comerr.ChannelClosed
	case <-time.After(this.sendTimeout):
		if this.suspending {
			return this.Send(msg)
		} else {
			return comerr.Overtime
		}
	case this.msgChan <- msg:
		return nil
	}
}

func (this *MessageQueue) Suspend() {
	this.suspending = true
}

func (this *MessageQueue) Resume() {
	this.resume <- 1
	this.suspending = false
}

func (this *MessageQueue) Close() {
	close(this.closer)
	this.msgChan <- nil
	<-this.msgChan
}
