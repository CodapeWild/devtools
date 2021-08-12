package msgque

import (
	"sync"
)

type Cache interface {
	Push(obj interface{}) bool
	Pop() interface{}
	Clear()
	Len() int
}

type MemCache struct {
	mem []interface{}
	sync.Mutex
}

func (this *MemCache) Push(obj interface{}) bool {
	this.Lock()
	defer this.Unlock()

	this.mem = append(this.mem, obj)

	return true
}

func (this *MemCache) Pop() interface{} {
	var tmp interface{}
	if len(this.mem) != 0 {
		this.Lock()
		defer this.Unlock()

		tmp = this.mem[len(this.mem)-1]
		this.mem = this.mem[:len(this.mem)-1]
	}

	return tmp
}

func (this *MemCache) Clear() {
	if len(this.mem) != 0 {
		this.Lock()
		defer this.Unlock()

		this.mem = this.mem[:0]
	}
}

func (this *MemCache) Len() int {
	return len(this.mem)
}
