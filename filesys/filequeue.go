package filesys

import (
	"database/sql"
	"devtools/idflaker"
	"devtools/msgque"
	"log"
	"qrpool/common/comerr"
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
	qBuf, maxThrds int
	idflk          *idflaker.IdFlaker
	fqdb           *sql.DB
}

type FileQueueSetting func(fq *FileQueue)

func SetThreads(qBuf int, maxThreads int) FileQueueSetting {
	return func(fq *FileQueue) {

	}
}

func NewFileQueue() *FileQueue {

}

func (this *FileQueue) StartUp() {
	this.MessageQueue.StartUp(fileFanout)
}

func (this *FileQueue) Send(msg msgque.Message) error {
	return this.MessageQueue.Send(msg)
}

func (this *FileQueue) Close() error {

}

func fileFanout(ticket interface{}, msg msgque.Message) {
	switch msg.Type() {
	case save_file_msg:
		saveFile()
	case find_file_msg:
	case del_file_msg:
	default:
		log.Println(comerr.ParamTypeInvalid)
	}
}

func saveFile(ticket *DirTicket, msg *SaveMsg) {

}

func findFile(msg *FindMsg) {

}

func delFile(msg *DelMsg) {

}
