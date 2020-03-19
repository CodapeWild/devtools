package filesystem

import (
	"database/sql"
	"devtools/file"
	"devtools/idflaker"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Def_Top_Dir                  = "./upload"
	Def_Dir_Mode     os.FileMode = 0744
	Def_File_Mode    os.FileMode = 0744
	Def_Max_Contains             = 3
	Def_FSDB_Path                = "./fs.db"
	def_fstab_file               = "tab_file"
	def_chan_size                = 6
	def_timeout                  = 3
)

type File struct {
	http.File
	fsdb                      *sql.DB
	showHidden, showForbidden bool
}

func (this File) Readdir(count int) ([]os.FileInfo, error) {
	finfo, err := this.Stat()
	if err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	}
	if !finfo.IsDir() {
		return nil, os.ErrInvalid
	}

	fstate := strconv.Itoa(File_Normal)
	if this.showHidden {
		fstate += "," + strconv.Itoa(File_Hidden)
	}
	if this.showForbidden {
		fstate += "." + strconv.Itoa(File_Forbidden)
	}

	return findMFiles(this.fsdb, fmt.Sprintf("dir_code=%s and state in(%s) order by created desc group by is_dir", finfo.Name(), fstate))
}

type FileSystem struct {
	http.Dir
	topDir                               string
	listFiles, showHidden, showForbidden bool
	maxContains                          int
	dirMode, fileMode                    os.FileMode
	fsdbPath                             string
	fsdb                                 *sql.DB
	idflk                                *idflaker.IdFlaker
}

type FileSystemSetting func(fs *FileSystem)

func SetAccessibility(listFiles, showHidden, showForbidden bool) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.listFiles = listFiles
		fs.showHidden = showHidden
		fs.showForbidden = showForbidden
	}
}

func SetDirProperties(topDir string, maxContains int, dirMode, fileMode os.FileMode) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.topDir = topDir
		fs.maxContains = maxContains
		fs.dirMode = dirMode
		fs.fileMode = fileMode
	}
}

func SetFSDBPath(dbPath string) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.fsdbPath = dbPath
	}
}

func NewFileSystem(opt ...FileSystemSetting) (*FileSystem, error) {
	fs := &FileSystem{
		listFiles:     true,
		showHidden:    false,
		showForbidden: false,
		topDir:        Def_Top_Dir,
		maxContains:   Def_Max_Contains,
		dirMode:       Def_Dir_Mode,
		fileMode:      Def_File_Mode,
		fsdbPath:      Def_FSDB_Path,
	}
	for _, v := range opt {
		v(fs)
	}

	if !file.IsDirExists(fs.topDir) {
		if err := os.MkdirAll(fs.topDir, fs.dirMode); err != nil {
			return nil, err
		}
	}
	fs.Dir = http.Dir(fs.topDir)

	db, err := sql.Open("sqlite3", fs.fsdbPath)
	if err != nil {
		return nil, err
	}
	if err = createTable(db); err != nil {
		return nil, err
	}
	fs.fsdb = db

	if fs.idflk, err = idflaker.NewIdFlaker(1); err != nil {
		return nil, err
	}

	return fs, nil
}

func (this *FileSystem) Open(filePath string) (http.File, error) {
	codes := strings.Split(strings.TrimRight(filePath, "/"), "/")
	if len(codes) > 2 || len(codes) < 1 {
		return nil, os.ErrInvalid
	}

	m, err := findOneMFile(this.fsdb, "path="+filePath)
	if err != nil {
		log.Println(err.Error())

		return nil, os.ErrNotExist
	}

	if (m.IsDirectory && !this.listFiles) || (m.State == File_Hidden && !this.showHidden) || (m.State == File_Forbidden && !this.showForbidden) {
		return nil, os.ErrPermission
	}

	if f, err := this.Dir.Open(filePath); err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	} else {
		return &File{File: f, fsdb: this.fsdb, showHidden: this.showHidden, showForbidden: this.showForbidden}, nil
	}
}
