package file

import (
	"devtools/comerr"
	"os"
	"path"
)

func IsFileExists(filePath string) bool {
	if finfo, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	} else {
		if md := finfo.Mode(); !md.IsRegular() {
			return false
		}
	}

	return true
}

func IsDirExists(dir string) bool {
	if finfo, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	} else {
		if md := finfo.Mode(); !md.IsDir() {
			return false
		}
	}

	return true
}

func CreateFileAnyway(filePath string, perm os.FileMode) (*os.File, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, comerr.ErrPathAlreadyExists
		}
	}

	if err = os.MkdirAll(path.Dir(filePath), perm); err != nil {
		return nil, err
	} else {
		return os.OpenFile(filePath, os.O_CREATE, perm)
	}
}
