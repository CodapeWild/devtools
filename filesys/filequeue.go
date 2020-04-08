package filesys

import (
	"database/sql"
	"devtools/comerr"
	"devtools/idflaker"
	"devtools/msgque"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	def_top_dir      = "./data"
	def_dir_mode     = 0744
	def_dir_capacity = 30
	def_fqdb_path    = "./data/fq.db"
	def_que_buf      = 3
	def_max_threads  = 3
	def_timeout      = time.Second
)

type FileQueue struct {
	topDir         string
	dirMode        os.FileMode
	dirCapacity    int
	fqDbPath       string
	fqdb           *sql.DB
	idflk          *idflaker.IdFlaker
	qBuf, maxThrds int
	timeout        time.Duration
	*msgque.MessageQueue
}

type FileQueueSetting func(fq *FileQueue)

func SetDirectory(topDir string, dirMode os.FileMode, dirCapacity int) FileQueueSetting {
	return func(fq *FileQueue) {
		fq.topDir = topDir
		fq.dirMode = dirMode
		fq.dirCapacity = dirCapacity
	}
}

func SetDbPath(dbPath string) FileQueueSetting {
	return func(fq *FileQueue) {
		fq.fqDbPath = dbPath
	}
}

func SetThreads(qBuf, maxThrds int) FileQueueSetting {
	return func(fq *FileQueue) {
		fq.qBuf = qBuf
		fq.maxThrds = maxThrds
	}
}

func SetTimeout(timeout time.Duration) FileQueueSetting {
	return func(fq *FileQueue) {
		fq.timeout = timeout
	}
}

func NewFileQueue(opt ...FileQueueSetting) (*FileQueue, error) {
	fq := &FileQueue{
		topDir:      def_top_dir,
		dirMode:     def_dir_mode,
		dirCapacity: def_dir_capacity,
		fqDbPath:    def_fqdb_path,
		qBuf:        def_que_buf,
		maxThrds:    def_max_threads,
		timeout:     def_timeout,
	}
	for _, v := range opt {
		v(fq)
	}

	err := os.MkdirAll(fq.topDir, fq.dirMode)
	if err != nil {
		return nil, err
	}

	if fq.fqdb, err = sql.Open("sqlite3", fq.fqDbPath); err != nil {
		return nil, err
	}
	if err = createTable(fq.fqdb); err != nil {
		return nil, err
	}

	if fq.idflk, err = idflaker.NewIdFlaker(1); err != nil {
		return nil, err
	}

	dirque, err := NewDirTicketQueue(fq.maxThrds, fq.topDir, fq.dirMode, fq.dirCapacity, fq.fqdb)
	if err != nil {
		return nil, err
	}
	fq.MessageQueue = msgque.NewMessageQueue(msgque.SetQueueBuffer(fq.qBuf), msgque.SetTicket(dirque), msgque.SetSendTimeout(fq.timeout))

	return fq, nil
}

func (this *FileQueue) StartUp() {
	this.MessageQueue.StartUp(this.fileFanout)
}

func (this *FileQueue) Send(msg msgque.Message) error {
	return this.MessageQueue.Send(msg)
}

func (this *FileQueue) Close() {

}

func (this *FileQueue) fileFanout(ticket interface{}, msg msgque.Message) {
	switch msg.Type() {
	case save_file_msg:
		this.saveFile(ticket.(*DirTicket), msg.(*SaveMsg))
	case find_file_msg:
		this.findFile(msg.(*FindMsg))
	case del_file_msg:
		this.delFile(msg.(*DelMsg))
	default:
		log.Println(comerr.ParamTypeInvalid)
	}
}

