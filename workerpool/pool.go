package workerpool

import (
	"errors"
	"log"
	"time"
)

type Process func(interface{}) interface{}

type Job struct {
	input   interface{}
	output  chan interface{}
	p       Process
	timeout time.Duration
}

func NewJob(input interface{}, output chan interface{}, p Process, timeout time.Duration) (*Job, error) {
	if output == nil {
		return nil, errors.New("output chan cannot be nil.")
	}
	if p == nil {
		return nil, errors.New("process chan not be nil.")
	}

	return &Job{
		input:   input,
		output:  output,
		p:       p,
		timeout: timeout,
	}, nil
}

func (j *Job) WaitResult() (interface{}, error) {
	r, ok := <-j.output
	if ok {
		return r, nil
	}

	return nil, errors.New("wait job result timeout.")
}

type WorkerPool chan *Job

func NewWorkerPool() WorkerPool {
	return make(WorkerPool)
}

func (wp WorkerPool) Start(threads int) error {
	if wp == nil {
		return errors.New("worker pool is not ready.")
	}
	if threads < 1 {
		return errors.New("worker pool needs at least has one thread.")
	}

	for i := 0; i < threads; i++ {
		go wp.worker()
	}

	return nil
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
				return errors.New("send job timeout.")
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

	for job := range wp {
		if job == nil {
			continue
		}
		if job.timeout < 1 {
			job.output <- job.p(job.input)
		} else {
			tick := time.NewTicker(job.timeout)
			select {
			case <-tick.C:
				log.Println("job timeout in worker.")
			case job.output <- job.p(job.input):
				tick.Stop()
			}
		}
		close(job.output)
	}
}
