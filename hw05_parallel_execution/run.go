package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrGoroutinesCount     = errors.New("goroutines count must be greater than 0")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrGoroutinesCount
	}

	var errors int64
	var shouldTrackErrors bool

	wg := &sync.WaitGroup{}
	shouldTrackErrors = m > 0
	tasksChan := make(chan Task)

	worker := func(tasksChan chan Task) {
		defer wg.Done()

		for task := range tasksChan {
			if err := task(); err != nil {
				atomic.AddInt64(&errors, 1)
			}
		}
	}

	wg.Add(n)

	for i := 0; i < n; i++ {
		go worker(tasksChan)
	}

	for _, task := range tasks {
		if !shouldTrackErrors {
			tasksChan <- task
			continue
		}

		errors := int(atomic.LoadInt64(&errors))

		if errors >= m {
			break
		}

		tasksChan <- task
	}

	close(tasksChan)

	wg.Wait()

	if shouldTrackErrors && int(errors) >= m {
		return ErrErrorsLimitExceeded
	}

	// Place your code here.
	return nil
}
