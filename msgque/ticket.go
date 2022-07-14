package msgque

import (
	"math/rand"
	"sync"
	"time"
)

type Ticket struct{ id uint64 }

type TicketQueue interface {
	// fetch one ticket from queue.
	Fetch() (Ticket, bool)
	// put one ticket back to queue.
	Restore(Ticket) error
	// return the current number of tickets in queue.
	Len() int
	// return the capacity tickets in queue if the given argument
	// is less than or equal to zero, otherwise, set the capacity to the given argument.
	Cap(int) int
}

type SimpleTicketQueue struct {
	tickets chan Ticket
	booth   map[Ticket]bool
	timeout time.Duration
	sync.Mutex
}

func NewSimpleTicketQueue(threads int, fetchTimeout time.Duration) SimpleTicketQueue {
	stq := SimpleTicketQueue{
		tickets: make(chan Ticket, threads),
		booth:   make(map[Ticket]bool, threads),
		timeout: fetchTimeout,
	}
	for threads > 0 {
		t := Ticket{id: rand.Uint64()}
		stq.tickets <- t
		stq.booth[t] = true

		threads--
	}

	return stq
}

func (stq SimpleTicketQueue) Fetch() (Ticket, bool) {
	if stq.timeout <= 0 {
		t := <-stq.tickets
		stq.booth[]
	}

	tick := time.NewTicker(stq.timeout)
	select {
	case <-tick.C:
		return Ticket{}, false
	case t := <-stq.tickets:
		return t, true
	}
}

func (stq SimpleTicketQueue) Restore() {

}

func (stq SimpleTicketQueue) Len() int {
	return len(stq.tickets)
}

func (stq SimpleTicketQueue) Cap(n int) int {
	if n > 0 {
		stq.Lock()
		defer stq.Unlock()

		stq.tickets = make(chan Ticket, n)
		var l = len(stq.tickets)
		if l > n {
			l = n
		}
		for i := 0; i < l; i++ {
			stq.tickets <- Ticket{}
		}

		return n
	} else {
		return cap(stq.tickets)
	}
}
