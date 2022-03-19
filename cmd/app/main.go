package main

import (
	"sync"

	"github.com/qsoyq/goqueue/pkg/queue"
)

func main() {
	q := queue.NewQueue()
	wg := sync.WaitGroup{}

	n := 10

	for i := 0; i < n; i++ {
		wg.Add(1)
		q.Add(i, 0)
	}

	go func() {
		for i := 0; i < n; i++ {
			q.Pop(func(t queue.Task) error {
				wg.Done()
				return nil
			})
		}
	}()
	wg.Wait()
}
