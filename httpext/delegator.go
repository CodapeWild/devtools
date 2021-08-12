package httpext

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/CodapeWild/devtools/comerr"
)

var defClient = http.Client{Timeout: time.Second}

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

// post and receive json data
func PostJson(rawurl string, req, resp interface{}) (status int, err error) {
	if resp == nil || reflect.TypeOf(resp).Kind() != reflect.Ptr {
		return http.StatusBadRequest, comerr.ErrParamInvalid
	}

	var u *url.URL
	if u, err = url.Parse(rawurl); err != nil {
		return
	}

	var buf []byte
	if buf, err = json.Marshal(req); err != nil {
		return
	}

	var r *http.Request
	if r, err = http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(buf)); err != nil {
		return
	}
	r.Header.Set("Content-Type", "application/json")

	var p *http.Response
	if p, err = defClient.Do(r); err != nil {
		return 0, err
	}

	status = p.StatusCode
	if status != http.StatusOK {
		return status, errors.New(p.Status)
	}

	if buf, err = ioutil.ReadAll(p.Body); err != nil {
		return
	}
	defer p.Body.Close()

	return status, json.Unmarshal(buf, resp)
}

func HandleJsonResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return buf, nil
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
