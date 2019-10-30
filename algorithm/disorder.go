package algorithm

import (
	"math/rand"
)

func Disorder(data SeqData, start, end int) {
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
