package msgque

import (
	"devtools/code"
	"errors"
	"log"
	"testing"
	"time"
)

type FooMsgType int

const (
	Foo1 FooMsgType = iota + 1
	Foo2
)

type Foo1Msg struct {
	MsgId string
	CbC   chan interface{}
}

func (this *Foo1Msg) Id() interface{} {
	return this.MsgId
}

func (this *Foo1Msg) Type() interface{} {
	return Foo1
}

func (this *Foo1Msg) Callback(cbMsg interface{}, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case this.CbC <- cbMsg:
		return true
	}
}

type Foo2Msg struct {
	MsgId string
	CbC   chan interface{}
}

func (this *Foo2Msg) Id() interface{} {
	return this.MsgId
}

func (this *Foo2Msg) Type() interface{} {
	return Foo2
}

func (this *Foo2Msg) Callback(cbMsg interface{}, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case this.CbC <- cbMsg:
		return true
	}
}

type FooCbMsg struct {
	MsgId      string
	MsgState   int
	Err        error
	MsgPayload interface{}
}

func (this *FooCbMsg) Id() interface{} {
	return this.MsgId
}

func (this *FooCbMsg) State() interface{} {
	return this.MsgState
}

func (this *FooCbMsg) Error() error {
	return this.Err
}

func (this *FooCbMsg) Payload() interface{} {
	return this.Payload
}

type FooTicket struct {
	Code   string
	Expire int
}

type FooTicketQueue struct {
	*TicketQueue
}

func (this *FooTicketQueue) Generate() interface{} {
	return &FooTicket{
		Code:   code.RandBase64(15),
		Expire: 6,
	}
}

func (this *FooTicketQueue) Recede(ticket interface{}) {
	ftick, ok := ticket.(*FooTicket)
	if !ok {
		return
	}

	if ftick.Expire <= 0 {
		return
	}
	ftick.Expire--

	this.TicketQueue.Recede(ftick)
}

var count int = 0

func FooFanout(ticket interface{}, msg Message) {
	count++
	switch msg.Type() {
	case Foo1:
		log.Println("******Foo1******")
		msg.Callback(&FooCbMsg{
			MsgId:      msg.Id().(string),
			MsgState:   200,
			Err:        nil,
			MsgPayload: nil,
		}, 6)
	case Foo2:
		log.Println("******Foo2******")
		msg.Callback(&FooCbMsg{
			MsgId:      msg.Id().(string),
			MsgState:   404,
			Err:        errors.New("can not found page"),
			MsgPayload: nil,
		}, 3)
	default:
	}
}

func TestMsgQ(t *testing.T) {
	ftickQ := &FooTicketQueue{NewTicketQueue(3)}
	msgQ := NewMessageQueue(SetTicket(ftickQ), SetQueueBuffer(6), SetQueueTimeout(time.Second))
	msgQ.StartUp(FooFanout)

	go func() {
		for {
			cbc := make(chan interface{})
			msgQ.Send(&Foo1Msg{MsgId: code.RandBase64(9), CbC: cbc})
			log.Println(<-cbc)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			cbc := make(chan interface{})
			msgQ.Send(&Foo2Msg{MsgId: code.RandBase64(9), CbC: cbc})
			log.Println(<-cbc)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			log.Println(count)
			time.Sleep(3 * time.Second)
		}
	}()
	log.Println("###############", count)

	select {}
}
