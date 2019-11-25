package algorithm

import (
	"math/rand"
)

type Swapable interface {
	Len() int
	Swap(i, j int)
}

func Disorder(data Swapable, start, end int) {
	if data == nil || start > end || data.Len() < end-start {
		return
	}

	switch end - start {
	case 0, 1:
	case 2:
		data.Swap(start, start+rand.Intn(2))
	default:
		for i := end - 1; i > start; i-- {
			data.Swap(start+rand.Intn(i-start), i)
		}
	}
}
