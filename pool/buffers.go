package pool

import (
	"bytes"
	"sync"
)

var bufpool sync.Pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

func GetBuffer() *bytes.Buffer {
	temp := bufpool.Get().(*bytes.Buffer)
	temp.Reset()

	return temp
}

func RestoreBuffer(temp *bytes.Buffer) {
	bufpool.Put(temp)
}
