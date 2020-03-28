package filesystem

import (
	"database/sql"
	"devtools/comerr"
	"devtools/file"
	"devtools/idflaker"
	"devtools/msgque"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type File struct {
	http.File
	path                      string
	showHidden, showForbidden bool
	fsdb                      *sql.DB
}

func (this File) Readdir(count int) ([]os.FileInfo, error) {
	var where string
	if this.path == "/" {
		where += "is_dir=1"
	} else {
		where += "dir_code='" + this.path[1:] + "'"
	}

	fstate := strconv.Itoa(File_Normal)
	if this.showHidden {
		fstate += "," + strconv.Itoa(File_Hidden)
	}
	if this.showForbidden {
		fstate += "," + strconv.Itoa(File_Forbidden)
	}

	return findMFiles(this.fsdb, where+fmt.Sprintf(" and state in(%s) order by created desc limit %d\n", fstate, count))
}

type OpenTag int

const (
	Open_WithPath = iota + 1
	Open_WithCode
)

const (
	Def_Top_Dir                  = "./upload"
	Def_Dir_Mode     os.FileMode = 0744
	Def_File_Mode    os.FileMode = 0744
	Def_Max_Contains             = 3
	Def_Open_Tag                 = Open_WithPath
	Def_FSDB_Path                = "./fs.db"
	def_tab_mfile                = "tab_file"
)

type FileSystem struct {
	http.Dir
	topDir                               string
	maxContains                          int
	dirMode, fileMode                    os.FileMode
	openTag                              OpenTag
	listFiles, showHidden, showForbidden bool
	fsdbPath                             string
	fsdb                                 *sql.DB
	idflk                                *idflaker.IdFlaker
	*msgque.MessageQueue
}

type FileSystemSetting func(fs *FileSystem)

func SetDirProperties(topDir string, maxContains int, dirMode, fileMode os.FileMode) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.topDir = topDir
		fs.maxContains = maxContains
		fs.dirMode = dirMode
		fs.fileMode = fileMode
	}
}

func SetOpenTag(tag OpenTag) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.openTag = tag
	}
}

func SetAccessibility(listFiles, showHidden, showForbidden bool) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.listFiles = listFiles
		fs.showHidden = showHidden
		fs.showForbidden = showForbidden
	}
}

func SetFSDBPath(dbPath string) FileSystemSetting {
	return func(fs *FileSystem) {
		fs.fsdbPath = dbPath
	}
}

func NewFileSystem(opt ...FileSystemSetting) (*FileSystem, error) {
	fs := &FileSystem{
		topDir:        Def_Top_Dir,
		maxContains:   Def_Max_Contains,
		dirMode:       Def_Dir_Mode,
		fileMode:      Def_File_Mode,
		openTag:       Def_Open_Tag,
		listFiles:     true,
		showHidden:    false,
		showForbidden: false,
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

	fs.MessageQueue = msgque.NewMessageQueue(msgque.SetTicket(NewDirectoryQueue(fs.topDir, fs.maxContains, fs.dirMode, fs.idflk, fs.fsdb, 6)), msgque.SetQueueBuffer(6), msgque.SetQueueTimeout(time.Second))
	fs.MessageQueue.StartUp(fs.fileMsgFanout)

	return fs, nil
}

func (this *FileSystem) Open(s string) (http.File, error) {
	if len(strings.TrimSpace(s)) == 0 {
		return nil, os.ErrInvalid
	}

	switch this.openTag {
	case Open_WithPath:
		return this.openWithPath(s)
	case Open_WithCode:
		return this.openWithCode(s)
	default:
		return nil, os.ErrInvalid
	}
}

func (this *FileSystem) openWithPath(filePath string) (http.File, error) {
	if filePath == "/" {
		if this.listFiles {
			if f, err := this.Dir.Open(filePath); err != nil {
				log.Println(err.Error())

				return nil, os.ErrInvalid
			} else {
				return &File{File: f, path: filePath, showHidden: this.showHidden, showForbidden: this.showForbidden, fsdb: this.fsdb}, nil
			}
		} else {
			return nil, os.ErrPermission
		}
	}

	m, err := findMFile(this.fsdb, "path='"+filePath+"'")
	if err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	}

	if (m.IsDirectory && !this.listFiles) || (m.State == File_Hidden && !this.showHidden) || (m.State == File_Forbidden && !this.showForbidden) {
		return nil, os.ErrPermission
	}

	if f, err := this.Dir.Open(m.Path); err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	} else {
		return &File{File: f, path: m.Path, showHidden: this.showHidden, showForbidden: this.showForbidden, fsdb: this.fsdb}, nil
	}
}

