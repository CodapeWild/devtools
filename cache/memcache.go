package cache

import "sync"

type MemCache struct {
	mem []interface{}
	sync.Mutex
}

func (this *MemCache) Push(v interface{}) error {
	this.Lock()
	defer this.Unlock()

	this.mem = append(this.mem, v)

	return nil
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
	this.mem = this.mem[:0]
}

func (this *MemCache) Len() int {
	return len(this.mem)
}
