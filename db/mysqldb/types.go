package mysqldb

import (
	"reflect"
	"time"
)

var (
	_uint8      uint8
	_uint16     uint16
	_uint32     uint32
	_uint       uint
	_uint64     uint64
	_int8       int8
	_int16      int16
	_int32      int32
	_int        int
	_int64      int64
	_float32    float32
	_float64    float64
	_complex64  complex64
	_complex128 complex128
	_byte       byte
	_string     string
	_bool       bool
	_time       time.Time
)

var (
	GoUInt8      = reflect.TypeOf(_uint8)
	GoUInt16     = reflect.TypeOf(_uint16)
	GoUInt32     = reflect.TypeOf(_uint32)
	GoUInt       = reflect.TypeOf(_uint)
	GoUInt64     = reflect.TypeOf(_uint64)
	GoInt8       = reflect.TypeOf(_int8)
	GoInt16      = reflect.TypeOf(_int16)
	GoInt32      = reflect.TypeOf(_int32)
	GoInt        = reflect.TypeOf(_int)
	GoInt64      = reflect.TypeOf(_int64)
	GoFloat32    = reflect.TypeOf(_float32)
	GoFloat64    = reflect.TypeOf(_float64)
	GoComplex64  = reflect.TypeOf(_complex64)
	GoComplex128 = reflect.TypeOf(_complex128)
	GoByte       = reflect.TypeOf(_byte)
	GoString     = reflect.TypeOf(_string)
	GoBool       = reflect.TypeOf(_bool)
	GoTime       = reflect.TypeOf(_time)
)

var (
	GoUInt8Ptr      = reflect.PtrTo(GoUInt8)
	GoUInt16Ptr     = reflect.PtrTo(GoUInt16)
	GoUInt32Ptr     = reflect.PtrTo(GoUInt32)
	GoUIntPtr       = reflect.PtrTo(GoUInt)
	GoUInt64Ptr     = reflect.PtrTo(GoUInt64)
	GoInt8Ptr       = reflect.PtrTo(GoInt8)
	GoInt16Ptr      = reflect.PtrTo(GoInt16)
	GoInt32Ptr      = reflect.PtrTo(GoInt32)
	GoIntPtr        = reflect.PtrTo(GoInt)
	GoInt64Ptr      = reflect.PtrTo(GoInt64)
	GoFloat32Ptr    = reflect.PtrTo(GoFloat32)
	GoFloat64Ptr    = reflect.PtrTo(GoFloat64)
	GoComplex64Ptr  = reflect.PtrTo(GoComplex64)
	GoComplex128Ptr = reflect.PtrTo(GoComplex128)
	GoBytePtr       = reflect.PtrTo(GoByte)
	GoStringPtr     = reflect.PtrTo(GoString)
	GoBoolPtr       = reflect.PtrTo(GoByte)
	GoTimePtr       = reflect.PtrTo(GoTime)
)

func GoTypeToMySqlType() {

}

func MySqlTypeToGoType() {

}
