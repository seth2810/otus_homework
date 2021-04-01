package hw05parallelexecution

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.ErrorIs(t, err, ErrErrorsLimitExceeded, "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("no tasks", func(t *testing.T) {
		tasksCount := 0
		workersCount := 2
		tasks := make([]Task, tasksCount)

		require.NoError(t, Run(tasks, workersCount, 1))
	})

	t.Run("ignoring errors", func(t *testing.T) {
		tasksCount := 10
		workersCount := 2
		var runTasksCount int32
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		require.NoError(t, Run(tasks, workersCount, -1))
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

		atomic.StoreInt32(&runTasksCount, 0)

		require.NoError(t, Run(tasks, workersCount, 0))
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("negative goroutines count", func(t *testing.T) {
		tasksCount := 10
		workersCount := -1
		var runTasksCount int32
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		require.Equal(t, runTasksCount, int32(0), "some of tasks were completed")
		require.EqualError(t, Run(tasks, workersCount, 1), "goroutines count must be greater than 0")
	})

	t.Run("concurrency without sleeps", func(t *testing.T) {
		tasksCount := 50
		workersCount := 5
		maxErrorsCount := 1
		completeChan := make(chan error)
		tasks := make([]Task, 0, tasksCount)

		var err error
		var runTasksCount int32
		var totalTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			totalTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		go func() {
			defer close(completeChan)

			completeChan <- Run(tasks, workersCount, maxErrorsCount)
		}()

		require.Eventually(t, func() bool {
			select {
			case err = <-completeChan:
				return true
			default:
				return false
			}
		}, totalTime/2, time.Millisecond, "tasks were run sequentially?")

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}
