package httpext

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/CodapeWild/devtools/comerr"
)

// RemoteAddr get remote ip who send the request and try to get the authentic ip:port
// hidden behind behind proxy.
func RemoteAddr(req *http.Request) (ip, port string) {
BREAKPOINT:
	for _, h := range []string{"x-forwarded-for", "X-FORWARDED-FOR", "X-Forwarded-For", "x-real-ip", "X-REAL-IP", "X-Real-Ip", "proxy-client-ip", "PROXY-CLIENT-IP", "Proxy-Client-Ip"} {
		addrs := strings.Split(req.Header.Get(h), ",")
		for _, addr := range addrs {
			if ip, port, _ = net.SplitHostPort(addr); ip == "" {
				continue
			}
			break BREAKPOINT
		}
	}
	if ip == "" {
		ip, port, _ = net.SplitHostPort(req.RemoteAddr)
	}

	return
}

func UnmarshalJsonReq(req *http.Request, out interface{}) error {
	if rv := reflect.ValueOf(out); rv.IsNil() || rv.Kind() != reflect.Ptr {
		return comerr.ErrParamInvalid
	}

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil
	}

	return json.Unmarshal(body, out)
}
