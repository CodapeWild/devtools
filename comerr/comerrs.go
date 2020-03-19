package comerr

import "errors"

var (
	ParamInvalid      = errors.New("invalid parameter")
	ParamTypeInvalid  = errors.New("invalid type of parameter")
	DataConvertFailed = errors.New("convert data failed")
	DataOutOfRange    = errors.New("data out of range")
	NullAddress       = errors.New("null address")
	EmptyValue        = errors.New("empty vlaue")
	IncompleteModule  = errors.New("incomplete module")
	ConnectFailed     = errors.New("connecting failed")
	NotFound          = errors.New("not found")
	Overtime          = errors.New("process overtime")
)
