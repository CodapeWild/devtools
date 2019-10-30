package algorithm

import (
	"log"
	"sort"
	"testing"
)

var (
	// 						0	 1   2  3  4  5   6    7  8  9  10 11  12
	data  = []int{34, 5, 67, 8, 9, 8, 76, 543, 2, 3, 45, 6, 89}
	sdata = sort.IntSlice(data)
)

func TestDisorder(t *testing.T) {
	log.Println(sdata)
	Disorder(sdata, 7, 12)
	log.Println(sdata)
}

func TestQuickSort(t *testing.T) {
	log.Println(sdata)
	QuickSort(sdata, 0, sdata.Len())
	log.Println(sdata)
}

func TestQuickLocate(t *testing.T) {
	order := 9
	log.Println(QuickLocate(sdata, order), order, sdata[order])
	QuickSort(sdata, 0, sdata.Len())
	log.Println(sdata)
}

func TestBinarySearch(t *testing.T) {
	QuickSort(sdata, 0, sdata.Len())
	ordered := OrderedInts(sdata)
	log.Println(ordered)
	log.Println(BinarySearch(ordered, 34))
}
