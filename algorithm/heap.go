package algorithm

type Heap struct {
	data       []interface{}
	comparator Comparator
	compare    CompareFunc
}

func NewHeap(comparator Comparator, compare CompareFunc, data ...interface{}) *Heap {
	heap := &Heap{
		comparator: CompareByte,
		compare:    compare,
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

func (this *Heap) parent(child int) int {
	return (child - 1) / 2
}

func (this *Heap) lchild(parent int) int {
	return parent*2 + 1
}

func (this *Heap) rchild(parent int) int {
	return parent*2 + 2
}

func (this *Heap) heapifyUp() {
	var (
		child  = this.Count() - 1
		parent = this.parent(child)
	)
	for parent > 0 {
		if this.compare(this.comparator)(this.data[child], this.data[parent]) {
			this.data[child], this.data[parent] = this.data[parent], this.data[child]
			child = parent
			parent = this.parent(child)
		} else {
			break
		}
	}
}

func (this *Heap) heapifyDown() {
	var (
		parent = 0
		lchild = this.lchild(parent)
		rchild = this.rchild(parent)
	)
	for {
		tmp := parent
		if lchild <= this.Count()-1 && this.compare(this.comparator)(this.data[parent], this.data[lchild]) {
			tmp = lchild
		}
		if rchild <= this.Count()-1 && this.compare(this.comparator)(this.data[parent], this.data[rchild]) {
			tmp = rchild
		}
		if tmp != parent {
			this.data[parent], this.data[tmp] = this.data[tmp], this.data[parent]
			parent = tmp
			lchild = this.lchild(parent)
			rchild = this.rchild(parent)
		} else {
			break
		}
	}
}
