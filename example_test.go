package runner_test

import (
	"errors"
	"fmt"
	"sync"
	"time"

	runner "github.com/sotirispl/job-runner"
)

func Example_NoError() {
	mu := &sync.Mutex{}
	counter := 1
	printOutSmth := func() error {
		mu.Lock()
		defer mu.Unlock()
		fmt.Printf("executed: %d\n", counter)
		counter++
		return nil
	}

	runner := runner.New(2, 3*time.Second, printOutSmth)
	runner.Start()
	time.Sleep(4 * time.Second)
	runner.Stop()
	// Output: job 0 starting
	// job 1 starting
	// executed: 1
	// executed: 2
	// executed: 3
	// executed: 4
	// job 0 clossing
	// job 1 clossing
}

func Example_WithError() {
	mu := &sync.Mutex{}
	counter := 1
	printOutSmth := func() error {
		mu.Lock()
		defer mu.Unlock()
		if counter == 1 {
			return errors.New("some error")
		}
		fmt.Printf("executed: %d\n", counter)
		counter++
		return nil
	}

	runner := runner.New(1, 3*time.Second, printOutSmth)
	runner.Start()
	time.Sleep(1 * time.Second)
	runner.Stop()
	// Output: job 0 starting
	// job 0 got error: some error
	// job 0 clossing
}
