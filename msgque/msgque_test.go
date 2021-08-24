package msgque

import (
	"log"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	log.SetFlags(log.LstdFlags | log.LstdFlags)
	log.SetOutput(os.Stdout)

	os.Exit(t.Run())
}

func TestMsgQue(t *testing.T) {
	tkque := NewSimpleTicketQueue(10)
	tkque.Fill()

	for i := 0; i < 100; i++ {
		go func(ticket Ticket, i int) {
			log.Printf("got ticket: %d\n", i)
		}(tkque.Fetch(), i)
	}

	select {}
}
