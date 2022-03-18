package queue

import (
	"errors"
	"time"
)

type Task struct {
	Attempts int
	Duration time.Duration
	RunAt    time.Time
	Data     interface{}
}

func (t Task) NextRunAt() (*time.Time, error) {
	t.Attempts++
	if t.Attempts >= MaxAttempts {
		return nil, errors.New("max attempts")
	}

	return nil, nil
}
