package comerr

import "errors"

var (
	NilPointer           = errors.New("nil pointer")
	ParamInvalid         = errors.New("invalid parameter")
	TypeInvalid          = errors.New("invalid type")
	AssertiionFailed     = errors.New("type assertion failed")
	IndexOutOfRange      = errors.New("index out of range")
	ChannelClosed        = errors.New("channel closed")
	DataConvertFailed    = errors.New("convert data failed")
	EmptyData            = errors.New("empty data")
	ProcessCanceled      = errors.New("process canceled")
	ProcessOvertime      = errors.New("process overtime")
	ProcessFailed        = errors.New("process failed")
	UnrecognizedProtocol = errors.New("unrecognized protocol")
	NotFound             = errors.New("not found")
)
