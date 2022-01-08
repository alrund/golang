package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorCounter struct {
	mu sync.Mutex
	i  int
}

func (c *errorCounter) inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.i++
}

func (c *errorCounter) moreThen(limit int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.i >= limit
}

func taskChannel(tasks []Task) <-chan Task {
	taskNum := len(tasks)
	ch := make(chan Task, taskNum)
	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	return ch
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errCounter := errorCounter{}
	taskCh := taskChannel(tasks)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(errCounter *errorCounter, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				task, ok := <-taskCh
				if !ok {
					break
				}
				result := task()
				if result != nil {
					errCounter.inc()
				}
				if errCounter.moreThen(m) {
					break
				}
			}
		}(&errCounter, &wg)
	}

	wg.Wait()

	if errCounter.moreThen(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
