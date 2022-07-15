package workerpool

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	wpool := NewWorkerPool(10)
	err := wpool.Start(3)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			for j := 0; j < 10; j++ {
				job, err := NewJob(func(input interface{}) (output interface{}, err error) {
					log.Printf("input: %v\n", input)

					return fmt.Sprintf("finish job %d\n", j), nil
				},
					WithInput(fmt.Sprintf("goroutine %d, job %d\n", i, j)),
					WithCallback(func(input, output interface{}, cost time.Duration, err error) {
						log.Printf("result: %v cost: %dns, err: %v\n", output, cost, err)
					}),
					WithTimeout(time.Second))
				if err != nil {
					t.Errorf(err.Error())
					continue
				}

				wpool.MoreWork(time.Second, job)
			}
		}(i)
	}

	select {}
}
