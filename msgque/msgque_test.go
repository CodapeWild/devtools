package msgque

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

const (
	fm1 = "foomsg1"
	fm2 = "foomsg2"
)

type FooMsg1 struct {
	MsgId int
	Callback
}

func (this *FooMsg1) Id() interface{} {
	return this.MsgId
}

func (this *FooMsg1) Type() interface{} {
	return fm1
}

type FooMsg2 struct {
	MsgId int
	Callback
}

func (this *FooMsg2) Id() interface{} {
	return this.MsgId
}

func (this *FooMsg2) Type() interface{} {
	return fm2
}

func fooMsgFanout(ticket interface{}, msg Message) {
	switch msg.Type() {
	case fm1:
		log.Println("foomsg1:", msg.Id())
		if msg.Put("process foomsg1 success") {
			log.Println("foomsg1 callback success")
		}
	case fm2:
		log.Println("foomsg2:", msg.Id())
		if msg.Put("process foomsg2 success") {
			log.Println("foomsg2 callback success")
		}
	}
}

func TestMsgQue(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	msgq := NewMessageQueue(SetTicket(NewTicketQueue(6)))
	msgq.StartUp(fooMsgFanout)

	for {
		if rand.Intn(100) > 49 {
			fm1 := &FooMsg1{
				MsgId:    rand.Intn(1000),
				Callback: NewCallbackQueue(time.Second),
			}
			if err := msgq.Send(fm1); err != nil {
				log.Println(err.Error())

				return
			}
			log.Println(fm1.Wait())
		} else {
			fm2 := &FooMsg2{
				MsgId:    rand.Intn(1000),
				Callback: NewCallbackQueue(time.Second),
			}
			if err := msgq.Send(fm2); err != nil {
				log.Println(err.Error())

				return
			}
			log.Println(fm2.Wait())
		}

		time.Sleep(time.Second)
	}
}
