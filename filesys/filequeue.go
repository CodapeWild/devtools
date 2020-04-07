package filesys

import (
	"database/sql"
	"devtools/idflaker"
	"devtools/msgque"
	"time"
)

const (
	def_top_dir      = "./data"
	def_max_threads  = 3
	def_timeout      = time.Second
	def_dir_contains = 30
)

type FileQueue struct {
	*msgque.MessageQueue
	idflk *idflaker.IdFlaker
	fqdb  *sql.DB
}

type FileQueueSetting func(fq *FileQueue)

func NewFileQueue() *FileQueue {
	fq := &FileQueue{}
	fq.MessageQueue = msgque.NewMessageQueue()
}

func (this *FileQueue) StartUp() {
	this.MessageQueue.StartUp(this.fileFanout)
}

func (this *FileQueue) Send(msg msgque.Message) error {
	return this.MessageQueue.Send(msg)
}

func (this *FileQueue) Close() error {

}

func (this *FileQueue) fileFanout(ticket interface{}, msg msgque.Message) {

}

func (this *FileQueue) saveFile(ticket *DirTicket, msg *SaveMsg) {

}

func (this *FileQueue) findFile(msg *FindMsg) {

}
