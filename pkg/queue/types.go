package queue

const (
	MaxAttempts = 65535
)

const (
	QueueClosedError = "Queue Already Closed"
)

type TaskQueue interface {
	Add(t *Task) error
	Pop(handler Handler) error
	Close() error
	SafeClose() error
}

type Handler func(t Task) error

type TaskInterface interface {
	Incr()
}
