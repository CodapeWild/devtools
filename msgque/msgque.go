package msgque

import (
	"log"
	"sync"
	"time"

	"github.com/CodapeWild/devtools/cache"
)

const (
	def_que_buffer int = 6
)

type MsgQState int

const (
	MsgQ_Open MsgQState = iota + 1
	MsgQ_Congest
	MsgQ_Suspend
	MsgQ_Close
)

type MsgQueueStatus struct {
	MsgQueueState   MsgQState
	TicketsCapacity int
	TicketsLen      int
	IsCacheEnabled  bool
	CacheDepth      int
}

/*
	ticket: ticket queue token
		 msg: message
  closer: message queue main gorotine closer
*/
type FanoutHandler func(ticket Ticket, msg Message)

type MessageQueue struct {
	tq              TicketQueue
	cache           cache.Cache
	msgQue          chan Message
	queBuf          int
	timeout         time.Duration
	retryTimes      int
	retryDelay      func(ith int) time.Duration
	suspend, resume chan struct{}
	closer          chan struct{}
	msgqState       MsgQState
	sync.Mutex
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetMsgQTicket(tickq TicketQueue) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.tq = tickq
	}
}

func SetMsgQCache(cache cache.Cache) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.cache = cache
	}
}

func SetMsgQBuffer(size int) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.queBuf = size
	}
}

func SetMsgQTimeout(timeout time.Duration) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.timeout = timeout
	}
}

func SetMsgQRetry(times int, delay func(ith int) time.Duration) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.retryTimes = times
		msgq.retryDelay = delay
	}
}

func SetMsgQState(state MsgQState) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.msgqState = state
	}
}

func NewMessageQueue(opt ...MessageQueueSetting) (*MessageQueue, error) {
	msgq := &MessageQueue{
		queBuf:    def_que_buffer,
		suspend:   make(chan struct{}),
		resume:    make(chan struct{}),
		closer:    make(chan struct{}),
		msgqState: MsgQ_Open,
	}
	for _, v := range opt {
		v(msgq)
	}

	if msgq.tq == nil {
		return nil, ErrTicketQueueNil
	}
	msgq.tq.FillUp()

	msgq.msgQue = make(chan Message, msgq.queBuf)

	return msgq, nil
}

func (this *MessageQueue) StartUp(handler FanoutHandler) {
	// message queue main fanout goroutine
	go func() {
		for msg := range this.msgQue {
			// check close
			select {
			case <-this.closer:
				log.Println("message queue closed")

				return
			default:
			}
			// check suspend and wait for resume
			select {
			case <-this.suspend:
				log.Println("message queue suspended")
				<-this.resume
				log.Println("message queue resumed")
			default:
			}
			// check timeout cache up buffer
			go this.cleanCache()

			// spin go routine only if get ticket success
			ticket := <-this.getTicket()
			if ticket != ValidTicket {
				log.Println(ErrMsgQueClosed.Error())

				return
			}
			go func(ticket Ticket, msg Message) {
				handler(ticket, msg)
				this.tq.Restore() <- ticket
			}(ticket, msg)
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	// span go routine only if get ticket success
	ticket := <-this.getTicket()
	if ticket != InvalidTicket {
		return ErrMsgQueClosed
	}
	defer func() { this.tq.Restore() <- ticket }()

	go func() {
		ith := 0
		for {
			timer := time.NewTimer(this.timeout)
			select {
			case this.msgQue <- msg:
				timer.Stop()

				return
			case <-timer.C:
				if ith++; ith > this.retryTimes {
					log.Println(ErrMsgSendFailedAll.Error())

					return
				}
				time.Sleep(this.retryDelay(ith))
			}
		}
	}()

	return nil
}

//	Suspend will halt the message queue, no message can send or consume.
func (this *MessageQueue) Suspend() {
	if this.msgqState == MsgQ_Suspend {
		return
	}

	this.Lock()
	defer this.Unlock()

	if this.msgqState == MsgQ_Suspend {
		return
	}

	this.suspend <- struct{}{}
	this.msgqState = MsgQ_Suspend
}

// Resume will reopen message queue.
func (this *MessageQueue) Resume() {
	if this.msgqState == MsgQ_Open {
		return
	}

	this.Lock()
	defer this.Unlock()

	if this.msgqState == MsgQ_Open {
		return
	}

	// check cache queue
	go this.cleanCache()

	this.resume <- struct{}{}
}

func (this *MessageQueue) Close() {
	close(this.closer)
	this.msgqState = MsgQ_Close
}

// TODO:
func (this *MessageQueue) Status() *MsgQueueStatus {
	return &MsgQueueStatus{}
}

// getTicket will always get a ticket only if message queue closed.
func (this *MessageQueue) getTicket() <-chan Ticket {
	res := make(chan Ticket)
	for {
		timer := time.NewTimer(this.tq.Timeout())
		select {
		case <-this.closer:
			log.Println("message queue closed")
			res <- InvalidTicket

			return res
		case <-this.suspend:
			timer.Stop()
			log.Println("message queue suspended")
			<-this.resume
			log.Println("message queue resumed")
		case ticket := <-this.tq.Fetch():
			timer.Stop()
			res <- ticket

			return res
		case <-timer.C:
			log.Printf("fetch ticket timeout, sleep %ds and retry", this.tq.RetryDelay()/time.Second)
			time.Sleep(this.tq.RetryDelay())
		}
	}
}

func (this *MessageQueue) cleanCache() {
	if this.cache != nil && this.cache.Len() != 0 {
		for msg, ok := this.cache.Pop().(Message); ok && msg != nil; {
			this.Send(msg)
			msg, ok = this.cache.Pop().(Message)
		}
	}
}
