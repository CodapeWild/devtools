package httpext

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	upload_file_header = "upload_file_field"
)

const (
	root_dir             = "./upload"
	max_mem_usage        = 10 << 20
	max_upload_file_size = 100 << 20
)

var (
	ErrFileSizeTooBig = errors.New("file size too large")
	ErrInvalidMIME    = errors.New("invalid file type")
)

var (
	defMkDir = func(req *http.Request, root string) (string, error) {
		dir := path.Join(root, time.Now().Local().Format("2006-01-02"))

		return dir, os.MkdirAll(dir, 0755)
	}
	defWriteFile = func(dir, filename, extension string, buf []byte) error {
		return os.WriteFile(fmt.Sprintf("%s%c%s.%s", dir, os.PathSeparator, filename, extension), buf, 0755)
	}
)

type FileServer struct {
	http.ServeMux
	http.FileSystem
	root              string
	listDir           bool
	uploadFileHeader  string
	maxMemoryUsage    int64
	maxUploadFileSize int64
	validContentTypes map[string]string
	mkdir             func(req *http.Request, root string) (string, error)
	writefile         func(dir, filename, extension string, buf []byte) error
}

type Option func(fsrv *FileServer)

func ConfigRootDir(root string) Option {
	return func(fsrv *FileServer) {
		fsrv.root = root
	}
}

func ConfigListDir(list bool) Option {
	return func(fsrv *FileServer) {
		fsrv.listDir = list
	}
}

func ConfigUploadFileHeader(field string) Option {
	return func(fsrv *FileServer) {
		fsrv.uploadFileHeader = field
	}
}

func ConfigMaxMemoryUsage(max int64) Option {
	return func(fsrv *FileServer) {
		fsrv.maxMemoryUsage = max
	}
}

func ConfigMaxUploadFileSize(max int64) Option {
	return func(fsrv *FileServer) {
		fsrv.maxUploadFileSize = max
	}
}

func ConfigValidContentTypes(types map[string]string) Option {
	return func(fsrv *FileServer) {
		fsrv.validContentTypes = types
	}
}

func ConfigMkDirFunc(mkdir func(req *http.Request, root string) (string, error)) Option {
	return func(fsrv *FileServer) {
		fsrv.mkdir = mkdir
	}
}

func ConfigWriteFileFunc(writefile func(dir, filename, extension string, buf []byte) error) Option {
	return func(fsrv *FileServer) {
		fsrv.writefile = writefile
	}
}

func RegisterUploadPattern(pattern string, middleware MiddlewareFunc) Option {
	return func(fsrv *FileServer) {
		if middleware == nil {
			fsrv.HandleFunc(pattern, fsrv.Upload)
		} else {
			fsrv.Handle(pattern, middleware(http.HandlerFunc(fsrv.Upload)))
		}
	}
}

func RegisterDownloadPattern(pattern string, middleware MiddlewareFunc) Option {
	return func(fsrv *FileServer) {
		if middleware == nil {
			fsrv.HandleFunc(pattern, fsrv.Download)
		} else {
			fsrv.Handle(pattern, middleware(http.HandlerFunc(fsrv.Download)))
		}
	}
}

func RegisterOpenPattern(pattern string) Option {
	return func(fsrv *FileServer) {
		fsrv.Handle(pattern, http.StripPrefix(pattern, http.FileServer(fsrv)))
	}
}

func NewFileServer(opts ...Option) *FileServer {
	fsrv := &FileServer{
		root:              root_dir,
		uploadFileHeader:  upload_file_header,
		maxMemoryUsage:    max_mem_usage,
		maxUploadFileSize: max_upload_file_size,
		mkdir:             defMkDir,
		writefile:         defWriteFile,
	}
	for _, opt := range opts {
		opt(fsrv)
	}

	return fsrv
}

