package comerr

import "errors"

var (
	ParamInvalid      = errors.New("invalid parameter")
	ParamTypeInvalid  = errors.New("invalid type of parameter")
	IndexOutOfRange   = errors.New("index out of range")
	DataConvertFailed = errors.New("convert data failed")
	DataOutOfRange    = errors.New("data out of range")
	EmptyData         = errors.New("empty data")
	NullAddress       = errors.New("null address")
	IncompleteModule  = errors.New("incomplete module")
	ConnectFailed     = errors.New("connecting failed")
	NotFound          = errors.New("not found")
	Overtime          = errors.New("process overtime")
	ProcessFailed     = errors.New("process failed")
)
