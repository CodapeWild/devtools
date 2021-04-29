package msgque

import (
	"devtools/code"
	"sync"
	"time"
)

const (
	def_que_buffer   int           = 6
	def_msgq_timeout time.Duration = time.Second
)

type MsgQStatus int

const (
	MsgQ_Open MsgQStatus = iota + 1
	MsgQ_Congest
	MsgQ_Suspend
	MsgQ_Close
)

type FanoutHandler func(ticket interface{}, msg Message, closer chan struct{})

type critical struct {
	suspend, resume chan struct{}
	token           string
	sync.Mutex
}

type MessageQueue struct {
	tq      TicketQueue
	cache   Cache
	msgChan chan Message
	queBuf  int
	timeout time.Duration
	crtl    critical
	closer  chan struct{}
	status  MsgQStatus
}

type MessageQueueSetting func(msgQ *MessageQueue)

func SetMsgQTicket(tickq TicketQueue) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.tq = tickq
	}
}

func SetMsgQCache(cache Cache) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.cache = cache
	}
}

func SetMsgQBuffer(queBuf int) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.queBuf = queBuf
	}
}

func SetMsgQTimeout(timeout time.Duration) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.timeout = timeout
	}
}

func SetMsgQStatus(status MsgQStatus) MessageQueueSetting {
	return func(msgq *MessageQueue) {
		msgq.status = status
	}
}

func NewMessageQueue(opt ...MessageQueueSetting) *MessageQueue {
	msgQ := &MessageQueue{
		queBuf:  def_que_buffer,
		timeout: def_msgq_timeout,
		crtl: critical{
			suspend: make(chan struct{}),
			resume:  make(chan struct{}),
		},
		closer: make(chan struct{}),
		status: MsgQ_Open,
	}
	for _, v := range opt {
		v(msgQ)
	}

	msgQ.msgChan = make(chan Message, msgQ.queBuf)

	return msgQ
}

func (this *MessageQueue) StartUp(fanout FanoutHandler) {
	// populate ticket queue
	this.tq.Fill()

	// message queue main goroutine
	go func() {
		for v := range this.msgChan {
			// check close
			select {
			case <-this.closer:
				return
			default:
			}
			// check suspend
			select {
			case <-this.crtl.suspend:
				select {
				case <-this.crtl.resume:
					// clean up cache if not empty
					if this.cache != nil && this.cache.Len() != 0 {
						this.cleanCache()
					}
				}
			default:
			}

			if v.MustInvoice() {
				go func(ticket interface{}, msg Message) {
					fanout(ticket, msg, this.closer)
					this.tq.Restore(ticket)
				}(this.tq.Fetch(), v)
			} else {
				go fanout(nil, v, this.closer)
			}
		}
	}()
}

func (this *MessageQueue) Send(msg Message) error {
	if this.status == MsgQ_Close {
		return ErrMsgQClosed
	}

	var err error
	select {
	case <-time.After(this.timeout):
		err = ErrMsgQEnqueOvertime
		// cache up if exists
		if this.cache != nil {
			if !this.cache.Push(msg) {
				err = ErrCachePushFailed
			}
		}
	case this.msgChan <- msg:
		return nil
	}

	return err
}

func (this *MessageQueue) Suspend() (string, bool) {
	var (
		token string
		ok    bool
	)
	if this.crtl.token == "" {
		this.crtl.Lock()
		defer this.crtl.Unlock()

		this.crtl.suspend <- struct{}{}
		this.crtl.token = code.RandBase64(32)

		token = this.crtl.token
		ok = true

		this.status = MsgQ_Suspend
	}

	return token, ok
}

func (this *MessageQueue) Resume(token string) bool {
	var ok bool
	if this.crtl.token == token {
		this.crtl.Lock()
		defer this.crtl.Unlock()

		this.crtl.resume <- struct{}{}
		this.crtl.token = ""
		this.status = MsgQ_Open

		ok = true
	}

	return ok
}

func (this *MessageQueue) Close() {
	close(this.closer)
	this.status = MsgQ_Close
}

func (this *MessageQueue) Status() MsgQStatus {
	return this.status
}

func (this *MessageQueue) cleanCache() {
	go func() {
		msg := this.cache.Pop()
		for msg != nil {
			this.msgChan <- msg.(Message)
			msg = this.cache.Pop()
		}
	}()
}