func (this *FileServer) Upload(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(this.maxMemoryUsage)
	if err != nil {
		log.Println(err.Error())
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	mulf, mulh, err := req.FormFile(this.uploadFileHeader)
	if err != nil {
		log.Println(err.Error())
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	if mulh.Size > this.maxUploadFileSize {
		log.Println(ErrFileSizeTooBig.Error())
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	first512 := make([]byte, 512)
	if _, err := mulf.Read(first512); err != nil {
		log.Println(err.Error())
		resp.WriteHeader(http.StatusBadRequest)

		return
	}
	mimestr := http.DetectContentType(first512)

	var ext string
	if len(this.validContentTypes) != 0 {
		var ok bool
		if ext, ok = this.validContentTypes[mimestr]; !ok {
			log.Println(ErrInvalidMIME.Error())
			resp.WriteHeader(http.StatusBadRequest)

			return
		}
	} else {
		if extensions, err := mime.ExtensionsByType(mimestr); err != nil {
			log.Println(err.Error())
			resp.WriteHeader(http.StatusBadRequest)

			return
		} else {
			ext = extensions[0]
		}
	}

	if this.mkdir == nil {
		this.mkdir = defMkDir
	}
	dir, err := this.mkdir(req, this.root)
	if err != nil {
		log.Println(err.Error())
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	buf, err := io.ReadAll(mulf)
	if err != nil {
		log.Println(err.Error())
		resp.WriteHeader(http.StatusBadRequest)

		return
	}

	if this.writefile == nil {
		this.writefile = defWriteFile
	}
	if err = this.writefile(dir, mulh.Filename, ext, buf); err != nil {
		log.Println(err.Error())
		resp.WriteHeader(http.StatusBadRequest)
	} else {
		resp.WriteHeader(http.StatusOK)
	}
}

func (this *FileServer) Download(resp http.ResponseWriter, req *http.Request) {

}

func (this *FileServer) Open(filePath string) (http.File, error) {

}

// func CheckMultiFile(handler http.Handler, formFileKey string, mval MIMEValidator) http.Handler {
// 	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
// 		mf, mfh, err := req.FormFile(formFileKey)
// 		if err != nil {
// 			log.Println(err.Error())
// 			NewStdResp(StateProcessError, nil).WriteJson(respw)

// 			return
// 		}

// 		// sniffer := make([]byte, 512)
// 		// if _, err = mf.Read(sniffer); err != nil {
// 		// 	log.Println(err.Error())
// 		// 	NewStdResp(StateProcessError, nil).WriteJson(respw)

// 		// 	return
// 		// }
// 		// iana = http.DetectContentType(sniffer)

// 		var (
// 			iana  string
// 			media string
// 			ext   string
// 		)
// 		iana, ext, err = mimetype.DetectReader(mf)
// 		if err != nil {
// 			log.Println(err.Error())
// 			NewStdResp(StateProcessError, nil).WriteJson(respw)

// 			return
// 		}
// 		if iana == "application/octet-stream" {
// 		} else {
// 			media = iana[:strings.Index(iana, "/")]
// 		}

// 		if conds, ok := mval[media]; !ok {
// 			NewStdResp(StateDataMediaInvalid, nil).WriteJson(respw)
// 		} else {
// 			// // check max size
// 			// req.Body = http.MaxBytesReader(respw, req.Body, 5800)
// 			// if err := req.ParseMultipartForm(5800); err != nil {
// 			// 	log.Println(err.Error())
// 			// 	NewStdResp(StateDataSizeInvalid, nil).WriteJson(respw)

// 			// 	return
// 			// }
// 			if mfh.Size > conds.MaxSize<<10 {
// 				NewStdResp(StateDataSizeInvalid, nil).WriteJson(respw)
// 			} else {
// 				for _, v := range conds.AllowedTypes {
// 					if ext == v {
// 						mf.Seek(0, io.SeekStart)
// 						req = req.WithContext(context.WithValue(req.Context(), IANA_ContextKey, iana))
// 						req = req.WithContext(context.WithValue(req.Context(), Media_ContextKey, media))
// 						req = req.WithContext(context.WithValue(req.Context(), MultipartFile_ContextKey, mf))
// 						req = req.WithContext(context.WithValue(req.Context(), MultipartFileHeader_ContextKey, mfh))
// 						handler.ServeHTTP(respw, req)

// 						return
// 					}
// 				}
// 				NewStdResp(StateDataTypeInvalid, nil).WriteJson(respw)
// 			}
// 		}
// 	})
// }
