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
