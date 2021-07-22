package configure

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/CodapeWild/devtools/comerr"
)

type Configure interface {
	ReadFrom(r io.Reader) error
	WriteTo(w io.Writer) error
	Marshal(v interface{})
	Unmarshal(v interface{})
}

type JsonFileConfigure struct {
	buf []byte
}

func (this *JsonFileConfigure) ReadFrom(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	this.buf = buf

	return nil
}

func (this *JsonFileConfigure) WriteTo(w io.Writer) error {
	var err error
	if len(this.buf) != 0 {
		_, err = w.Write(this.buf)
	} else {
		err = comerr.ErrEmptyData
	}

	return err
}

func (this *JsonFileConfigure) Marshal(v interface{}) ([]byte, error) {
	var err error
	this.buf, err = json.Marshal(v)

	return this.buf, err
}

func (this *JsonFileConfigure) Unmarshal(v interface{}) error {
	if rv := reflect.ValueOf(v); rv.Kind() != reflect.Ptr || rv.IsNil() {
		return comerr.ErrTypeInvalid
	}
	if !json.Valid(this.buf) {
		return comerr.ErrParamInvalid
	}

	return json.Unmarshal(this.buf, v)
}
