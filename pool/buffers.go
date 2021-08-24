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

// GetBuffer get the empty Buffer from pool, the Buffer could contain underlying storage.
func GetBuffer() *bytes.Buffer {
	temp := bufpool.Get().(*bytes.Buffer)
	temp.Reset()

	return temp
}

// RestoreBuffer restore the Buffer into pool.
func RestoreBuffer(temp *bytes.Buffer) {
	bufpool.Put(temp)
}
