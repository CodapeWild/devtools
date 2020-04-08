package msgque

import (
	"devtools/code"
)

type TicketQueue interface {
	MaxThreads() int
	Fill()
	Generate() interface{}
	Fetch() interface{}
	Restore(ticket interface{})
	Traverse(eachWithBreak func(ticket interface{}) bool)
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

func (this *SimpleTicketQueue) MaxThreads() int {
	return this.maxThrds
}

func (this *SimpleTicketQueue) Fill() {
	for i := len(this.tickets); i < this.MaxThreads(); i++ {
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

func (this *SimpleTicketQueue) Traverse(eachWithBreak func(ticket interface{}) bool) {
	for i := 0; i < len(this.tickets); i++ {
		ticket := <-this.tickets
		if eachWithBreak(ticket) {
			return
		}
		this.tickets <- ticket
	}
}
