package algorithm

import (
	"devtools/comerr"
	"sort"
)

type ReverseSort struct {
	sort.Interface
}

func (this *ReverseSort) Less(i, j int) bool {
	return !this.Interface.Less(i, j)
}

func QuickSort(data sort.Interface, start, end int) {
	if start >= end {
		return
	}

	Disorder(data, start, end)
	var i, j, pivot = start, start + 1, start
	for j < end {
		if data.Less(j, pivot) {
			data.Swap(i+1, j)
			i++
		}
		j++
	}
	data.Swap(pivot, i)

	QuickSort(data, start, i)
	QuickSort(data, i+1, end)
}

// quick sort data from start to end
func QuickSortOverall(data sort.Interface) {
	QuickSort(data, 0, data.Len())
}

func QuickLocate(data sort.Interface, ith int) error {
	if data == nil || data.Len() <= ith {
		return comerr.ParamInvalid
	}

	var (
		s, e  = 0, data.Len()
		pivot int
	)
RELOCATE:
	pivot = quickSortOnce(data, s, e)
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
func quickSortOnce(data sort.Interface, start, end int) int {
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
