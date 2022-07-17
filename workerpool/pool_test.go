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
	wpool := NewWorkerPool(200)
	err := wpool.Start(8)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	var n, m = 1000, 1000
	wg := sync.WaitGroup{}
	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < m; j++ {
				job, err := NewJob(
					WithInput(fmt.Sprintf("goroutine %d, job %d\n", i, j)),
					WithTimeout(time.Second),
					WithProcess(func(input interface{}) (output interface{}) {
						log.Printf("start process input %d:%d\n", i, j)

						return fmt.Sprintf("finish process %d:%d\n", i, j)
					}),
					WithProcessCallback(func(input, output interface{}, cost time.Duration, isTimeout bool) {
						log.Printf("finish process and callback, input: %v, output: %v, cost: %d isTimeout: %v\n", input, output, cost, isTimeout)
					}),
				)
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
