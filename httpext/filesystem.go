package httpext

import (
	"devtools/directory"
	"net/http"
	"os"
	"strings"
)

type File struct {
	http.File
	listFiles   bool
	hideDotFile bool
}

func (this *File) Readdir(count int) (fis []os.FileInfo, err error) {
	if this.listFiles {
		var all []os.FileInfo
		all, err = this.File.Readdir(count)
		if this.hideDotFile {
			for _, v := range all {
				if ns := v.Name(); len(ns) > 0 && ns[0] != '.' {
					fis = append(fis, v)
				}
			}
		} else {
			fis = all
		}
	}

	return
}

type FileSystem struct {
	http.Dir
	listFiles   bool
	hideDotFile bool
	finder      directory.Finder
}

// finder could be nil, if so file system will be used
func NewFileSystem(dir string, listFiles, hideDotFile bool) *FileSystem {
	return &FileSystem{
		Dir:         http.Dir(dir),
		listFiles:   listFiles,
		hideDotFile: hideDotFile,
	}
}

func NewFinderSystem(finder directory.Finder) *FileSystem {
	return &FileSystem{finder: finder}
}

func (this *FileSystem) Open(name string) (http.File, error) {
	if this.finder == nil {
		if this.hideDotFile {
			for _, v := range strings.Split(name, "/") {
				if len(v) > 0 && v[0] == '.' {
					return nil, os.ErrPermission
				}
			}
		}
	} else {
		if mf, err := this.finder.Get(name[1:]); err != nil {
			return nil, os.ErrNotExist
		} else {
			name = mf.Path
		}
	}

	if f, err := this.Dir.Open(name); err != nil {
		return nil, err
	} else {
		return &File{f, this.listFiles, this.hideDotFile}, nil
	}
}
