package msgque

import "devtools/code"

type Ticket interface {
	Threads() int
	Fill()
	Generate() interface{}
	Fetch() interface{}
	Restore(ticket interface{})
}

type TicketQueue struct {
	maxThrds int
	tickets  chan interface{}
}

func NewTicketQueue(maxThrds int) *TicketQueue {
	return &TicketQueue{
		maxThrds: maxThrds,
		tickets:  make(chan interface{}, maxThrds),
	}
}

func (this *TicketQueue) Threads() int {
	return this.maxThrds
}

func (this *TicketQueue) Fill() {
	for i := 0; i < this.Threads(); i++ {
		this.Restore(this.Generate())
	}
}

func (this *TicketQueue) Generate() interface{} {
	return code.RandBase64(15)
}

func (this *TicketQueue) Fetch() interface{} {
	return <-this.tickets
}

func (this *TicketQueue) Restore(ticket interface{}) {
	this.tickets <- ticket
}
