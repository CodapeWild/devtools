package msgque

import (
	"devtools/code"
)

type TicketQueue interface {
	Threads() int
	Fill()
	Generate() interface{}
	Fetch() interface{}
	Restore(ticket interface{})
	Traverse(process func(ticket interface{}))
}

type SimpleTicketQueue struct {
	maxThrds int
	tickets  chan interface{}
}

func NewSimpleTicketQueue(maxThrds int) *SimpleTicketQueue {
	return &SimpleTicketQueue{
		maxThrds: maxThrds,
		tickets:  make(chan interface{}, maxThrds),
	}
}

func (this *SimpleTicketQueue) Threads() int {
	return this.maxThrds
}

func (this *SimpleTicketQueue) Fill() {
	for i := len(this.tickets); i < this.Threads(); i++ {
		this.Restore(this.Generate())
	}
}

func (this *SimpleTicketQueue) Generate() interface{} {
	return code.RandBase64(15)
}

func (this *SimpleTicketQueue) Fetch() interface{} {
	return <-this.tickets
}

func (this *SimpleTicketQueue) Restore(ticket interface{}) {
	this.tickets <- ticket
}

func (this *SimpleTicketQueue) Traverse(process func(ticket interface{})) {
	for i := 0; i < len(this.tickets); i++ {
		ticket := <-this.tickets
		process(ticket)
		this.tickets <- ticket
	}
}
