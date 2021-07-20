package algorithm

import (
	"devtools/comerr"
	"math/rand"
	"sort"
)

type ReverseSort struct {
	sort.Interface
}

func (this *ReverseSort) Less(i, j int) bool {
	return !this.Interface.Less(i, j)
}

func Quicksort(data sort.Interface, left, right int) {
	if left >= right {
		return
	}

	// create random pivot
	data.Swap(left+rand.Intn(right-left), left)
	var pivot, candi, cur = left, left, left + 1
	for cur < right {
		if data.Less(cur, pivot) {
			candi++
			data.Swap(cur, candi)
		}
		cur++
	}
	data.Swap(pivot, candi)

	Quicksort(data, left, candi)
	Quicksort(data, candi+1, right)
}

// quick sort data from start to end
func QuicksortOverall(data sort.Interface) {
	Quicksort(data, 0, data.Len())
}

func QuickLocate(data sort.Interface, ith int) error {
	if data == nil || data.Len() <= ith {
		return comerr.ErrParamInvalid
	}

	var (
		s, e  = 0, data.Len()
		pivot int
	)
RELOCATE:
	pivot = quicksortOnce(data, s, e)
	if pivot < ith {
		s++
		goto RELOCATE
	} else if pivot > ith {
		e = pivot
		goto RELOCATE
	}

	return nil
}

// return the pivot index after sorting once
func quicksortOnce(data sort.Interface, start, end int) int {
	var i, j, pivot = start, start + 1, start
	for j < end {
		if data.Less(j, pivot) {
			data.Swap(i+1, j)
			i++
		}
		j++
	}
	data.Swap(pivot, i)

	return i
}
