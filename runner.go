// Package runner spawns new goroutines to call a function. The goroutines run endlessly until the stop function is called. The number of goroutines to be created is required along with the functionality that each routine will execute. Also, an interval between each execution must be defined.
package runner

import (
	"fmt"
	"sync"
	"time"
)

// Runner orchestrates and handles jobs
type Runner struct {
	jobs     []*job
	errors   []chan error
	callback func() error
	quit     chan bool
	mu       *sync.Mutex
}

// job defines a channel to run on something
type job struct {
	id       int
	interval time.Duration
	quit     chan bool
	callback func() error
	mu       *sync.Mutex
}

// New creates a runner with its jobs
func New(count int, interval time.Duration, callback func() error) *Runner {
	jobs := make([]*job, count)
	for i := 0; i < count; i++ {
		job := newJob(i, interval, callback)
		jobs[i] = job
	}
	return &Runner{
		jobs:   jobs,
		errors: make([]chan error, count),
		quit:   make(chan bool),
		mu:     &sync.Mutex{},
	}
}

// Start executes all the jobs declared
func (r *Runner) Start() {
	wg := &sync.WaitGroup{}
	for _, j := range r.jobs {
		fmt.Printf("job %d starting\n", j.id)
		go j.execute(wg)
	}

	wg.Wait()
}

// Stop gracefully sends the quit signal to jobs
func (r *Runner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, j := range r.jobs {
		j.stop()
	}
}

func newJob(id int, interval time.Duration, callback func() error) *job {
	return &job{
		id:       id,
		interval: interval,
		quit:     make(chan bool),
		callback: callback,
		mu:       &sync.Mutex{},
	}
}

func (j *job) execute(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
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
	j.mu.Lock()
	defer j.mu.Unlock()
	fmt.Printf("job %d clossing\n", j.id)
	j.quit <- true
}
