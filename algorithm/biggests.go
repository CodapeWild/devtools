package algorithm

type Comparable interface {
	Bigger(lval, rval interface{}) bool
}
