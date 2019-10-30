package algorithm

import "devtools/comerr"

func QuickLocate(data SeqData, ith int) error {
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
func quickSortOnce(data SeqData, start, end int) int {
	var (
		i     = start
		j     = i + 1
		pivot = start
	)
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
