package msgque

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

type LaughMsg struct {
	id int
}

func (this *LaughMsg) Id() interface{} {
	return this.id
}

func (this *LaughMsg) Type() interface{} {
	return 1
}

func (this *LaughMsg) Callback(*CallbackMsg, time.Duration) bool {
	return false
}

type SmileMsg struct {
}

func (this *SmileMsg) Id() interface{} {
	return 2
}

func (this *SmileMsg) Type() interface{} {
	return 2
}

func (this *SmileMsg) Callback(*CallbackMsg, time.Duration) bool {
	return false
}

func Director(tick interface{}, msg Message) {
	switch msg.Type() {
	case 1:
		log.Println(tick)
		log.Println(msg.Id())
	case 2:
		log.Println(tick)
		log.Println(msg.Id())
	}
}

type FooTicketQueue struct {
	*TicketQueue
}

func (this *FooTicketQueue) Generate() interface{} {
	return 123
}

func TestMsgQ(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	msgQ := NewMessageQueue(SetTicketManager(&FooTicketQueue{NewTicketQueue(6)}))
	msgQ.StartMsgQueue(Director)

	// log.Println(msgQ.Retrieve())
	// log.Println(msgQ.Retrieve())
	// log.Println(msgQ.Retrieve())

	for i := 0; i < 3; i++ {
		go func(i int) {
			for {
				if rand.Intn(100) > 49 {
					if err := msgQ.Send(&LaughMsg{}); err != nil {
						log.Println(err.Error())
					}
				} else {
					if err := msgQ.Send(&SmileMsg{}); err != nil {
						log.Println(err.Error())
					}
				}
			}
		}(i)
	}

	// for i := 0; i < 100; i++ {
	// 	if err := msgQ.Send(&LaughMsg{id: i}); err != nil {
	// 		log.Println(err.Error())
	// 	}
	// }

	select {}
}
