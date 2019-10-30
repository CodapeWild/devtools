package glue

import (
	"devtools/db/redisdb"
	"log"
	"os"
	"testing"
	"time"
)

var (
	rdsConf = &redisdb.RedisConfig{
		Host: "127.0.0.1",
		Port: "6389",
	}
	remote_server = "testserver"
	r             = "foo"
	local_client  = "testclient"
)

type Param struct {
	Name string
	Age  int
}

var p = &Param{
	Name: "tnt",
	Age:  123,
}

type Ticket struct {
	Permission bool
	Count      int
}

func routine(req *Request, resp *Response) {
	param := &Param{}
	err := req.GetPayload(param)
	if err != nil {
		log.Fatalln(err.Error())
	}

	t := &Ticket{
		Permission: true,
		Count:      321,
	}
	resp.SetStatus(StatusOK)
	if err = resp.WriteBack(t); err != nil {
		log.Fatalln(err.Error())
	}
}

func callback(resp *Response) {
	t := &Ticket{}
	if err := resp.GetPayload(t); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(t.Count, t.Permission)
}

func TestGlue(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	server, err := NewServer(remote_server, 100, time.Second)
	if err != nil {
		log.Fatalln(err.Error())
	}
	server.Register(r, HandlerFunc(routine))
	go func() {
		if err := server.Listen(rdsConf); err != nil {
			log.Fatalln(err.Error())
		}
	}()
	time.Sleep(time.Second)

	req := NewRequest()
	req.SetRemote(remote_server, r).SetLocal(remote_server, local_client).SetPayload(p).SetCallback(callback)
	c, err := NewClient(rdsConf)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for i := 0; i < 10; i++ {
		if err = c.Do(req); err != nil {
			log.Fatalln(err.Error())
		}
		if i == 3 {
			server.Terminate()
		}
		time.Sleep(time.Second)
	}

	select {}
}
