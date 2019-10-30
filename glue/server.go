package glue

import (
	"devtools/db/redisdb"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"
)

const (
	terminator_ = "###glue_terminator"
)

var (
	processorUnregistered = errors.New("unregistered processor")
)

type Handler interface {
	Process(req *Request, resp *Response)
}

type HandlerFunc func(req *Request, resp *Response)

func (this HandlerFunc) Process(req *Request, resp *Response) {
	this(req, resp)
}

type Server struct {
	key         string
	timeout     time.Duration
	workers     []chan []byte
	multiplexer map[string][]Handler
	rdsWrapper  *redisdb.RedisWrapper
}

func NewServer(server string, maxThreads int, timeout time.Duration) (*Server, error) {
	if maxThreads <= 0 {
		maxThreads = 1000
	}
	s := &Server{
		timeout:     timeout,
		workers:     make([]chan []byte, maxThreads),
		multiplexer: make(map[string][]Handler),
	}
	var err error
	s.key, err = formatServerKey(server)

	return s, err
}

func (this *Server) Register(tag string, handlers ...Handler) {
	for k := range handlers {
		this.multiplexer[tag] = append(this.multiplexer[tag], handlers[k])
	}
}

func (this *Server) workRoutine(i int) {
	for buf := range this.workers[i] {
		req := &Request{}
		err := json.Unmarshal(buf, req)
		if err != nil {
			log.Println(err.Error())

			continue
		}

		var tag string
		if _, tag, err = req.GetRemote(); err != nil {
			log.Println(err.Error())

			continue
		}

		if handlers, ok := this.multiplexer[tag]; !ok {
			log.Println(processorUnregistered.Error())
		} else {
			resp := NewResponse(req, this.rdsWrapper)
			for k := range handlers {
				handlers[k].Process(req, resp)
			}
		}
	}
}

func (this *Server) timeoutResp(req []byte) {
	// todo: add glue server timeout
}

func (this *Server) Listen(conf *redisdb.RedisConfig) error {
	pool, err := conf.NewPool()
	if err != nil {
		return err
	}
	this.rdsWrapper = redisdb.NewWrapper(pool)

	for i := 0; i < len(this.workers); i++ {
		go func(i int) {
			this.workers[i] = make(chan []byte)
			this.workRoutine(i)
		}(i)
	}

	// clear server pip
	this.rdsWrapper.DelKey(this.key)

	log.Println("glue server online")
	rand.Seed(time.Now().UnixNano())
	for {
		rply, err := this.rdsWrapper.BLPop(this.key, 0)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		buf := rply.([]interface{})[1].([]byte)
		if len(buf) == len(terminator_) && string(buf) == terminator_ {
			this.rdsWrapper.DelKey(this.key)
			log.Println("glue server offline")
			break
		}

		i := rand.Intn(len(this.workers))
		select {
		case this.workers[i] <- buf:
		case <-time.After(this.timeout):
			log.Println("glue server timeout")
		}
	}

	return nil
}

func (this *Server) Terminate() {
	if this.rdsWrapper != nil {
		this.rdsWrapper.LPush(this.key, terminator_)
	}
}
