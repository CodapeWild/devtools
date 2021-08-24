package httpext

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/CodapeWild/devtools/file"
)

const (
	root_dir             = "./upload"
	max_mem_usage        = 10 << 20
	max_upload_file_size = 100 << 20
)

var (
	ErrFileSizeTooBig = errors.New("file size too large")
	ErrInvalidMIME    = errors.New("invalid file media type")
	ErrFileExists     = errors.New("file already exists")
)

var (
	defMakeDir = func(req *http.Request, root string, perm os.FileMode) (string, error) {
		dir := path.Join(root, time.Now().Local().Format("2006-01-02"))

		return dir, os.MkdirAll(dir, perm)
	}
	defSaveFile = func(dir, filename, extension string, perm os.FileMode, fh *multipart.FileHeader) error {
		filePath := fmt.Sprintf("%s.%s", filepath.Join(dir, filename), extension)
		if file.IsFileExists(filePath) {
			return ErrFileExists
		}

		temp, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, perm)
		if err != nil {
			return err
		}
		defer temp.Close()

		src, err := fh.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(temp, src)

		return err
	}
)

type FileServer struct {
	http.ServeMux
	http.FileSystem
	root              string
	perm              os.FileMode
	listDir           bool
	uploadFileHeader  string
	maxMemoryUsage    int64
	maxUploadFileSize int64
	validContentTypes map[string]string
	mkdir             func(req *http.Request, root string, perm os.FileMode) (string, error)
	saveFile          func(dir, filename, extension string, perm os.FileMode, fh *multipart.FileHeader) error
}

type FSrvOption func(fsrv *FileServer)

func ConfigRootDir(root string) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.root = root
	}
}

func ConfigFileMode(perm os.FileMode) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.perm = perm
	}
}

func ConfigListDir(list bool) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.listDir = list
	}
}

func ConfigUploadFileHeader(field string) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.uploadFileHeader = field
	}
}

func ConfigMaxMemoryUsage(max int64) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.maxMemoryUsage = max
	}
}

func ConfigMaxUploadFileSize(max int64) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.maxUploadFileSize = max
	}
}

func ConfigValidContentTypes(types map[string]string) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.validContentTypes = types
	}
}

func ConfigMakeDirFunc(mkdir func(req *http.Request, root string, perm os.FileMode) (string, error)) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.mkdir = mkdir
	}
}

func ConfigSaveFileFunc(savefile func(dir, filename, extension string, perm os.FileMode, fh *multipart.FileHeader) error) FSrvOption {
	return func(fsrv *FileServer) {
		fsrv.saveFile = savefile
	}
}

func RegisterUploadPattern(pattern string, middleware MiddlewareFunc) FSrvOption {
	return func(fsrv *FileServer) {
		if middleware == nil {
			fsrv.HandleFunc(pattern, fsrv.Upload)
		} else {
			fsrv.Handle(pattern, middleware(http.HandlerFunc(fsrv.Upload)))
		}
	}
}

func RegisterDownloadPattern(pattern string, middleware MiddlewareFunc) FSrvOption {
	return func(fsrv *FileServer) {
		if middleware == nil {
			fsrv.HandleFunc(pattern, fsrv.Download)
		} else {
			fsrv.Handle(pattern, middleware(http.HandlerFunc(fsrv.Download)))
		}
	}
}

func RegisterOpenPattern(pattern string) FSrvOption {
	return func(fsrv *FileServer) {
		if len(pattern) > 0 && pattern[len(pattern)-1] != '/' {
			pattern += "/"
		}
		fsrv.Handle(pattern, http.StripPrefix(pattern, http.FileServer(fsrv)))
	}
}

// NewFileServer will return valid file server, use this function to get file server object pointer
// to avoid uncomplete creation.
func NewFileServer(opts ...FSrvOption) *FileServer {
	fsrv := &FileServer{
		root:              root_dir,
		maxMemoryUsage:    max_mem_usage,
		maxUploadFileSize: max_upload_file_size,
		mkdir:             defMakeDir,
		saveFile:          defSaveFile,
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

	// read multiple files
	for k, fhs := range req.MultipartForm.File {
		log.Printf("read files from field key %s\n", k)
		for _, fh := range fhs {
			// check file size
			if fh.Size > this.maxUploadFileSize {
				log.Println(ErrFileSizeTooBig.Error())
				continue
			}
			// detect file media type
			ext, err := this.isValidMedia(fh)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			dir, err := this.mkdir(req, this.root, this.perm)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			if err = this.saveFile(dir, fh.Filename, ext, this.perm, fh); err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}

	resp.WriteHeader(http.StatusOK)
}

func (this *FileServer) Download(resp http.ResponseWriter, req *http.Request) {

}

func (this *FileServer) ReadDir(name string) ([]fs.DirEntry, error) {

}

func (this *FileServer) Open(filePath string) (http.File, error) {

}

// isValidMedia detect MIME type of file if valid the extension string of file will return
func (this *FileServer) isValidMedia(fh *multipart.FileHeader) (ext string, err error) {
	var mulf multipart.File
	mulf, err = fh.Open()
	if err != nil {
		return
	}
	defer mulf.Close()

	first512 := make([]byte, 512)
	if _, err = mulf.Read(first512); err != nil {
		return
	}

	mimestr := http.DetectContentType(first512)

	var ok bool
	if len(this.validContentTypes) != 0 {
		if ext, ok = this.validContentTypes[mimestr]; !ok {
			err = ErrInvalidMIME
		}
	} else {
		var exts []string
		if exts, err = mime.ExtensionsByType(mimestr); err == nil {
			ext = exts[0]
		}
	}

	return
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
