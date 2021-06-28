package algorithm

var (
	CompareInt IntComp = func(lopr, ropr int) int {
		if lopr > ropr {
			return 1
		} else if lopr < ropr {
			return -1
		} else {
			return 0
		}
	}
	CompareByte ByteComp = func(lopr, ropr byte) int {
		if lopr > ropr {
			return 1
		} else if lopr < ropr {
			return -1
		} else {
			return 0
		}
	}
)

type Comparator interface {
	// if lopr > ropr return 1
	// if lopr < ropr return -1
	// if lopr == ropr return 0
	Compare(lopr, ropr interface{}) int
}

type IntComp func(lopr, ropr int) int

func (this IntComp) Compare(lopr, ropr interface{}) int {
	return this(lopr.(int), ropr.(int))
}

type ByteComp func(lopr, ropr byte) int

func (this ByteComp) Compare(lopr, ropr interface{}) int {
	return this(lopr.(byte), ropr.(byte))
}

type CompareFunc func(comparator Comparator) func(lopr, ropr interface{}) bool

func Bigger(comparator Comparator) func(lopr, ropr interface{}) bool {
	return func(lopr, ropr interface{}) bool {
		if comparator.Compare(lopr, ropr) >= 0 {
			return true
		} else {
			return false
		}
	}
}

func Smaller(comparator Comparator) func(lopr, ropr interface{}) bool {
	return func(lopr, ropr interface{}) bool {
		if comparator.Compare(lopr, ropr) <= 0 {
			return true
		} else {
			return false
		}
	}
}
