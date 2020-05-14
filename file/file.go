package file

import (
	"devtools/comerr"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"reflect"
)

func ReadJsonFile(filePath string, out interface{}) error {
	if !IsFileExists(filePath) {
		return FileNotExists
	}
	if out == nil {
		return comerr.ParamInvalid
	}
	if v := reflect.TypeOf(out); v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return comerr.TypeInvalid
	}

	if buf, err := ioutil.ReadFile(filePath); err != nil {
		return err
	} else {
		return json.Unmarshal(buf, out)
	}
}

func WriteJsonFile(filePath string, in interface{}) error {
	buf, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, buf, 0644)
}

func ReadGoTplFile(param interface{}, funcMap template.FuncMap, w io.Writer, names ...string) error {
	if param != nil {
		t := reflect.TypeOf(param)
		k := t.Kind()
		if k == reflect.Ptr {
			k = t.Elem().Kind()
		}
		if k != reflect.Struct && k != reflect.Map {
			return comerr.ParamInvalid
		}
	}

	if len(funcMap) != 0 {
		return template.Must(template.New(path.Base(names[0])).Funcs(funcMap).ParseFiles(names...)).Execute(w, param)
	} else {
		return template.Must(template.ParseFiles(names...)).Execute(w, param)
	}
}
