package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorCounter := 0
	wg := sync.WaitGroup{}
	wg.Add(n)
	ch := Tasker(tasks)
	mu := sync.Mutex{}

	for i := 0; i < n; i++ {
		go func(ch <-chan Task, errorCounter *int, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			for {
				task, ok := <-ch
				if !ok {
					break
				}
				result := task()

				mu.Lock()
				if result != nil {
					*errorCounter++
				}
				if *errorCounter >= m {
					mu.Unlock()
					break
				}
				mu.Unlock()
			}
		}(ch, &errorCounter, &wg, &mu)
	}

	wg.Wait()

	if errorCounter >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func Tasker(tasks []Task) <-chan Task {
	taskNum := len(tasks)
	ch := make(chan Task, taskNum)
	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	return ch
}
