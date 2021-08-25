package httpext

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/CodapeWild/devtools/comerr"
)

type StdResp interface {
	Encode() ([]byte, error)
	Decode(buf []byte) error
	Response(respw http.ResponseWriter) (int, error)
}

type StdStatus struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type JsonResp struct {
	*StdStatus
	Payload interface{} `json:"payload"`
}

func NewJsonResp(status *StdStatus, payload interface{}) *JsonResp {
	return &JsonResp{
		StdStatus: status,
		Payload:   payload,
	}
}

func (this *JsonResp) Encode() ([]byte, error) {
	return json.Marshal(this)
}

func (this *JsonResp) Decode(buf []byte) error {
	return json.Unmarshal(buf, this)
}

func (this *JsonResp) Response(respw http.ResponseWriter) (int, error) {
	var (
		buf []byte
		err error
		n   int
	)
	if buf, err = this.Encode(); err != nil {
		respw.WriteHeader(http.StatusInternalServerError)
	} else {
		respw.Header().Set("Content-Type", "application/json")

		n, err = respw.Write(buf)
	}

	return n, err
}

type GobResp struct {
	*StdStatus
	Payload interface{}
}

func NewGobResp(status *StdStatus, payload interface{}) *GobResp {
	return &GobResp{
		StdStatus: status,
		Payload:   payload,
	}
}

func (this *GobResp) Encode() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(this)

	return buf.Bytes(), err
}

func (this *GobResp) Decode(buf []byte) error {
	if this != nil {
		return gob.NewDecoder(bytes.NewReader(buf)).Decode(this)
	} else {
		return comerr.ErrNilPointer
	}
}

func (this *GobResp) Response(respw http.ResponseWriter) (int, error) {
	var (
		buf []byte
		err error
		n   int
	)
	if buf, err = this.Encode(); err != nil {
		respw.WriteHeader(http.StatusInternalServerError)
	} else {
		respw.Header().Set("Content-Type", "application/octet-stream")

		n, err = respw.Write(buf)
	}

	return n, err
}
