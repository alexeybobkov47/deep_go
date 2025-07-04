package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type allowed interface {
	int | int8 | int16 | int32 | int64
}
type CircularQueue[T allowed] struct {
	len    int
	first  int
	last   int
	values []T
}

// создать очередь с определенным размером буффера
func NewCircularQueue[T allowed](size int) CircularQueue[T] {
	return CircularQueue[T]{
		len:    size,
		first:  0,
		last:   size - 1,
		values: make([]T, size),
	}
}

func (q *CircularQueue[T]) nextLast() {
	next := q.last + 1
	if next == q.len {
		next = 0
	}

	q.last = next
}
func (q *CircularQueue[T]) nextFirst() {
	next := q.first + 1
	if next == q.len {
		next = 0
	}

	q.first = next
}

// добавить значение в конец очереди (false, если очередь заполнена)
func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.nextLast()
	q.values[q.last] = value

	return true
}

// удалить значение из начала очереди (false, если очередь пустая)
func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	q.values[q.first] = 0
	q.nextFirst()

	return true
}

// получить значение из начала очереди (-1, если очередь пустая)
func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.first]
}

// получить значение из конца очереди (-1, если очередь пустая)
func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.last]
}

// проверить пустая ли очередь
func (q *CircularQueue[T]) Empty() bool {
	for _, v := range q.values {
		if v != 0 {
			return false
		}
	}

	return true
}

// проверить заполнена ли очередь
func (q *CircularQueue[T]) Full() bool {
	for _, v := range q.values {
		if v == 0 {
			return false
		}
	}

	return true

}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	//

	queue64 := NewCircularQueue[int64](3)
	assert.True(t, queue64.Push(1))
	assert.True(t, queue64.Push(2))
	assert.True(t, queue64.Push(3))
	assert.True(t, queue64.Full())
	assert.Equal(t, int64(1), queue64.Front())
	assert.True(t, queue64.Pop())
	assert.Equal(t, int64(2), queue64.Front())
	assert.True(t, queue64.Pop())
	assert.Equal(t, int64(3), queue64.Front())
	assert.True(t, queue64.Pop())
	assert.Equal(t, int64(-1), queue64.Front())
	assert.False(t, queue64.Pop())
	assert.True(t, queue64.Empty())

}
