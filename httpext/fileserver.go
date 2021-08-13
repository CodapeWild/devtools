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
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {

	}
	http.DefaultTransport
	req.FormFile()
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
