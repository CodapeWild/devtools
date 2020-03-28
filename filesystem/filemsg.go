package filesystem

import (
	"devtools/msgque"
)

const (
	Save_File int = iota + 1
	Del_File
)

const (
	FileSys_Proc_Success int = iota + 1
	FileSys_Proc_Failed
)

type SaveFileMsg struct {
	MsgId string
	Name  string
	Buf   []byte
	Size  int64
	Media MediaType
	Span  int64
	State int
	msgque.Callback
}

func (this *SaveFileMsg) Id() interface{} {
	return this.MsgId
}

func (this *SaveFileMsg) Type() interface{} {
	return Save_File
}

type DeleteFileMsg struct {
	MsgId   string
	Code    string
	DirCode string
	msgque.Callback
}

func (this *DeleteFileMsg) Id() interface{} {
	return this.Id
}

func (this *DeleteFileMsg) Type() interface{} {
	return Del_File
}

type FileSysCallbackMsg struct {
	MsgId string
	State int
	Err   error
}
