package queue

import (
	"math/rand"
	"time"
)

type Task struct {
	Attempts int
	Duration time.Duration
	RunAt    time.Time
	Data     interface{}
}

func (t *Task) Incr() {
	t.Attempts++
	now := time.Now()
	// TODO:　Add backoff
	t.RunAt = now.Add(time.Duration(time.Second * time.Duration(rand.Intn(3)+1)))
}

func NewTask(data interface{}, duration time.Duration) Task {
	now := time.Now()
	runAt := now.Add(duration)
	return Task{
		Attempts: 0,
		Duration: duration,
		RunAt:    runAt,
		Data:     data,
	}
}
