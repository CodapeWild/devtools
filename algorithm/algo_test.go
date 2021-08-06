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
	Quicksort(sdata, 0, len(sdata))
	log.Println(sdata)
	log.Println(BinarySearch(IntSearchable(sdata), 2))
	log.Println(BinarySearch(IntSearchable(sdata), 9))
	log.Println(BinarySearch(IntSearchable(sdata), 543))
	log.Println(BinarySearch(IntSearchable(sdata), 544))
}

func TestDisorder(t *testing.T) {
	Quicksort(sdata, 0, len(sdata))
	log.Println(sdata)
	Disorder(sdata, 7, 12)
	log.Println(sdata)
}

func TestQuicksort(t *testing.T) {
	log.Println(sdata)
	Quicksort(sdata, 0, sdata.Len())
	log.Println(sdata)
}

func TestQuickLocate(t *testing.T) {
	order := 9
	log.Println(QuickLocate(sdata, order), order, sdata[order])
	Quicksort(sdata, 0, sdata.Len())
	log.Println(sdata)
}

func makeChange(denomi []int, n int, remain int) int {
	if remain == 0 {
		return 1
	}
	if n < 0 {
		return 0
	}
	if remain < denomi[n] {
		return makeChange(denomi, n-1, remain)
	}

	var count int
	for i := 0; remain-i*denomi[n] >= 0; i++ {
		count += makeChange(denomi, n-1, remain-i*denomi[n])
	}

	return count
}

func TestMakeChange(t *testing.T) {
	denomi := []int{1, 2, 5, 10, 15, 50}
	log.Println(makeChange(denomi, len(denomi)-1, 3))
}

func permuteHandlerTest(data []int, p []int, output chan []int) {
	if len(data) == 0 {
		output <- p
	} else {
		for i := 0; i < len(data); i++ {
			ptmp := make([]int, len(p))
			copy(ptmp, p)
			ptmp = append(ptmp, data[i])

			datatmp := make([]int, len(data))
			copy(datatmp, data)
			datatmp[i] = datatmp[len(datatmp)-1]

			permuteHandlerTest(datatmp[:len(datatmp)-1], ptmp, output)
		}
	}
}

func permute(data []int, output chan []int) {
	permuteHandlerTest(data, nil, output)
}

func TestPermute(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	output := make(chan []int)
	go func() {
		for v := range output {
			log.Println(v)
		}
	}()
	permute(data, output)
}

func TestPermute1(t *testing.T) {
	data := IntsPermutable([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	output := make(chan interface{})
	go func() {
		for v := range output {
			log.Println(v)
		}
	}()
	Permute(data, output)
}

func TestTrie(t *testing.T) {
	cases := []string{"123", "12", "1234", "234", "321", "1256"}
	trie := NewTrie()
	for _, v := range cases {
		trie.Add(v)
	}

	trie.Add("12")

	for _, v := range cases {
		log.Println(trie.Find(v))
	}

	log.Println(trie.origin)
}

func TestHeapify(t *testing.T) {
	hp := NewHeap(CompareInt, Bigger, func(data []int) []interface{} {
		inters := make([]interface{}, len(data))
		for k, v := range data {
			inters[k] = v
		}

		return inters
	}([]int{3, 68, 7, 65, 45, 6, 9, 8})...)

	log.Printf("heapified: %v\n", hp.data)

	removed := hp.Remove()
	for removed != nil {
		log.Printf("removed: %v, heap: %v", removed, hp.data)
		removed = hp.Remove()
	}
}

func TestInterfacesAdaptor(t *testing.T) {
	data := []int{12, 2, 3, 4, 432, 56, 42, 327, 8, 9}
	ira := NewIntsRandAcc(data)
	k, v := ira.Next()
	for v != nil {
		log.Println(k, v)
		k, v = ira.Next()
	}

	m := map[string]string{"name": "tnt", "age": "123", "rank": "999"}
	mra := NewStrStrMapRandAcc(m)
	k, v = mra.Next()
	for v != nil {
		log.Println(k, v)
		k, v = mra.Next()
	}
}
