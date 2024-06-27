//go:build unit_test

package structure_test

import (
	"testing"

	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	t.Parallel()

	expQueue := &structure.Queue[string]{}
	actQueue := structure.NewQueue[string]()
	assert.Equal(t, expQueue, actQueue)
}

func TestAddingAndRemovingItemToQueue(t *testing.T) {
	t.Parallel()

	expItem := "test"
	queue := structure.NewQueue[string]()
	queue.Add(expItem)

	actItem := queue.Drop()
	assert.Equal(t, expItem, actItem)
}

func TestQueueLen(t *testing.T) {
	t.Parallel()

	queue := structure.NewQueue[string]()
	queue.Add("testA")
	queue.Add("testB")
	queue.Add("testC")

	expLength := 3
	assert.Equal(t, expLength, queue.Len())
}

func TestFiFoProperty(t *testing.T) {
	t.Parallel()

	expItem := "testA"
	queue := structure.NewQueue[string]()
	queue.Add(expItem)
	queue.Add("testB")

	actItem := queue.Drop()
	assert.Equal(t, expItem, actItem)
}

func TestQueueDropOnZeroLength(t *testing.T) {
	t.Parallel()

	queue := structure.NewQueue[string]()
	actResult := queue.Drop()
	expResult := ""
	assert.Equal(t, expResult, actResult)
}
