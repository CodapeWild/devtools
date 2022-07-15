package workerpool

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	m.Run()
}

func TestWorkerPool(t *testing.T) {
	wpool := NewWorkerPool(100)
	err := wpool.Start(16)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	wg := sync.WaitGroup{}
	wg.Add(10000)
	start := time.Now()
	for i := 0; i < 10000; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				job, err := NewJob(func(input interface{}) (output interface{}, err error) {
					log.Printf("get new job, input: %v\n", input)

					return fmt.Sprintf("finish job %d in goroutine %d\n", j, i), nil
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
	wg.Wait()

	log.Printf("send jobs finished, cost %ds\n", time.Since(start)/time.Second)

	time.Sleep(1 * time.Second)
	wpool.Shutdown()
}
