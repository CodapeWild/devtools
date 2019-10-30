package file

import (
	"errors"
	"os"
	"path"
)

var (
	DirNotExists   = errors.New("directory not exists")
	FileNotExists  = errors.New("file not exists")
	OpenFileFailed = errors.New("unable to open file")
)

func IsFileExists(filePath string) bool {
	finfo, err := os.Stat(filePath)

	return err == nil && !finfo.IsDir()
}

func IsDirExists(dir string) bool {
	finfo, err := os.Stat(dir)

	return err == nil && finfo.IsDir()
}

func CreateAnyway(filePath string, perm os.FileMode) (*os.File, error) {
	if dir := path.Dir(filePath); !IsDirExists(dir) {
		if err := os.MkdirAll(dir, perm); err != nil {
			return nil, err
		}
	}

	return os.Create(filePath)
}
