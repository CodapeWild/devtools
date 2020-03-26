package filesystem

import (
	"database/sql"
	"devtools/idflaker"
	"devtools/msgque"
	"encoding/base64"
	"log"
	"os"
)

type DirectoryTicket string

const (
	def_max_threads = 6
)

type DirectoryQueue struct {
	*msgque.TicketQueue
	topDir  string
	dirMode os.FileMode
	idflk   *idflaker.IdFlaker
	fsdb    *sql.DB
}

func NewDirectoryQueue(topDir string, dirMode os.FileMode, idflk *idflaker.IdFlaker, fsdb *sql.DB, maxThrds int) *DirectoryQueue {
	if maxThrds <= 0 {
		maxThrds = def_max_threads
	}

	return &DirectoryQueue{
		TicketQueue: msgque.NewTicketQueue(maxThrds),
		topDir:      topDir,
		dirMode:     dirMode,
		idflk:       idflk,
		fsdb:        fsdb,
	}
}

func (this *DirectoryQueue) Generate() interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	code := this.idflk.NextBase64Id(base64.RawURLEncoding)
	path := this.topDir + string(os.PathSeparator) + code
	if err := os.MkdirAll(path, this.dirMode); err != nil {
		panic(err)
	}
	if err := insertMFile(this.fsdb, &MFile{
		Code:        code,
		IsDirectory: true,
		Path:        path,
		FileMode:    this.dirMode,
		State:       File_Normal,
	}); err != nil {
		panic(err)
	}

	return DirectoryTicket(code)
}

func (this *DirectoryQueue) Recede(ticket interface{}) {
	dirTick, ok := ticket.(DirectoryTicket)
	if !ok {
		return
	}

	m, err := findMFile(this.fsdb, "code='"+string(dirTick)+"'")
	if err != nil {
		log.Println(err.Error())

		return
	}

	if m.Contains >= Def_Max_Contains {
		this.TicketQueue.Recede(this.Generate())
	} else {
		this.TicketQueue.Recede(dirTick)
	}
}
