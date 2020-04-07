package filesys

import (
	"devtools/msgque"
)

const (
	save_file_msg int = iota + 1
	find_file_msg
	del_file_msg
)

type SaveMsg struct {
	Buf []byte
	*msgque.CallbackQueue
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
	*msgque.CallbackQueue
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
	*msgque.CallbackQueue
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
	filesys_success int = iota + 1
	filesys_failed
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
