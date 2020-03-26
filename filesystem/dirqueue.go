package filesystem

import (
	"database/sql"
	"devtools/idflaker"
	"devtools/msgque"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

type DirectoryTicket string

const (
	def_max_threads = 6
)

type DirectoryQueue struct {
	*msgque.TicketQueue
	topDir      string
	maxContains int
	dirMode     os.FileMode
	idflk       *idflaker.IdFlaker
	fsdb        *sql.DB
}

func NewDirectoryQueue(topDir string, maxContains int, dirMode os.FileMode, idflk *idflaker.IdFlaker, fsdb *sql.DB, maxThrds int) *DirectoryQueue {
	if maxThrds <= 0 {
		maxThrds = def_max_threads
	}

	return &DirectoryQueue{
		TicketQueue: msgque.NewTicketQueue(maxThrds),
		topDir:      topDir,
		maxContains: maxContains,
		dirMode:     dirMode,
		idflk:       idflk,
		fsdb:        fsdb,
	}
}

func (this *DirectoryQueue) Fill() {
	cs, err := findDirCodes(this.fsdb, fmt.Sprintf("is_dir=1 and contains<%d limit %d", this.maxContains, this.Threads()))
	if err != nil {
		panic(err)
	}

	for _, v := range cs {
		this.Recede(DirectoryTicket(v))
	}
	for i := 0; i < this.Threads()-len(cs); i++ {
		this.Recede(this.Generate())
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
