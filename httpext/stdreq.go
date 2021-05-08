package httpext

import (
	"devtools/comerr"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
)

func RemoteIp(req *http.Request) (ip, port string) {
	var err error
	for _, h := range []string{"x-forwarded-for", "x-real-ip", "proxy-client-ip"} {
		addrs := strings.Split(req.Header.Get(h), ",")
		for i := len(addrs) - 1; i >= 0; i-- {
			if ip, port, err = net.SplitHostPort(addrs[i]); err != nil || ip == "" {
				continue
			} else {
				goto FOUND
			}
		}
	}
FOUND:
	if ip == "" {
		ip, port, _ = net.SplitHostPort(req.RemoteAddr)
	}

	return
}

// parse json data from request
func ReadJson(req *http.Request, param interface{}) error {
	if param == nil || reflect.TypeOf(param).Kind() != reflect.Ptr {
		return comerr.ErrParamInvalid
	}

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	return json.Unmarshal(buf, param)
}
