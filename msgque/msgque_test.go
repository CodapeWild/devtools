package msgque

import (
	"log"
	"testing"
	"time"
)

type LaughMsg struct {
	id int
}

func (this *LaughMsg) Type() interface{} {
	return 1
}

func (this *LaughMsg) Id() interface{} {
	return this.id
}

func (this *LaughMsg) Callback(*CallbackMsg, time.Duration) bool {
	return false
}

type SmileMsg struct {
}

func (this *SmileMsg) Type() interface{} {
	return 2
}

func (this *SmileMsg) Id() interface{} {
	return 2
}

func (this *SmileMsg) Callback(*CallbackMsg, time.Duration) bool {
	return false
}

func Directory(msg Message) {
	switch msg.Type() {
	case 1:
		log.Println(msg.Id())
	case 2:
		log.Println(msg.Id())
	}
}

func TestMsgQ(t *testing.T) {
	msgQ := NewMessageQueue(SetQueueBuffer(3), SetMaxThreads(1), SetTimeout(time.Millisecond))
	msgQ.StartMsgQueue(Directory)

	// for i := 0; i < 3; i++ {
	// 	go func(i int) {
	// 		for {
	// 			if rand.Intn(100) > 49 {
	// 				msgQ.Send(&LaughMsg{})
	// 			} else {
	// 				msgQ.Send(&SmileMsg{})
	// 			}
	// 		}
	// 	}(i)
	// }

	for i := 0; i < 100; i++ {
		if err := msgQ.Send(&LaughMsg{id: i}); err != nil {
			log.Println(err.Error())
		}
	}

	select {}
}
