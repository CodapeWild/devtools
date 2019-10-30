package algorithm

func QuickSort(data SeqData, start, end int) {
	if data == nil || data.Len() == 0 || start >= end {
		return
	}

	Disorder(data, start, end)
	var (
		pivot = start
		i     = start
		j     = start + 1
	)
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