func (this *FileQueue) saveFile(ticket *DirTicket, msg *SaveMsg) {
	fid := this.idflk.NextBase64Id(base64.RawURLEncoding)
	filePath := fmt.Sprintf("%s/%s/%s.%s", this.topDir, ticket.Dir, fid, msg.Ext)
	err := ioutil.WriteFile(filePath, msg.Buf, msg.fileMode)
	if err != nil {
		msg.Put(&CallbackMsg{
			Status: filesys_failed,
			Msg:    err.Error(),
		})

		return
	}

	err = addFile(this.fqdb, &MFile{
		FId:  fid,
		DId:  ticket.Dir,
		Path: filePath,
	})
	if err != nil {
		msg.Put(&CallbackMsg{
			Status: filesys_failed,
			Msg:    err.Error(),
		})

		return
	}

	ticket.Capacity++

	msg.Put(&CallbackMsg{
		Status: filesys_success,
		Payload: &SaveCbMsg{
			FId:  fid,
			DId:  ticket.Dir,
			Path: filePath,
		},
	})
}

func (this *FileQueue) findFile(msg *FindMsg) {
	ms, err := findFiles(this.fqdb, "f_id='"+msg.FId+"'")
	if err != nil {
		msg.Put(&CallbackMsg{
			Status: filesys_failed,
			Msg:    err.Error(),
		})
	} else if len(ms) != 1 {
		msg.Put(&CallbackMsg{
			Status: filesys_failed,
			Msg:    comerr.NotFound.Error(),
		})
	} else {
		m := ms[0]
		if m.IsDir {
			this.Suspend()
			defer this.Resume()

			var (
				found    = false
				capacity = m.Capacity
			)
			this.Traverse(func(ticket interface{}) bool {
				if ticket.(*DirTicket).Dir == m.FId {
					found = true
					capacity = ticket.(DirTicket).Capacity
					this.Restore(ticket)
				}

				return found
			})
			msg.Put(&CallbackMsg{
				Status: filesys_success,
				Payload: &FindCbMsg{
					FId:      m.FId,
					DId:      m.DId,
					IsDir:    m.IsDir,
					Capacity: capacity,
					Path:     m.Path,
				},
			})
		} else {
			msg.Put(&CallbackMsg{
				Status: filesys_success,
				Payload: &FindCbMsg{
					FId:      m.FId,
					DId:      m.DId,
					IsDir:    m.IsDir,
					Capacity: m.Capacity,
					Path:     m.Path,
				},
			})
		}
	}
}

func (this *FileQueue) delFile(msg *DelMsg) {
	ms, err := findFiles(this.fqdb, "f_id='"+msg.FId+"'")
	if err != nil {
		msg.Put(&CallbackMsg{
			Status: filesys_failed,
			Msg:    err.Error(),
		})
	} else if len(ms) != 1 {
		msg.Put(&CallbackMsg{
			Status: filesys_failed,
			Msg:    comerr.NotFound.Error(),
		})
	} else {
		this.Suspend()
		defer this.Resume()

		m := ms[0]
		found := false
		this.Traverse(func(ticket interface{}) bool {
			if ticket.(*DirTicket).Dir == m.FId {
				if err = deleteFile(this.fqdb, fmt.Sprintf("f_id='%s' and d_id='%s'", m.FId, m.FId)); err != nil {
					msg.Put(&CallbackMsg{
						Status: filesys_failed,
						Msg:    err.Error(),
					})
				}
				if err = os.RemoveAll(m.Path); err != nil {
					msg.Put(&CallbackMsg{
						Status: filesys_failed,
						Msg:    err.Error(),
					})
				}
				this.Restore(this.Generate())
				found = true
			} else if ticket.(*DirTicket).Dir == m.DId {
				ticket.(*DirTicket).Capacity--
				if err = deleteFile(this.fqdb, "f_id='"+m.FId+"'"); err != nil {
					msg.Put(&CallbackMsg{
						Status: filesys_failed,
						Msg:    err.Error(),
					})
				}
				if err = os.Remove(m.Path); err != nil {
					msg.Put(&CallbackMsg{
						Status: filesys_failed,
						Msg:    err.Error(),
					})
				}
				this.Restore(ticket)
				found = true
			}

			return found
		})

		if !
	}
}
