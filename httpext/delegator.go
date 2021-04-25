package httpext

import (
	"bytes"
	"devtools/comerr"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

var defClient = &http.Client{Timeout: time.Second}

// post and receive json data
func PostJson(rawurl string, req, resp interface{}) (status int, err error) {
	if resp == nil || reflect.TypeOf(resp).Kind() != reflect.Ptr {
		return http.StatusBadRequest, comerr.ParamInvalid
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