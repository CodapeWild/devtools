package directory

import (
	"devtools/comerr"
	"devtools/file"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	def_top_level                 = "./upload"
	def_dir_mode      os.FileMode = 0744
	def_file_mode     os.FileMode = 0744
	def_sub_max_files             = 300000
	def_copy_threads              = 30
)

const (
	WriteFailed = iota
	WriteSuccess
)

type AutoIncrementFileInfo struct {
	FilePath string
	Content  []byte
	Callback chan int
}

type AutoIncrementDirectory struct {
	TopLevelDir string
	DirMode     os.FileMode
	FileMode    os.FileMode
	SubMaxFiles int
	MaxThreads  int
	curSubIndex int
	curDirFiles int
	queue       chan *AutoIncrementFileInfo
	sync.Mutex
}

type AutoIncrementDirSetting func(*AutoIncrementDirectory)

func SetTopLevelDir(dir string) AutoIncrementDirSetting {
	return func(updir *AutoIncrementDirectory) {
		updir.TopLevelDir = dir
	}
}

func SetMode(dirMode, fileMode os.FileMode) AutoIncrementDirSetting {
	return func(updir *AutoIncrementDirectory) {
		updir.DirMode = dirMode
		updir.FileMode = fileMode
	}
}

func SetSubLevelMaxFiles(max int) AutoIncrementDirSetting {
	return func(updir *AutoIncrementDirectory) {
		updir.SubMaxFiles = max
	}
}

func SetMaxThreads(max int) AutoIncrementDirSetting {
	return func(updir *AutoIncrementDirectory) {
		updir.MaxThreads = max
	}
}

var (
	WriteFileFailed = errors.New("write file failed")
)

func NewAutoIncrementDirectory(settings ...AutoIncrementDirSetting) *AutoIncrementDirectory {
	updir := &AutoIncrementDirectory{
		TopLevelDir: def_top_level,
		DirMode:     def_dir_mode,
		FileMode:    def_file_mode,
		SubMaxFiles: def_sub_max_files,
		MaxThreads:  def_copy_threads,
	}

	for _, v := range settings {
		v(updir)
	}

	return updir
}

func (this *AutoIncrementDirectory) Init() error {
	fis, err := ioutil.ReadDir(this.TopLevelDir)
	if c := len(fis); c == 0 {
		if err = os.MkdirAll(fmt.Sprintf("%s/%d", this.TopLevelDir, this.curSubIndex), this.DirMode); err != nil {
			return err
		}
	} else {
		this.curSubIndex = c - 1
		if fis, err = ioutil.ReadDir(fmt.Sprintf("%s/%d", this.TopLevelDir, this.curSubIndex)); err != nil {
			return err
		}
		this.curDirFiles = len(fis)
	}

	this.queue = make(chan *AutoIncrementFileInfo)
	for i := 0; i < this.MaxThreads; i++ {
		go func() {
			for v := range this.queue {
				go writeFileRoutine(v, this.FileMode)
			}
		}()
	}

	return nil
}

// only pluse the file counter by 1 or -1
func (this *AutoIncrementDirectory) AddFileCounter(i int) (curDir string, err error) {
	if i != 1 && i != -1 {
		return "", comerr.ParamInvalid
	}

	this.Lock()
	defer this.Unlock()

	c := this.curDirFiles + i
	if c > this.SubMaxFiles {
		if err = this.addSubDir(this.curSubIndex + 1); err != nil {
			return
		}
	} else if c < 0 {
		this.curSubIndex -= 1
		this.curDirFiles = this.SubMaxFiles - 1
	}

	curDir = fmt.Sprintf("%s/%d", this.TopLevelDir, this.curSubIndex)
	this.curDirFiles += 1

	return
}

// create empty file as placeholder, call WriteToFile to wirte content to file
func (this *AutoIncrementDirectory) AddEmptyFile(name string) (filePath string, err error) {
	if name == "" {
		return "", comerr.ParamInvalid
	}

	this.Lock()
	defer this.Unlock()

	if c := this.curDirFiles + 1; c > this.SubMaxFiles {
		if err = this.addSubDir(this.curSubIndex + 1); err != nil {
			return
		}
	}

	filePath = fmt.Sprintf("%s/%d/%s", this.TopLevelDir, this.curSubIndex, name)
	if !file.IsFileExists(filePath) {
		if f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, this.FileMode); err != nil {
			return "", err
		} else {
			f.Close()
			this.curDirFiles++
		}
	}

	return filePath, nil
}

func (this *AutoIncrementDirectory) addSubDir(index int) error {
	curDir := fmt.Sprintf("%s/%d", this.TopLevelDir, index)
	if !file.IsDirExists(curDir) {
		if err := os.MkdirAll(curDir, this.DirMode); err != nil {
			return err
		}
		this.curSubIndex = index
		this.curDirFiles = 0
	}

	return nil
}

func (this *AutoIncrementDirectory) WriteToFile(finfo *AutoIncrementFileInfo) {
	if this.queue != nil && finfo != nil {
		go func() { this.queue <- finfo }()
	}
}

func writeFileRoutine(finfo *AutoIncrementFileInfo, fmode os.FileMode) {
	defer func() {
		if e := recover(); e != nil {
			if finfo.Callback != nil {
				finfo.Callback <- WriteFailed
			}
			log.Println(e)
		}
	}()

	f, err := os.OpenFile(finfo.FilePath, os.O_RDWR|os.O_TRUNC, fmode)
	if err != nil {
		log.Panic(err.Error())
	}
	if _, err = f.Write(finfo.Content); err != nil {
		log.Panic(err.Error())
	}
	if err = f.Close(); err != nil {
		log.Panic(err.Error())
	}
	if finfo.Callback != nil {
		finfo.Callback <- WriteSuccess
	}
}
