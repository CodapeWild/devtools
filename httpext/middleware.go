package httpext

import (
	"context"
	"devtools/session"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

const (
	Cookie_Header_Key = "Authorization"
)

const (
	Id_Context_Key                  = "ID_CTX"
	MIME_Context_Key                = "MIME_CTX"
	Media_Context_Key               = "Media_CTX"
	MultipartFile_Context_Key       = "MultipartFILE_CTX"
	MultipartFileHeader_Context_Key = "MultipartFileHeader_CTX"
)

func PostOnly(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			respw.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(respw, req)
		}
	})
}

func GetOnly(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			respw.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(respw, req)
		}
	})
}

func AfterLogin(handler http.Handler, keeper *session.Keeper) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		if cookie := req.Header.Get(Cookie_Header_Key); keeper.GetSessionTTL(cookie) <= 0 {
			NewStdResp(StateDataExpired, nil).WriteJson(respw)
		} else {
			if idstr, err := keeper.CookieValue(cookie); err != nil {
				log.Println(err.Error())
				NewStdResp(StateProcessError, nil).WriteJson(respw)
			} else {
				handler.ServeHTTP(respw, req.WithContext(context.WithValue(req.Context(), Id_Context_Key, idstr)))
			}
		}
	})
}

// map[media]struct{AllowTypes: [aac, mp4, x-flv, jpeg...], MaxSize: kb}
type MIMEValidator map[string]struct {
	AllowedTypes []string `json:"allowed_types"`
	MaxSize      int64    `json:"max_size"`
}

func CheckMultiFile(handler http.Handler, formFileKey string, mval MIMEValidator) http.Handler {
	return http.HandlerFunc(func(respw http.ResponseWriter, req *http.Request) {
		mf, mfh, err := req.FormFile(formFileKey)
		if err != nil {
			log.Println(err.Error())
			NewStdResp(StateProcessError, nil).WriteJson(respw)

			return
		}

		// sniffer := make([]byte, 512)
		// if _, err = mf.Read(sniffer); err != nil {
		// 	log.Println(err.Error())
		// 	NewStdResp(StateProcessError, nil).WriteJson(respw)

		// 	return
		// }
		// iana = http.DetectContentType(sniffer)

		var (
			iana  string
			media string
			ext   string
		)
		iana, ext, err = mimetype.DetectReader(mf)
		if err != nil {
			log.Println(err.Error())
			NewStdResp(StateProcessError, nil).WriteJson(respw)

			return
		}
		if iana == "application/octet-stream" {
		} else {
			media = iana[:strings.Index(iana, "/")]
		}

		if conds, ok := mval[media]; !ok {
			NewStdResp(StateDataMediaInvalid, nil).WriteJson(respw)
		} else {
			// // check max size
			// req.Body = http.MaxBytesReader(respw, req.Body, 5800)
			// if err := req.ParseMultipartForm(5800); err != nil {
			// 	log.Println(err.Error())
			// 	NewStdResp(StateDataSizeInvalid, nil).WriteJson(respw)

			// 	return
			// }
			if mfh.Size > conds.MaxSize<<10 {
				NewStdResp(StateDataSizeInvalid, nil).WriteJson(respw)
			} else {
				for _, v := range conds.AllowedTypes {
					if ext == v {
						mf.Seek(0, io.SeekStart)
						req = req.WithContext(context.WithValue(req.Context(), MIME_Context_Key, iana))
						req = req.WithContext(context.WithValue(req.Context(), Media_Context_Key, media))
						req = req.WithContext(context.WithValue(req.Context(), MultipartFile_Context_Key, mf))
						req = req.WithContext(context.WithValue(req.Context(), MultipartFileHeader_Context_Key, mfh))
						handler.ServeHTTP(respw, req)

						return
					}
				}
				NewStdResp(StateDataTypeInvalid, nil).WriteJson(respw)
			}
		}
	})
}
