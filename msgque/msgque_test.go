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

func TestSimpleTicketQueue(t *testing.T) {
	tq := NewSimpleTicketQueue(10)
	for i := 0; i < 10; i++ {
		go func() {
			tq.Fetch()
		}()
	}
}
