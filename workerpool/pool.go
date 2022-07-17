package workerpool

import (
	"errors"
	"fmt"
	"time"
)

type Process func(input interface{}) (output interface{})

type ProcessCallback func(input, output interface{}, cost time.Duration, isTimeout bool)

type Job struct {
	input   interface{}
	timeout time.Duration
	p       Process
	cb      ProcessCallback
}

type Option func(job *Job)

func WithInput(input interface{}) Option {
	return func(job *Job) {
		job.input = input
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(job *Job) {
		job.timeout = timeout
	}
}

func WithProcess(p Process) Option {
	return func(job *Job) {
		job.p = p
	}
}

func WithProcessCallback(cb ProcessCallback) Option {
	return func(job *Job) {
		job.cb = cb
	}
}

func NewJob(options ...Option) (*Job, error) {
	job := &Job{}
	for i := range options {
		options[i](job)
	}
	if job.p == nil {
		return nil, errors.New("process can not be nil")
	}

	return job, nil
}

type WorkerPool chan *Job

func NewWorkerPool(buffer int) WorkerPool {
	if buffer < 0 {
		return nil
	}

	return make(WorkerPool, buffer)
}

func (wp WorkerPool) Start(threads int) error {
	if wp == nil {
		return errors.New("worker pool is not ready")
	}
	if threads < 1 {
		return errors.New("worker pool needs at least one thread")
	}

	for i := 0; i < threads; i++ {
		go wp.worker()
	}

	return nil
}

func (wp WorkerPool) Shutdown() {
	close(wp)
}

func (wp WorkerPool) MoreJobsSync(jobs ...*Job) error {
	if wp == nil {
		return errors.New("worker pool is not ready")
	}
	if len(jobs) == 0 {
		return nil
	}

	for i := range jobs {
		wp <- jobs[i]
	}

	return nil
}

func (wp WorkerPool) MoreJobsWithoutTimeout(jobs ...*Job) error {
	if wp == nil {
		return errors.New("worker pool is not ready")
	}
	if len(jobs) == 0 {
		return nil
	}

	var busyCount int
	for i := range jobs {
		select {
		case wp <- jobs[i]:
		default:
			busyCount++
		}
	}
	if busyCount != 0 {
		return fmt.Errorf("worker pool is busy, failed send jobs %d", busyCount)
	}

	return nil
}

func (wp WorkerPool) MoreJobsWithTimeout(timeout time.Duration, jobs ...*Job) error {
	if wp == nil {
		return errors.New("worker pool is not ready")
	}
	if len(jobs) == 0 {
		return nil
	}

	var (
		tick    = time.NewTicker(timeout)
		success int
	)
	for i := range jobs {
		select {
		case <-tick.C:
			return fmt.Errorf("send jobs to worker pool timeout, successfully send %d", success)
		case wp <- jobs[i]:
			success++
		}
	}
	tick.Stop()

	return nil
}

func (wp WorkerPool) worker() {
	if wp == nil {
		return
	}

	for {
		job, ok := <-wp
		if !ok {
			break
		}
		if job == nil {
			continue
		}

		var start = time.Now()
		if job.timeout < 1 {
			r := job.p(job.input)
			if job.cb != nil {
				job.cb(job.input, r, time.Since(start), false)
			}
		} else {
			result := make(chan interface{}, 1)
			go func(job *Job, r chan interface{}) {
				r <- job.p(job.input)
				close(r)
			}(job, result)

			var tick = time.NewTicker(job.timeout)
			select {
			case <-tick.C:
				if job.cb != nil {
					job.cb(job.input, nil, time.Since(start), true)
				}
			case r := <-result:
				if job.cb != nil {
					job.cb(job.input, r, time.Since(start), false)
				}
				tick.Stop()
			}
		}
	}
}
