package code

import (
	"bytes"
	"reflect"

	"github.com/CodapeWild/devtools/comerr"
	"github.com/CodapeWild/devtools/pool"
	"github.com/hashicorp/go-msgpack/codec"
)

var msgpHandler = &codec.MsgpackHandle{}

func MsgPackMarshal(v interface{}) ([]byte, error) {
	buf := pool.GetBuffer()
	defer pool.RestoreBuffer(buf)
	err := codec.NewEncoder(buf, msgpHandler).Encode(v)

	return buf.Bytes(), err
}

func MsgPackUnmarshal(buf []byte, out interface{}) error {
	if out == nil {
		return comerr.ErrNilPointer
	}
	if rt := reflect.ValueOf(out).Kind(); rt != reflect.Ptr || rt != reflect.Map {
		return comerr.ErrTypeInvalid
	}

	return codec.NewDecoder(bytes.NewBuffer(buf), msgpHandler).Decode(out)
}
