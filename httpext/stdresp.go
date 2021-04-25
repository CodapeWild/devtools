package httpext

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
)

var (
	StateSuccess = &StdStatus{Status: 1000, Msg: "success"}
)

var (
	StateServiceAccessBlocked = &StdStatus{Status: 2000, Msg: "service access blocked"}
	StateParseParamFailed     = &StdStatus{Status: 2001, Msg: "parse parameter failed"}
	StateParamInvalid         = &StdStatus{Status: 2002, Msg: "invalid parameter for request"}
	StateProcessTimeout       = &StdStatus{Status: 2003, Msg: "processing timeout"}
	StateProcessError         = &StdStatus{Status: 2004, Msg: "processing error"}
	StateDataNotFound         = &StdStatus{Status: 2005, Msg: "data can not find"}
	StateDataVerifyFailed     = &StdStatus{Status: 2006, Msg: "data verification failed"}
	StateDataAccessBlocked    = &StdStatus{Status: 2007, Msg: "data access blocked"}
	StateDataModifyForbidden  = &StdStatus{Status: 2008, Msg: "data modification forbidden"}
	StateDataExpired          = &StdStatus{Status: 2009, Msg: "data expired"}
	StateDataSizeInvalid      = &StdStatus{Status: 2010, Msg: "data size invalid"}
	StateDataTypeInvalid      = &StdStatus{Status: 2011, Msg: "data type invalid"}
	StateDataMediaInvalid     = &StdStatus{Status: 2012, Msg: "data media invalid"}
)

type StdStatus struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type StdResp interface {
	Encode() ([]byte, error)
	Decode(buf []byte) error
	WriteTo(respw http.ResponseWriter) (int, error)
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

func (this *JsonResp) WriteTo(respw http.ResponseWriter) (int, error) {
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
	return gob.NewDecoder(bytes.NewReader(buf)).Decode(this)
}

func (this *GobResp) WriteTo(respw http.ResponseWriter) (int, error) {
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
