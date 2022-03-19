package main

import (
	"fmt"
	"sync"

	"github.com/qsoyq/goqueue/pkg/queue"
)

func main() {
	q := queue.NewQueue()
	wg := sync.WaitGroup{}

	n := 10

	for i := 0; i < n; i++ {
		wg.Add(1)
		q.Add(i, uint(i))
	}

	go func() {
		for i := 0; i < n; i++ {
			q.Pop(func(t queue.Task) error {
				fmt.Printf("task %+v, %s\n", t.Data, t.RunAt.String())
				wg.Done()
				return nil
			})
		}
	}()
	wg.Wait()
}
