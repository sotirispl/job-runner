## Description

Runner that spawns new goroutines to call a function. The goroutines run endlessly until the stop function is called. The number of goroutines to be created is required along with the functionality that each routine will execute. Also, an interval between each execution must be defined.

[Documentation](http://godoc.org/github.com/sotirispl/job-runner)

## Usage

To start the job-runner just define the callback function to be executed in every process.

```
callback := func() error {
		fmt.Println("executed")
		return nil
	}

count := 2 // the number of processes that run simultaneously
interval := 3*time.Second // the time to wait for a process to be executed
runner := runner.New(count, interval, callback)
runner.Start()
```

You can stop the job-runner with:

```
runner.Stop()
```
