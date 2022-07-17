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
	wpool := NewWorkerPool(500000)
	err := wpool.Start(16)
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
						log.Printf("finish process and callback, input: %v, output: %v, cost: %dms isTimeout: %v\n", input, output, cost/time.Millisecond, isTimeout)
					}),
				)
				if err != nil {
					t.Errorf(err.Error())
					continue
				}

				// if err = wpool.MoreJobsSync(job); err != nil {
				// 	log.Println(err.Error())
				// }
				if err = wpool.MoreJobsWithoutTimeout(job); err != nil {
					log.Println(err.Error())
				}
				// if err = wpool.MoreJobsWithTimeout(10*time.Millisecond, job); err != nil {
				// 	log.Println(err.Error())
				// }
			}
		}(i)
	}
	wg.Wait()

	log.Printf("send jobs finished, cost %dms\n", time.Since(start)/time.Millisecond)

	time.Sleep(10 * time.Second)
	wpool.Shutdown()
}
