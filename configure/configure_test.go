package configure

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var jsonbuf = `
{
	"name":"tnt",
	"data": [1,2,3,4,5,6],
	"kv": {"day":"Monday","month":"July"}
}
`

func TestJsonConfigure(t *testing.T) {
	conf := &JsonFileConfigure{}
	err := conf.ReadFrom(bytes.NewBuffer([]byte(jsonbuf)))
	if err != nil {
		log.Panicln(err.Error())
	}
	m := make(map[string]interface{})
	if err = conf.Unmarshal(&m); err != nil {
		log.Panicln(err.Error())
	}
	var buf []byte
	if buf, err = conf.Marshal(&m); err != nil {
		log.Panicln(err.Error())
	} else {
		log.Println(string(buf))
	}
	var f *os.File
	if f, err = os.OpenFile("./tmp.json", os.O_CREATE|os.O_RDWR, 0777); err != nil {
		f.Close()
		log.Panicln(err.Error())
	} else {
		if err = conf.WriteTo(f); err != nil {
			log.Panicln(err.Error())
		}
	}
}
