package httpext

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// post and receive json data
func PostJson(rawurl string, req, resp interface{}) (status int, err error) {
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
	if p, err = (&http.Client{Timeout: 3 * time.Second}).Do(r); err != nil {
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

// parse json data from request
func ReadJson(req *http.Request, param interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	return json.Unmarshal(buf, param)
}
