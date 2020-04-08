package filesys

import (
	"database/sql"
	"devtools/idflaker"
	"devtools/msgque"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

type DirTicket struct {
	Dir      string
	Capacity int
}

type DirTicketQueue struct {
	*msgque.SimpleTicketQueue
	topDir      string
	dirMode     os.FileMode
	dirCapacity int
	fqdb        *sql.DB
	idflk       *idflaker.IdFlaker
}

func NewDirTicketQueue(maxThrds int, topDir string, dirMode os.FileMode, dirCapacity int, fqdb *sql.DB) (*DirTicketQueue, error) {
	idflk, err := idflaker.NewIdFlaker(2)
	if err != nil {
		return nil, err
	}

	return &DirTicketQueue{
		SimpleTicketQueue: msgque.NewSimpleTicketQueue(maxThrds),
		dirMode:           dirMode,
		dirCapacity:       dirCapacity,
		fqdb:              fqdb,
		idflk:             idflk,
	}, nil
}

func (this *DirTicketQueue) Fill() {
	ms, err := findFiles(this.fqdb, fmt.Sprintf("is_dir=1 and capacity<%d order by capacity asc limit %d", this.dirCapacity, this.MaxThreads()))
	log.Fatalln(err.Error())
	for _, v := range ms {
		this.Restore(&DirTicket{Dir: v.FId})
	}
	for i := len(ms); i < this.MaxThreads(); i++ {
		this.Restore(this.Generate())
	}
}

func (this *DirTicketQueue) Generate() interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	dir := this.idflk.NextBase64Id(base64.RawURLEncoding)
	path := this.topDir + "/" + dir
	err := os.Mkdir(path, this.dirMode)
	if err != nil {
		log.Panicln(err.Error())
	}

	err = addFile(this.fqdb, &MFile{
		FId:   dir,
		IsDir: true,
		Path:  path,
	})
	if err != nil {
		log.Panicln(err.Error())
	}

	return &DirTicket{Dir: dir}
}

func (this *DirTicketQueue) Restore(ticket interface{}) {
	dirTick := ticket.(*DirTicket)
	if dirTick.Capacity >= this.dirCapacity {
		if err := updateDirCap(this.fqdb, dirTick.Dir, dirTick.Capacity); err != nil {
			log.Println(err.Error())
		}

		ticket = this.Generate()
	}

	this.SimpleTicketQueue.Restore(ticket)
}
