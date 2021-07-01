package algorithm

type RandomAccess interface {
	Get(key interface{}) interface{}
	Len() int
	Previous() (key, value interface{})
	Current() (key, value interface{})
	Next() (key, value interface{})
}

type IntsRandAcc struct {
	data []int
	i    int
}

func NewIntsRandAcc(data []int) IntsRandAcc {
	return IntsRandAcc{data: data}
}

func (this IntsRandAcc) Get(key interface{}) interface{} {
	return this.data[key.(int)]
}

func (this IntsRandAcc) Len() int {
	return len(this.data)
}

func (this IntsRandAcc) Previous() (key, value interface{}) {
	if this.i-1 >= 0 {
		this.i--
		key = this.i
		value = this.data[this.i]
	}

	return
}

func (this IntsRandAcc) Current() (key, value interface{}) {
	if this.i >= 0 && this.i < this.Len() {
		key = this.i
		value = this.data[this.i]
	}

	return
}

func (this IntsRandAcc) Next() (key, value interface{}) {
	if this.i+1 < this.Len() {
		this.i++
		key = this.i
		value = this.data[this.i]
	}

	return
}

type FloatsRandAcc struct {
	data []float64
	i    int
}

func NewFloatsRandAcc(data []float64) FloatsRandAcc {
	return FloatsRandAcc{data: data}
}

func (this FloatsRandAcc) Get(key interface{}) interface{} {
	return this.data[key.(int)]
}

func (this FloatsRandAcc) Len() int {
	return len(this.data)
}

func (this FloatsRandAcc) Previous() (key, value interface{}) {
	if this.i-1 >= 0 {
		this.i--
		key = this.i
		value = this.data[this.i]
	}

	return
}

func (this FloatsRandAcc) Current() (key, value interface{}) {
	if this.i >= 0 && this.i < this.Len() {
		key = this.i
		value = this.data[this.i]
	}

	return
}

func (this FloatsRandAcc) Next() (key, value interface{}) {
	if this.i+1 < this.Len() {
		this.i++
		key = this.i
		value = this.data[this.i]
	}

	return
}

type StringsRandAcc struct {
	data []string
	i    int
}

func NewStringsRandAcc(data []string) StringsRandAcc {
	return StringsRandAcc{data: data}
}

func (this StringsRandAcc) Get(key interface{}) interface{} {
	return this.data[key.(int)]
}

func (this StringsRandAcc) Len() int {
	return len(this.data)
}

func (this StringsRandAcc) Previous() (key, value interface{}) {
	if this.i-1 >= 0 {
		this.i--
		key = this.i
		value = this.data[this.i]
	}

	return
}

func (this StringsRandAcc) Current() (key, value interface{}) {
	if this.i >= 0 && this.i < this.Len() {
		key = this.i
		value = this.data[this.i]
	}

	return
}

func (this StringsRandAcc) Next() (key, value interface{}) {
	if this.i+1 < this.Len() {
		this.i++
		key = this.i
		value = this.data[this.i]
	}

	return
}

type StrStrMapRandAcc struct {
	data map[string]string
	keys []string
	i    int
}

func NewStrStrMapRandAcc(data map[string]string) StrStrMapRandAcc {
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	return StrStrMapRandAcc{
		data: data,
		keys: keys,
	}
}

func (this StrStrMapRandAcc) Get(key interface{}) interface{} {
	return this.data[key.(string)]
}

func (this StrStrMapRandAcc) Len() int {
	return len(this.data)
}

func (this StrStrMapRandAcc) Previous() (key, value interface{}) {
	if this.i-1 >= 0 {
		this.i--
		key = this.keys[this.i]
		value = this.data[this.keys[this.i]]
	}

	return
}

func (this StrStrMapRandAcc) Current() (key, value interface{}) {
	if this.i >= 0 && this.i < this.Len() {
		key = this.keys[this.i]
		value = this.data[this.keys[this.i]]
	}

	return
}

func (this StrStrMapRandAcc) Next() (key, value interface{}) {
	if this.i+1 < this.Len() {
		this.i++
		key = this.keys[this.i]
		value = this.data[this.keys[this.i]]
	}

	return
}
