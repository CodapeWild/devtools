package comerr

import "errors"

var (
	ParamInvalid         = errors.New("invalid parameter")
	TypeInvalid          = errors.New("invalid type")
	IndexOutOfRange      = errors.New("index out of range")
	DataConvertFailed    = errors.New("convert data failed")
	DataOutOfRange       = errors.New("data out of range")
	EmptyData            = errors.New("empty data")
	NilAddress           = errors.New("nil address")
	IncompleteModule     = errors.New("incomplete module")
	ConnectFailed        = errors.New("connecting failed")
	NotFound             = errors.New("not found")
	Overtime             = errors.New("process overtime")
	ProcessFailed        = errors.New("process failed")
	ChannelClosed        = errors.New("channel closed")
	UnrecognizedProtocol = errors.New("unrecognized protocol")
	HostLookupFailed     = errors.New("host lookup failed")
)
