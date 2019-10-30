package idflaker

import (
	"fmt"
	"log"
	"testing"
)

func TestIdFlaker(t *testing.T) {
	fk, err := NewIdFlaker(1)
	if err != nil {
		log.Fatalln(err.Error())
	}

	id := fk.NextInt64Id()
	fmt.Println(id)
	fmt.Println(fk.id>>53, fk.ts, fk.seq)
	fmt.Println(ParseInt64Id(id))
	fmt.Println(fk.NextStrId())
}

func TestEfficiency(t *testing.T) {
	fk1, _ := NewIdFlaker(1)
	c := 10
	unique := make(map[int64]bool)
	ch := make(chan int64, 1000)

	routine := func(fk *IdFlaker) {
		for i := 0; i < 1000000; i++ {
			id := fk.NextInt64Id()
			ch <- id
		}
		ch <- 0
	}

	for i := 0; i < c; i++ {
		go routine(fk1)
	}

	count := 0
	for v := range ch {
		if v == 0 {
			if count++; count == c {
				fmt.Println("efficiency test completed")

				return
			}
			continue
		}

		if !unique[v] {
			unique[v] = true
		} else {
			fmt.Println("produced duplicated id, error")
			fmt.Println(unique[v], v)

			return
		}
	}
}
