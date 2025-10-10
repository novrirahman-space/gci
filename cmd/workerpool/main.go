package main

import (
	"context"
	"math/rand"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
	Work func(ctx context.Context) (any, error)
}

type Result struct {
	JobID int
	Out any
	Err error
}

func worker(ctx context.Context, workerID int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		select {
		case <-ctx.Done():
			fmt.Printf("[worker %d] context canceled, stop\n", workerID)
			return
		default:
			out, err := job.Work(ctx)
			fmt.Printf("[worker %d] processed job %d\n", workerID, job.ID)
			results <- Result{JobID: job.ID, Out: out, Err: err}
		}
	}
}

func main() {
	const (
		numWorkers = 5
		numJobs = 12
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(ctx, w, jobs, results, &wg)
	}

	go func ()  {
		for i := 1; i <= numJobs; i++ {
			i := i
			jobs <- Job{
				ID: i,
				Work: func(ctx context.Context) (any, error) {
					time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)
					return fmt.Sprintf("job-%d done", i), nil
				},
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	okCount := 0
	errCount := 0
	start := time.Now()
	for res := range results {
		if res.Err != nil {
			errCount++
			fmt.Printf("[job %d] error: %v\n", res.JobID, res.Err)
			continue
		}
		okCount++
		fmt.Printf("[job %d] %v\n", res.JobID, res.Out)
	}

	elapsed := time.Since(start)
	fmt.Printf("completed: %d ok, %d error, elapsed: %s\n", okCount, errCount, elapsed.Truncate(time.Millisecond))

}