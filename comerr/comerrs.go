package comerr

import "errors"

// progress
var (
	ParamInvalid      = errors.New("invalid parameter")
	ParamTypeInvalid  = errors.New("invalid type of parameter")
	DataConvertFailed = errors.New("convert data failed")
	DataOutOfRange    = errors.New("data out of range")
	NilPointer        = errors.New("nil pointer")
	EmptyValue        = errors.New("empty vlaue")
	IncompleteModule  = errors.New("incomplete module")
)

// database
var (
	DbConnFailed  = errors.New("connect to database failed")
	DbNameInvalid = errors.New("invalid database name")
	NoEntryFound  = errors.New("no rows found")
)
