package fileque

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type FileSystem struct {
	db        *sql.DB
	listFiles bool
}

func NewFileSystem(dbPath string, listFiles bool) (*FileSystem, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &FileSystem{
		db:        db,
		listFiles: listFiles,
	}, nil
}

func (this *FileSystem) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(path.Clean(name), "/")
	if len(name) != 11 {
		return nil, os.ErrInvalid
	}

	ms, err := findFiles(this.db, "f_id='"+name+"'")
	if err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	}
	if len(ms) != 1 {
		return nil, os.ErrNotExist
	}

	m := ms[0]
	if m.IsDir && !this.listFiles {
		return nil, os.ErrPermission
	}

	return http.Dir(".").Open(m.Path)
}
