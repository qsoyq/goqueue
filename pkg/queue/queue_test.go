package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopAfterClose(t *testing.T) {

	q := NewQueue(10)
	q.SafeClose()
	err := q.Pop(nil)
	assert.EqualError(t, err, QueueClosedError)

}
