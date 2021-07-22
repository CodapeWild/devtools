package fileque

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/CodapeWild/devtools/msgque"
)

func TestFileQue(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fq, err := NewFileQueue(SetDirectory("./data", 0744, 3), SetDbPath("./data/fq.db"), SetThreads(6, 6), SetTimeout(time.Second))
	if err != nil {
		log.Fatalln(err.Error())
	}
	fq.StartUp()

	for i := 0; i < 10; i++ {
		go func() {
			for {
				msg := &SaveMsg{
					Buf:      []byte("hello,tnt"),
					FileMode: 0644,
					Ext:      "txt",
					Callback: msgque.NewSimpleCallback(time.Second, nil),
				}
				if err = fq.Send(msg); err != nil {
					log.Panicln(err.Error())
				}
				cb := msg.Wait()
				log.Println(cb.(*CallbackMsg))

				time.Sleep(time.Second)
			}
		}()
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-c
	log.Println("fq closing")
	fq.Close()
}
