package queue

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	DefaultMaxDuration = time.Second * 42
)

type Task struct {
	Attempts int
	Duration time.Duration
	RunAt    time.Time
	Data     interface{}
	backoff  backoff.BackOff
}

func (t *Task) Incr() {
	t.Attempts++
	now := time.Now()
	duration := t.nextBackOff()
	runAt := now.Add(duration)
	t.Duration = duration
	t.RunAt = runAt
}

func NewTask(data interface{}, duration time.Duration) Task {
	now := time.Now()
	runAt := now.Add(duration)
	t := Task{
		Attempts: 0,
		Duration: duration,
		RunAt:    runAt,
		Data:     data,
		backoff:  backoff.NewExponentialBackOff(),
	}
	t.backoff.(*backoff.ExponentialBackOff).MaxElapsedTime = 0
	return t
}

func (t *Task) nextBackOff() time.Duration {
	// backoff.ExponentialBackOff 并不严格单调递增, 而是按照因子在 interval 范围波动
	// 为了保证单调递增, 重置 InitialInterval 为 2 倍前值, 保证最小范围扔大于前值
	if t.Duration == 0 {
		t.backoff.(*backoff.ExponentialBackOff).InitialInterval = backoff.DefaultInitialInterval
	} else {
		t.backoff.(*backoff.ExponentialBackOff).InitialInterval = t.Duration * 2
	}
	t.backoff.Reset()
	duration := t.backoff.NextBackOff()

	// backoff.ExponentialBackOff 的 MaxInterval 只作为 currentInterval 的增长限制, 而不是实际延迟的上限
	if duration > DefaultMaxDuration {
		return DefaultMaxDuration
	}
	return duration
}
