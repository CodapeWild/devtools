package filesystem

import (
	"devtools/msgque"
	"time"
)

type FileMsgType int

const (
	Save_File = iota + 1
	Del_File
)

type SaveFileMsg struct {
	MsgId  string
	Buf    []byte
	Size   int64
	Media  MediaType
	Span   int64
	CbChan chan *msgque.CallbackMsg
}

func (this *SaveFileMsg) Type() interface{} {
	return Save_File
}

func (this *SaveFileMsg) Id() interface{} {
	return this.MsgId
}

func (this *SaveFileMsg) Callback(cbMsg *msgque.CallbackMsg, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case this.CbChan <- cbMsg:
		return true
	}
}

type DelFileMsg struct {
	MsgId  string
	Code   string
	CbChan chan *msgque.CallbackMsg
}

func (this *DelFileMsg) Id() interface{} {
	return this.Id
}

func (this *DelFileMsg) Type() interface{} {
	return Del_File
}

func (this *DelFileMsg) Callback(cbMsg *msgque.CallbackMsg, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case this.CbChan <- cbMsg:
		return true
	}
}
