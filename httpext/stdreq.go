package httpext

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/CodapeWild/devtools/comerr"
)

func RemoteIp(req *http.Request) (ip, port string) {
	var err error
BREAKPOINT:
	for _, h := range []string{"x-forwarded-for", "x-real-ip", "proxy-client-ip"} {
		addrs := strings.Split(req.Header.Get(h), ",")
		for _, addr := range addrs {
			if ip, port, err = net.SplitHostPort(addr); err != nil || ip == "" {
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

// parse json data from request
func ReadJson(req *http.Request, param interface{}) error {
	if rv := reflect.ValueOf(param); rv.Kind() != reflect.Ptr || rv.IsNil() {
		return comerr.ErrParamInvalid
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	return json.Unmarshal(buf, param)
}
