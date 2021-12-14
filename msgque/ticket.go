package msgque

import "time"

const (
	InvalidTicket Ticket = iota - 1
	ValidTicket          = 1
)

type Ticket int

type TicketQueue interface {
	FillUp()
	Fetch() <-chan Ticket
	Restore() chan<- Ticket
	Len() int
	Cap() int
	Timeout() time.Duration
	RetryDelay() time.Duration
}

type SimpleTicketQueue chan Ticket

func NewSimpleTicketQueue(maxThreads int) SimpleTicketQueue {
	return make(chan Ticket, maxThreads)
}

func (this SimpleTicketQueue) FillUp() {
	for this.Len() < this.Cap() {
		this.Restore() <- ValidTicket
	}
}

func (this SimpleTicketQueue) Fetch() <-chan Ticket {
	return this
}

func (this SimpleTicketQueue) Restore() chan<- Ticket {
	return this
}

func (this SimpleTicketQueue) Len() int {
	return len(this)
}

func (this SimpleTicketQueue) Cap() int {
	return cap(this)
}

func (this SimpleTicketQueue) Timeout() time.Duration {
	return 3 * time.Second
}

func (this SimpleTicketQueue) RetryDelay() time.Duration {
	return time.Second
}
