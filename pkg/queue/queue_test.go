package queue

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPopAfterClose(t *testing.T) {

	q := NewQueue()
	q.SafeClose()
	err := q.Pop(func(t Task) error {
		return nil
	})
	assert.EqualError(t, err, QueueClosedError)

}

func TestWorkflow(t *testing.T) {
	wg := sync.WaitGroup{}

	handler := func(t Task) error {
		time.Sleep(time.Duration(time.Microsecond * 20))
		if rand.Float32() >= 0.65 {
			return errors.New("handler error")
		}
		return nil
	}
	start := make(chan struct{})
	queue := NewQueue()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		queue.Add(nil, uint(i))
		go func() {
			<-start
			for {
				if err := queue.Pop(handler); err == nil {
					wg.Done()
					break
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}
