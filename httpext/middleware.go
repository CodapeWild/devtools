package httpext

import (
	"devtools/session"
	"net/http"
)

const (
	Authorization                  = "Authorization"
	Id_ContextKey                  = "ID_CTX"
	MultipartFile_ContextKey       = "MultipartFILE_CTX"
	MultipartFileHeader_ContextKey = "MultipartFileHeader_CTX"
	Media_ContextKey               = "Media_CTX"
	IANA_ContextKey                = "MIME_CTX"
)

func GetOnly(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			respw.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(respw, req)
		}
	})
}

func PostOnly(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			respw.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(respw, req)
		}
	})
}

func AfterLogin(handler http.Handler, sessToken session.SessToken) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		if token := req.Header.Get(Authorization); sessToken.Verify(token) {
			handler.ServeHTTP(respw, req)
		} else {
			NewJsonResp(StateDataExpired, nil).Response(respw)
		}
	})
}

// // map[media]struct{AllowTypes: [aac, mp4, x-flv, jpeg...], MaxSize: kb}
// type MIMEValidator map[string]struct {
// 	AllowedTypes []string `json:"allowed_types"`
// 	MaxSize      int64    `json:"max_size"`
// }

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
