// Package runner spawns new goroutines to call a function. The goroutines run endlessly until the stop function is called. The number of goroutines to be created is required along with the functionality that each routine will execute. Also, an interval between each execution must be defined.
package runner

import (
	"fmt"
	"sync"
	"time"
)

// Runner orchestrates and handles jobs.
type Runner struct {
	jobs     []*job
	callback func() error
	quit     chan struct{}
}

// job defines a channel to run on something.
type job struct {
	id       int
	interval time.Duration
	quit     chan struct{}
	callback func() error
}

// New creates a runner with its jobs.
func New(count int, interval time.Duration, callback func() error) *Runner {
	jobs := make([]*job, count)
	for i := 0; i < count; i++ {
		job := newJob(i, interval, callback)
		jobs[i] = job
	}
	return &Runner{
		jobs: jobs,
		quit: make(chan struct{}),
	}
}

// Start executes all the jobs declared.
func (r *Runner) Start() {
	go func() {
		wg := &sync.WaitGroup{}
		for _, j := range r.jobs {
			wg.Add(1)
			fmt.Printf("job %d starting\n", j.id)
			go func(j *job, wg *sync.WaitGroup) {
				j.execute()
				wg.Done()
			}(j, wg)
		}

		wg.Wait()
	}()
}

// Stop gracefully sends the quit signal to jobs.
func (r *Runner) Stop() {
	for _, j := range r.jobs {
		j.stop()
	}
}

func newJob(id int, interval time.Duration, callback func() error) *job {
	return &job{
		id:       id,
		interval: interval,
		quit:     make(chan struct{}),
		callback: callback,
	}
}

func (j *job) execute() {
	for {
		err := j.callback()
		if err != nil {
			fmt.Printf("job %d got error: %v\n", j.id, err)
		}
		select {
		case <-j.quit:
			return
		case <-time.After(j.interval):
			continue
		}
	}
}

func (j *job) stop() {
	fmt.Printf("job %d clossing\n", j.id)
	close(j.quit)
}
