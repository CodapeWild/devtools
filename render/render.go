package render

import "sync"

type Render struct {
	threads   []Thread
	worker    chan Thread
	maxWorker int
	sync.Mutex
}

func NewRender() *Render {

}

func (this *Render) AddTask(t Thread) {
	if t == nil {
		return
	}

	this.Lock()
	defer this.Unlock()

	this.threads = append(this.threads, t)
}

func (this *Render) Start() {

}

func (this *Render) Status() string {

}

func (this *Render) Close() {

}
