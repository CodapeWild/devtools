package fileque

import (
	"devtools/msgque"
	"os"
)

const (
	save_file_msg int = iota + 1
	find_file_msg
	del_file_msg
)

type SaveMsg struct {
	Buf      []byte
	FileMode os.FileMode
	Ext      string
	msgque.Callback
}

func (this *SaveMsg) Id() interface{} {
	return nil
}

func (this *SaveMsg) Type() interface{} {
	return save_file_msg
}

func (this *SaveMsg) MustFetch() bool {
	return true
}

type FindMsg struct {
	FId string
	msgque.Callback
}

func (this *FindMsg) Id() interface{} {
	return nil
}

func (this *FindMsg) Type() interface{} {
	return find_file_msg
}

func (this *FindMsg) MustFetch() bool {
	return false
}

type DelMsg struct {
	FId string
	msgque.Callback
}

func (this *DelMsg) Id() interface{} {
	return nil
}

func (this *DelMsg) Type() interface{} {
	return del_file_msg
}

func (this *DelMsg) MustFetch() bool {
	return false
}

const (
	FileQue_Success int = iota + 1
	FileQue_Failed
)

type CallbackMsg struct {
	Status  int
	Msg     string
	Payload interface{}
}

type SaveCbMsg struct {
	FId  string
	DId  string
	Path string
}

type FindCbMsg struct {
	FId      string
	DId      string
	IsDir    bool
	Capacity int
	Path     string
}
