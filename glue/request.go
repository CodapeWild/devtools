package glue

import (
	"devtools/comerr"
	"encoding/json"
	"errors"
	"reflect"
)

var (
	serverFieldAbsent = errors.New("remote server field absent")
	clientFiledAbsent = errors.New("local client field absent")
)

type CallbackFunc func(resp *Response)

type Request struct {
	Headers    map[string][]string
	Payload    []byte
	IsCallback uint8
	callback   CallbackFunc
	timeoutSec int
	err        error
}

func NewRequest() *Request {
	return &Request{Headers: make(map[string][]string)}
}

func (this *Request) SetRemote(server, routine string) *Request {
	if this.err == nil {
		var key string
		key, this.err = formatServerKey(server)
		this.Headers["remote"] = []string{key, routine}
	}

	return this
}

func (this *Request) SetLocal(server, client string) *Request {
	if this.err == nil {
		var key string
		key, this.err = formatClientKey(server, client)
		this.Headers["local"] = []string{key}
	}

	return this
}

func (this *Request) SetPayload(payload interface{}) *Request {
	if this.err == nil {
		if payload == nil {
			this.err = comerr.ParamInvalid

			return this
		}
		t := reflect.TypeOf(payload)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if k := t.Kind(); k != reflect.Struct && k != reflect.Map {
			this.err = comerr.ParamTypeInvalid

			return this
		}

		this.Payload, this.err = json.Marshal(payload)
	}

	return this
}

func (this *Request) SetCallback(callback CallbackFunc) *Request {
	this.IsCallback = 1
	this.callback = callback

	return this
}

func (this *Request) SetTimeout(sec int) *Request {
	this.timeoutSec = sec

	return this
}

func (this *Request) GetServerKey() (key string, err error) {
	if len(this.Headers["remote"]) != 2 {
		return "", serverFieldAbsent
	}

	return this.Headers["remote"][0], nil
}

func (this *Request) GetRemote() (server, routine string, err error) {
	if len(this.Headers["remote"]) != 2 {
		return "", "", serverFieldAbsent
	}

	var key string
	key, routine = this.Headers["remote"][0], this.Headers["remote"][1]
	server, err = parseKey(key)

	return
}

func (this *Request) GetClientKey() (key string, err error) {
	if len(this.Headers["local"]) != 1 {
		return "", clientFiledAbsent
	}

	return this.Headers["local"][0], nil
}

func (this *Request) GetLocal() (client string, err error) {
	if len(this.Headers["local"]) != 1 {
		return "", clientFiledAbsent
	}

	client, err = parseKey(this.Headers["local"][0])

	return
}

func (this *Request) GetPayload(payload interface{}) error {
	if payload == nil || reflect.TypeOf(payload).Kind() != reflect.Ptr {
		return comerr.ParamTypeInvalid
	}

	return json.Unmarshal(this.Payload, payload)
}

func (this *Request) GetCallback() CallbackFunc {
	return this.callback
}

func (this *Request) GetTimeout() int {
	return this.timeoutSec
}
