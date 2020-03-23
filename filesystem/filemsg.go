package filesystem

import (
	"time"
)

const (
	Save_File int = iota + 1
	Del_File
)

const (
	File_Opt_Success int = iota + 1
	File_Opt_Failed
)

type SaveFileMsg struct {
	MsgId  string
	Buf    []byte
	Size   int64
	Media  MediaType
	Span   int64
	State  int
	CbChan chan interface{}
}

func (this *SaveFileMsg) Id() interface{} {
	return this.MsgId
}

func (this *SaveFileMsg) Type() interface{} {
	return Save_File
}

func (this *SaveFileMsg) Callback(cbMsg interface{}, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case this.CbChan <- cbMsg:
		return true
	}
}

type DeleteFileMsg struct {
	MsgId   string
	Code    string
	DirCode string
	CbChan  chan interface{}
}

func (this *DeleteFileMsg) Id() interface{} {
	return this.Id
}

func (this *DeleteFileMsg) Type() interface{} {
	return Del_File
}

func (this *DeleteFileMsg) Callback(cbMsg interface{}, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case this.CbChan <- cbMsg:
		return true
	}
}

type FileCallbackMsg struct {
	MsgId string
	State int
	Err   error
}