func (this *FileSystem) openWithCode(code string) (http.File, error) {
	code = strings.Trim(code, "/")
	if len(code) != 11 {
		return nil, os.ErrInvalid
	}

	m, err := findMFile(this.fsdb, "code='"+code+"'")
	if err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	}

	if (m.IsDirectory && !this.listFiles) || (m.State == File_Hidden && !this.showHidden) || (m.State == File_Forbidden && !this.showForbidden) {
		return nil, os.ErrPermission
	}

	if f, err := this.Dir.Open(m.Path); err != nil {
		log.Println(err.Error())

		return nil, os.ErrInvalid
	} else {
		return &File{File: f, path: m.Path, showHidden: this.showHidden, showForbidden: this.showForbidden, fsdb: this.fsdb}, nil
	}
}

// func (this *FileSystem) Open(filePath string) (http.File, error) {
// 	log.Println("open file")
// 	codes := strings.Split(strings.TrimRight(filePath, "/"), "/")
// 	if len(codes) > 2 || len(codes) < 1 {
// 		return nil, os.ErrInvalid
// 	}

// 	m, err := findMFile(this.fsdb, "path='"+filePath+"'")
// 	if err != nil {
// 		log.Println(err.Error())

// 		return nil, os.ErrNotExist
// 	}

// 	if (m.IsDirectory && !this.listFiles) || (m.State == File_Hidden && !this.showHidden) || (m.State == File_Forbidden && !this.showForbidden) {
// 		return nil, os.ErrPermission
// 	}

// 	if f, err := this.Dir.Open(filePath); err != nil {
// 		log.Println(err.Error())

// 		return nil, os.ErrInvalid
// 	} else {
// 		return &File{File: f, fsdb: this.fsdb, showHidden: this.showHidden, showForbidden: this.showForbidden}, nil
// 	}
// }

func (this *FileSystem) fileMsgFanout(ticket interface{}, msg msgque.Message) {
	switch msg.Type() {
	case Save_File:
		this.saveFile(ticket, msg.(*SaveFileMsg))
	case Del_File:
		this.deleteFile(msg.(*DeleteFileMsg))
	default:
		log.Println(comerr.ParamInvalid.Error())
	}
}

func (this *FileSystem) saveFile(ticket interface{}, msg *SaveFileMsg) {
	code := this.idflk.NextBase64Id(base64.RawURLEncoding)
	dirCode := string(ticket.(DirectoryTicket))
	path := fmt.Sprintf("/%s/%s", dirCode, code)
	m := &MFile{
		Code:         code,
		DirCode:      dirCode,
		IsDirectory:  false,
		Path:         path,
		OriginalName: msg.Name,
		FileMode:     this.fileMode,
		FileSize:     msg.Size,
		Media:        msg.Media,
		Span:         msg.Span,
		State:        msg.State,
	}

	err := insertMFile(this.fsdb, m)
	if err != nil {
		if msg.CbChan != nil {
			callback(msg, File_Opt_Failed, err, 1)
		}

		return
	}

	f, err := os.OpenFile(this.topDir+path, os.O_CREATE|os.O_WRONLY, this.fileMode)
	if err != nil {
		if msg.CbChan != nil {
			callback(msg, File_Opt_Failed, err, 1)
		}

		return
	}
	_, err = f.Write(msg.Buf)
	if msg.CbChan != nil {
		if err != nil {
			callback(msg, File_Opt_Failed, err, 1)
		} else {
			callback(msg, File_Opt_Success, nil, 1)
		}
	}
}

func (this *FileSystem) deleteFile(msg *DeleteFileMsg) {
	err := deleteMFile(this.fsdb, msg.Code)
	if err != nil {
		if msg.CbChan != nil {
			callback(msg, File_Opt_Failed, err, 1)
		}

		return
	}

	var path string = "/" + msg.Code
	if msg.DirCode != "" {
		path = "/" + msg.DirCode + path
	}
	path = this.topDir + path
	err = os.RemoveAll(path)
	if msg.CbChan != nil {
		if err != nil {
			callback(msg, File_Opt_Failed, err, 1)
		} else {
			callback(msg, File_Opt_Success, nil, 1)
		}
	}
}

func callback(msg msgque.Message, state int, err error, timeout time.Duration) {
	msg.Callback(&FileCallbackMsg{
		MsgId: msg.Id().(string),
		State: state,
		Err:   err,
	}, 1)
}
