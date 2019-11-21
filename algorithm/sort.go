package algorithm

func QuickSort(data SeqData, start, end int) {
	if start >= end {
		return
	}

	Disorder(data, start, end)
	var i, j, pivot = start, start + 1, start
	for j < end {
		if data.Less(j, pivot) {
			if i++; i != j {
				data.Swap(i, j)
			}
		}
		j++
	}
	data.Swap(pivot, i)

	QuickSort(data, start, i)
	QuickSort(data, i+1, end)
}
