package msgque

import "devtools/code"

type Ticket interface {
	Threads() int
	Fill()
	Generate() interface{}
	Retrieve() interface{}
	Recede(ticket interface{})
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
		this.Recede(this.Generate())
	}
}

func (this *TicketQueue) Generate() interface{} {
	return code.RandBase64(16)
}

func (this *TicketQueue) Retrieve() interface{} {
	return <-this.tickets
}

func (this *TicketQueue) Recede(ticket interface{}) {
	this.tickets <- ticket
}
