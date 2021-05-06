package msgque

type TicketQueue interface {
	MaxThreads() int
	Fill()
	Generate() interface{}
	Fetch() interface{}
	Restore(ticket interface{})
	Traverse(eachWithBreak func(ticket interface{}) bool)
}

type SimpleTicketQueue chan struct{}

func NewSimpleTicketQueue(maxThreads int) SimpleTicketQueue {
	return make(chan struct{}, maxThreads)
}

func (this SimpleTicketQueue) MaxThreads() int {
	return cap(this)
}

func (this SimpleTicketQueue) Fill() {
	for i := len(this); i < cap(this); i++ {
		this.Restore(this.Generate())
	}
}

func (this SimpleTicketQueue) Generate() struct{} {
	return struct{}{}
}

func (this SimpleTicketQueue) Fetch() struct{} {
	return <-this
}

func (this SimpleTicketQueue) Restore(ticket struct{}) {
	this <- ticket
}

func (this SimpleTicketQueue) Traverse(eachWithBreak func(ticket struct{}) bool) {
	for i := 0; i < len(this); i++ {
		ticket := <-this
		if eachWithBreak(ticket) {
			return
		}
		this <- ticket
	}
}
