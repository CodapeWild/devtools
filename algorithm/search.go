package algorithm

type CompareResult int8

const (
	Less    CompareResult = -1
	Equal   CompareResult = 0
	Greater CompareResult = 1
)

type Searchable interface {
	Len() int
	Compare(i int, t interface{}) CompareResult
}

type IntSearchable []int

func (this IntSearchable) Len() int {
	return len(this)
}

func (this IntSearchable) Compare(i int, t interface{}) CompareResult {
	if this[i] > t.(int) {
		return Greater
	} else if this[i] < t.(int) {
		return Less
	} else {
		return Equal
	}
}

func BinarySearch(data Searchable, target interface{}) int {
	if data.Len() == 0 {
		return -1
	}

	var (
		i, j = 0, data.Len()
		m    = (i + j) / 2
	)
	for i < j {
		switch data.Compare(m, target) {
		case Greater:
			j = m - 1
		case Less:
			i = m + 1
		case Equal:
			return m
		}
		m = (i + j) / 2
	}

	return -1
}
