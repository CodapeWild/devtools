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

func TestBinarySearch(t *testing.T) {
	QuickSort(sdata, 0, len(sdata))
	log.Println(sdata)
	log.Println(BinarySearch(IntSearchable(sdata), 2))
	log.Println(BinarySearch(IntSearchable(sdata), 9))
	log.Println(BinarySearch(IntSearchable(sdata), 543))
	log.Println(BinarySearch(IntSearchable(sdata), 544))
}

func TestDisorder(t *testing.T) {
	QuickSort(sdata, 0, len(sdata))
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

func TestTrie(t *testing.T) {
	root1 := NewTrieNode(TrieString(""), 0, nil)
	root1.Add(TrieString("abc"), TrieString("ab"), TrieString("a"))
	ShowTrie(root1)
	log.Println("###################")
	root2 := NewTrieNode(TrieString("a"), 1, nil)
	root2.Add(TrieString("ab"), TrieString("abc"))
	ShowTrie(root2)
	log.Println("###################")
	root3 := NewTrieNode(TrieString("abc"), 1, nil)
	root3.Add(TrieString("ab"), TrieString("a"), TrieString("abc"), TrieString("abcdxy"))
	ShowTrie(root3)
}
