package comerr

import "errors"

// code related
var (
	ErrNilPointer           = errors.New("nil pointer")
	ErrParamInvalid         = errors.New("invalid parameter")
	ErrTypeInvalid          = errors.New("invalid type")
	ErrAssertionFailed      = errors.New("type assertion failed")
	ErrIndexOutOfRange      = errors.New("index out of range")
	ErrChannelClosed        = errors.New("channel closed")
	ErrDataConvertFailed    = errors.New("convert data failed")
	ErrEmptyData            = errors.New("empty data")
	ErrProcessCanceled      = errors.New("process canceled")
	ErrProcessOvertime      = errors.New("process overtime")
	ErrProcessFailed        = errors.New("process failed")
	ErrUnrecognizedProtocol = errors.New("unrecognized protocol")
	ErrEncodeInvalid        = errors.New("invalid encoding")
	ErrDecodeFailed         = errors.New("decode failed")
	ErrNotFound             = errors.New("not found")
)

// file related
var (
	ErrFileNotExists = errors.New("file does not exist")
	ErrDirNotExists  = errors.New("directory does not exist")
)
