package queue

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
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

func TestCloseWhilePop(t *testing.T) {
	q := NewQueue()
	wg := sync.WaitGroup{}
	wg.Add(1)

	start := make(chan struct{})

	q.Add(nil, 0)
	go func() {
		close(start)
		for {
			err := q.Pop(func(t Task) error {
				return errors.New("")
			})
			if err.Error() == QueueClosedError {
				break
			}
			time.Sleep(time.Duration(time.Second))
		}
		wg.Done()
	}()

	<-start
	go func() {
		time.Sleep(time.Second * 5)
		q.SafeClose()
	}()
	wg.Wait()
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

func TestBackoffRetryDuration(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	queue := NewQueue()
	queue.Add(nil, 0)
	go func() {
		var lastDuration time.Duration
		// 第一次 Pop 失败重新添加到队列后, Task.Duration 才会被重置为每次 bakcoff 对应的 Duration
		queue.Pop(func(t Task) error {
			return errors.New("backoff retry test")
		})
		queue.Pop(func(t Task) error {
			lastDuration = t.Duration
			return errors.New("backoff retry test")
		})

		for i := 0; i < 10; i++ {
			err := queue.Pop(func(task Task) error {
				t.Logf("Duration: %d, MaxInterval: %d", task.Duration.Milliseconds(), DefaultMaxDuration.Milliseconds())
				if bf, ok := task.backoff.(*backoff.ExponentialBackOff); ok {
					if task.Duration > bf.MaxInterval {
						t.Errorf("task druation biger than max interval")
						return errors.New("backoff duration error")
					}
				}

				if lastDuration >= task.Duration {
					if lastDuration != DefaultMaxDuration {
						t.Logf("attemps: %d, last duration: %d, current duration: %d", task.Attempts, lastDuration.Milliseconds(), task.Duration.Milliseconds())
						return errors.New("backoff duration error")
					}
				}
				lastDuration = task.Duration
				return errors.New("backoff retry test")
			})
			if err.Error() == "backoff duration error" {
				wg.Done()
				return
			}
		}
		queue.Pop(func(task Task) error {
			wg.Done()
			return nil
		})

	}()
	wg.Wait()

}
