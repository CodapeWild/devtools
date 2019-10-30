package algorithm

type SeqData interface {
	Swap(i, j int)
	Less(i, j int) bool
	Len() int
}

type OrderedData interface {
	Data(i int) interface{}
	Greater(m int, target interface{}) bool
	Less(m int, target interface{}) bool
	Len() int
}

type OrderedInts []int

func (this OrderedInts) Data(i int) interface{} {
	return this[i]
}

func (this OrderedInts) Greater(m int, target interface{}) bool {
	return this[m] > target.(int)
}

func (this OrderedInts) Less(m int, target interface{}) bool {
	return this[m] < target.(int)
}

func (this OrderedInts) Len() int {
	return len(this)
}

type OrderedStrings []string

func (this OrderedStrings) Data(i int) interface{} {
	return this[i]
}

func (this OrderedStrings) Greater(m int, target interface{}) bool {
	return this[m] > target.(string)
}

func (this OrderedStrings) Less(m int, target interface{}) bool {
	return this[m] < target.(string)
}

func (this OrderedStrings) Len() int {
	return len(this)
}
