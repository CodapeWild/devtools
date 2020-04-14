package httpext

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type StaticFile struct {
	http.File
	showHidden bool
}

func NewStaticFile(f http.File, showHidden bool) *StaticFile {
	return &StaticFile{
		File:       f,
		showHidden: showHidden,
	}
}

func (this *StaticFile) Readdir(count int) ([]os.FileInfo, error) {
	all, err := this.File.Readdir(count)
	if err != nil {
		return nil, err
	}

	var finfos []os.FileInfo
	for _, v := range all {
		if strings.HasPrefix(v.Name(), ".") && !this.showHidden {
			continue
		}
		finfos = append(finfos, v)
	}

	return finfos, nil
}

type StaticSystem struct {
	http.Dir
	listFiles, showHidden bool
}

func NewStaticSystem(topDir string, listFiles, showHidden bool) *StaticSystem {
	return &StaticSystem{
		Dir:        http.Dir(topDir),
		listFiles:  listFiles,
		showHidden: showHidden,
	}
}

func (this *StaticSystem) Open(name string) (http.File, error) {
	finfo, err := os.Stat(path.Clean(string(this.Dir) + name))
	if err != nil {
		return nil, err
	}

	if (!this.listFiles && finfo.IsDir()) || (!this.showHidden && strings.HasPrefix(finfo.Name(), ".")) {
		return nil, os.ErrPermission
	}

	f, err := this.Dir.Open(name)
	if err != nil {
		return nil, err
	}

	return NewStaticFile(f, this.showHidden), nil
}
