package queue

import (
	"container/heap"
	"errors"
	"sync"
	"time"
)

const (
	Init = iota
	Closed
)

type Queue struct {
	maxsize int
	state   int
	recv    chan struct{}
	sendq   chan Task
	closeCh chan struct{}
	once    sync.Once
	mutex   sync.Mutex
	heap    *TaskHeap
}

func NewQueue(maxsize int) *Queue {
	q := &Queue{
		maxsize: maxsize,
		state:   Init,
		recv:    make(chan struct{}),
		sendq:   make(chan Task),
		closeCh: make(chan struct{}),
		once:    sync.Once{},
		mutex:   sync.Mutex{},
		heap:    &TaskHeap{},
	}
	go q.loop()
	return q
}

func (q *Queue) tosend(task Task) (isClosed bool) {
	isClosed = false
	select {
	case q.sendq <- task:
	case <-q.closeCh:
		isClosed = true
	}
	return isClosed
}

func (q *Queue) loop() {
	ticker := &time.Ticker{}
	for {
		q.mutex.Lock()
		count := q.heap.Len()
		q.mutex.Unlock()

		if count == 0 {
			<-q.recv
		}

		q.mutex.Lock()
		task := q.heap.Pop().(Task)
		q.mutex.Unlock()

		now := time.Now()

		// 任务延迟未达到, 需要等待延迟到达或者有新的任务被添加到队列中
		if task.RunAt.After(now) {
			sub := task.RunAt.Sub(now)
			ticker.Reset(sub)

			select {
			case <-q.recv:
				q.mutex.Lock()
				q.heap.Push(task)
				q.mutex.Unlock()

			case <-ticker.C:
				if closed := q.tosend(task); closed {
					return
				}

			case <-q.closeCh:
				return
			}
		} else {
			// 当前任务可被 Pop 调用接收
			if closed := q.tosend(task); closed {
				return
			}
		}
	}
}

func (q *Queue) Add(t Task) error {
	if err := q.ifClosed(); err != nil {
		return err
	}
	q.mutex.Lock()
	heap.Push(q.heap, t)
	q.mutex.Unlock()
	// TODO: 需要写入不阻塞, 但每次写会覆盖
	go func() {
		select {
		case q.recv <- struct{}{}:
		case <-q.closeCh:
		}
	}()
	return nil
}

func (q *Queue) Pop(handler Handler) error {
	if err := q.ifClosed(); err != nil {
		return err
	}
	for {
		select {
		case <-q.closeCh:
			return q.ifClosed()
		case t := <-q.sendq:
			// TODO: recover 捕获 panic, 执行失败后要重新添加到队列
			handler(t)
			return nil
		}
	}
}

// 关闭队列, 阻止 Add 和 Pop
func (q *Queue) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if err := q.ifClosed(); err != nil {
		return err
	}
	close(q.closeCh)
	q.state = Closed
	return nil
}

// 安全关闭, 多次调用只执行一次
func (q *Queue) SafeClose() {
	q.once.Do(func() { q.Close() })
}

func (q *Queue) ifClosed() error {
	if q.state == Closed {
		return errors.New(QueueClosedError)
	}
	return nil
}
