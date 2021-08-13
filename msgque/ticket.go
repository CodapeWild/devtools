package msgque

type Ticket struct{}

type TicketQueue interface {
	Fill()
	Fetch() Ticket
	Restore(ticket Ticket)
	Len() int
	Cap() int
}

type SimpleTicketQueue chan Ticket

func NewSimpleTicketQueue(maxThreads int) SimpleTicketQueue {
	return make(chan Ticket, maxThreads)
}

func (this SimpleTicketQueue) Fill() {
	for i := this.Len(); i < this.Cap(); i++ {
		this.Restore(Ticket{})
	}
}

func (this SimpleTicketQueue) Fetch() Ticket {
	return <-this
}

func (this SimpleTicketQueue) Restore(ticket Ticket) {
	this <- ticket
}

func (this SimpleTicketQueue) Len() int {
	return len(this)
}

func (this SimpleTicketQueue) Cap() int {
	return cap(this)
}
