package filesystem

import (
	"database/sql"
	"devtools/comerr"
	"devtools/file"
	"devtools/idflaker"
	"devtools/msgque"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	Def_Top_Dir                  = "./upload"
	Def_Dir_Mode     os.FileMode = 0744
	Def_File_Mode    os.FileMode = 0744
	Def_Max_Contains             = 3
	Def_FSDB_Path                = "./fs.db"
	def_fstab_file               = "tab_file"
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
	*msgque.MessageQueue
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

	fs.MessageQueue = msgque.NewMessageQueue(msgque.SetTicket(NewDirectoryQueue(fs.topDir, fs.dirMode, fs.idflk, fs.fsdb, 6)), msgque.SetQueueBuffer(6), msgque.SetQueueTimeout(time.Second))
	fs.MessageQueue.StartUp(fs.fileMsgFanout)

	return fs, nil
}

func (this *FileSystem) Open(filePath string) (http.File, error) {
	codes := strings.Split(strings.TrimRight(filePath, "/"), "/")
	if len(codes) > 2 || len(codes) < 1 {
		return nil, os.ErrInvalid
	}

	m, err := findMFile(this.fsdb, "path="+filePath)
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

func (this *FileSystem) Send(msg msgque.Message) {
	this.MessageQueue.Send(msg)
}

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
	code := this.idflk.NextBase64Id()
	dirCode := string(ticket.(DirectoryTicket))
	path := fmt.Sprintf("/%s/%s", code, dirCode)
	m := &MFile{
		Code:        code,
		DirCode:     dirCode,
		IsDirectory: false,
		Path:        path,
		FileMode:    this.fileMode,
		FileSize:    msg.Size,
		Media:       msg.Media,
		Span:        msg.Span,
		State:       File_Normal,
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
