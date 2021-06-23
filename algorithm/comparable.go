package algorithm

type Comparable interface {
	Compare(lopr, ropr interface{}) int
}

type IntComp func(lopr, ropr int) int

func (this IntComp) Compare(lopr, ropr interface{}) int {
	return this(lopr.(int), ropr.(int))
}

func CompareInt(lopr, ropr int) int {
	if lopr > ropr {
		return 1
	} else if lopr < ropr {
		return -1
	} else {
		return 0
	}
}
