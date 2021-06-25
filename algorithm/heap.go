package algorithm

type Heap struct {
	data       []interface{}
	comparator Comparator
}

func NewHeap(comparator Comparator, data ...interface{}) *Heap {
	heap := &Heap{
		comparator: comparator,
	}
	heap.Insert(data...)

	return heap
}

func (this *Heap) Peek() interface{} {
	if len(this.data) != 0 {
		return this.data[0]
	} else {
		return nil
	}
}

func (this *Heap) Count() int {
	return len(this.data)
}

func (this *Heap) Insert(data ...interface{}) {
	for _, v := range data {
		this.data = append(this.data, v)
		this.heapifyUp()
	}
}

func (this *Heap) Remove() interface{} {
	var (
		c   = this.Count()
		tmp interface{}
	)
	switch c {
	case 0:
	case 1:
		tmp = this.data[0]
		this.data = nil
	case 2:
		this.data[0], this.data[1] = this.data[1], this.data[0]
		tmp = this.data[1]
		this.data = this.data[:1]
	default:
		this.data[0], this.data[c-1] = this.data[c-1], this.data[0]
		tmp = this.data[c-1]
		this.data = this.data[:c-1]
		this.heapifyDown()
	}

	return tmp
}

func (this *Heap) heapifyUp() {
	var (
		child  = this.Count() - 1
		parent = child / 2
	)
	for parent > 0 {
		if this.comparator.Compare(this.data[child])
	}
}

func (this *Heap) heapifyDown() {

}
