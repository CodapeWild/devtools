package glue

import (
	"devtools/comerr"
	"devtools/db/redisdb"
	"encoding/json"
	"reflect"
)

type StatusCode int

var (
	OK     StatusCode = 100
	Busy   StatusCode = 200
	Exit   StatusCode = 300
	Failed StatusCode = 400
)

type RespStatus struct {
	Status StatusCode
	Msg    string
}

var (
	StatusOK     = &RespStatus{Status: OK, Msg: "OK"}
	StatusBusy   = &RespStatus{Status: Busy, Msg: "BUSY"}
	StatusExit   = &RespStatus{Status: Exit, Msg: "EXIT"}
	StatusFailed = &RespStatus{Status: Failed, Msg: "FAILED"}
)

type Response struct {
	req     *Request
	wrapper *redisdb.RedisWrapper
	Status  *RespStatus
	Payload []byte
}

func NewResponse(req *Request, wrapper *redisdb.RedisWrapper) *Response {
	return &Response{
		req:     req,
		wrapper: wrapper,
	}
}

func (this *Response) SetStatus(status *RespStatus) *Response {
	this.Status = status

	return this
}

func (this *Response) WriteBack(payload interface{}) error {
	if this.req.IsCallback == 0 {
		return nil
	}
	clientKey, err := this.req.GetClientKey()
	if err != nil {
		return err
	}
	if this.wrapper == nil {
		return comerr.ParamInvalid
	}

	if payload == nil {
		this.Payload = []byte{}
	} else {
		if this.Payload, err = json.Marshal(payload); err != nil {
			return err
		}
	}

	buf, err := json.Marshal(this)
	if err != nil {
		return err
	}
	_, err = this.wrapper.RPush(clientKey, buf)

	return err
}

func (this *Response) GetPayload(payload interface{}) error {
	if payload == nil {
		return comerr.ParamInvalid
	} else if reflect.TypeOf(payload).Kind() != reflect.Ptr {
		return comerr.ParamTypeInvalid
	}

	return json.Unmarshal(this.Payload, payload)
}
