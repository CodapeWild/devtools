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
	topDir      string
	maxContains int
	dirMode     os.FileMode
	idflk       *idflaker.IdFlaker
	fsdb        *sql.DB
	msgque.Ticket
}

func NewDirectoryQueue(topDir string, maxContains int, dirMode os.FileMode, idflk *idflaker.IdFlaker, fsdb *sql.DB, maxThrds int) *DirectoryQueue {
	if maxThrds <= 0 {
		maxThrds = def_max_threads
	}

	return &DirectoryQueue{
		topDir:      topDir,
		maxContains: maxContains,
		dirMode:     dirMode,
		idflk:       idflk,
		fsdb:        fsdb,
		Ticket:      msgque.NewTicketQueue(maxThrds),
	}
}

func (this *DirectoryQueue) Fill() {
	cs, err := findDirCodes(this.fsdb, fmt.Sprintf("is_dir=1 and contains<%d limit %d", this.maxContains, this.Threads()))
	if err != nil {
		panic(err)
	}

	for _, v := range cs {
		this.Restore(DirectoryTicket(v))
	}
	for i := 0; i < this.Threads()-len(cs); i++ {
		this.Restore(this.Generate())
	}
}

func (this *DirectoryQueue) Generate() interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	code := this.idflk.NextBase64Id(base64.RawURLEncoding)
	path := string(os.PathSeparator) + code
	if err := os.Mkdir(this.topDir+path, this.dirMode); err != nil {
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

func (this *DirectoryQueue) Restore(ticket interface{}) {
	dirTick, ok := ticket.(DirectoryTicket)
	if !ok {
		return
	}

	m, err := findMFile(this.fsdb, "code='"+string(dirTick)+"'")
	if err != nil {
		log.Println(err.Error())

		return
	}

	if m.Contains >= this.maxContains {
		this.Ticket.Restore(this.Generate())
	} else {
		this.Ticket.Restore(dirTick)
	}
}
