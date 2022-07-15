package workerpool

import (
	"errors"
	"time"
)

type Process func(input interface{}) (output interface{}, err error)

type Callback func(input, output interface{}, cost time.Duration, err error)

type result struct {
	output interface{}
	err    error
}

type Job struct {
	input   interface{}
	p       Process
	cb      Callback
	timeout time.Duration
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

func WithCallback(cb Callback) Option {
	return func(job *Job) {
		job.cb = cb
	}
}

func NewJob(p Process, options ...Option) (*Job, error) {
	if p == nil {
		return nil, errors.New("process can not be nil.")
	}

	job := &Job{p: p}
	for i := range options {
		options[i](job)
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
		return errors.New("worker pool is not ready.")
	}
	if threads < 1 {
		return errors.New("worker pool needs at least one thread.")
	}

	for i := 0; i < threads; i++ {
		go wp.worker()
	}

	return nil
}

func (wp WorkerPool) Stop() {
	close(wp)
}

func (wp WorkerPool) MoreWork(sendTimeout time.Duration, jobs ...*Job) error {
	if wp == nil {
		return errors.New("worker pool is not ready.")
	}
	if len(jobs) == 0 {
		return errors.New("job list can not be empty.")
	}

	if sendTimeout < 1 {
		for i := range jobs {
			wp <- jobs[i]
		}
	} else {
		tick := time.NewTicker(sendTimeout)
		for i := range jobs {
			select {
			case <-tick.C:
				return errors.New("send job to worker pool timeout.")
			case wp <- jobs[i]:
			}
		}
		tick.Stop()
	}

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
			output, err := job.p(job.input)
			if job.cb != nil {
				job.cb(job.input, output, time.Since(start), err)
			}
		} else {
			var (
				tick = time.NewTicker(job.timeout)
				rslt = make(chan *result)
			)
			go func() {
				r, err := job.p(job.input)
				rslt <- &result{r, err}
			}()
			select {
			case <-tick.C:
				if job.cb != nil {
					job.cb(job.input, nil, time.Since(start), errors.New("job process timeout."))
				}
			case r := <-rslt:
				if job.cb != nil {
					job.cb(job.input, r.output, time.Since(start), r.err)
				}
				tick.Stop()
				close(rslt)
			}
		}
	}
}
