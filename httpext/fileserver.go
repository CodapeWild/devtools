package httpext

import (
	"net/http"
)

type FileServer struct {
	http.ServeMux
	http.FileSystem
	root    string
	listDir bool
}

type Option func(fsrv *FileServer)

func ConfigRootDir(root string) Option {
	return func(fsrv *FileServer) {
		fsrv.root = root
	}
}

func ConfigIsListDir(list bool) Option {
	return func(fsrv *FileServer) {
		fsrv.listDir = list
	}
}

func RegisterUploadPattern(path string) Option {
	return func(fsrv *FileServer) {
		fsrv.Handle(path, http.HandlerFunc(fsrv.Upload))
	}
}

func RegisterDownloadPattern(path string) Option {
	return func(fsrv *FileServer) {
		fsrv.Handle(path, http.HandlerFunc(fsrv.Download))
	}
}

func RegisterOpenPattern(path string) Option {
	return func(fsrv *FileServer) {
		fsrv.Handle(path, http.StripPrefix(path, http.FileServer(fsrv)))
	}
}

func NewFileServer(opts ...Option) *FileServer {
	fsrv := &FileServer{}
	for _, opt := range opts {
		opt(fsrv)
	}

	return fsrv
}

func (this *FileServer) Upload(resp http.ResponseWriter, req *http.Request) {

}

func (this *FileServer) Download(resp http.ResponseWriter, req *http.Request) {

}

func (this *FileServer) Open(filePath string) (http.File, error) {

}
