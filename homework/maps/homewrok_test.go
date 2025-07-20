package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// go test -v homework_test.go

/*
Идея упорядоченного словаря заключается в том, что он будет реализован на основе бинарного дерева поиска (BST).
Дерево будет строиться только по ключам элементов, значения элементов при построении дерева не учитываются.
Элементы с одинаковыми ключами в упорядоченном словаре хранить нельзя.

Поподробнее с бинарными деревьями поиска можно познакомиться [здесь.](https://habr.com/ru/articles/65617/)
*/

type node[K comparable, V any] struct {
	key   K
	value V

	l, r *node[K, V]
}

type OrderedMap[K comparable, V any] struct {
	root *node[K, V]
	size int

	less func(a, b K) bool
}

// создать упорядоченный словарь
func NewOrderedMap[K comparable, V any](less func(a, b K) bool) OrderedMap[K, V] {
	return OrderedMap[K, V]{less: less}
}

// добавить элемент в словарь
func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.root = m.insert(m.root, key, value)
	m.size += 1
}

func (m *OrderedMap[K, V]) insert(n *node[K, V], key K, value V) *node[K, V] {
	if n == nil {
		return &node[K, V]{key: key, value: value}
	}

	if m.less(n.key, key) {
		n.r = m.insert(n.r, key, value)
		return n
	}

	n.l = m.insert(n.l, key, value)
	return n
}

// удалить элемент из словари
func (m *OrderedMap[K, V]) Erase(key K) {
	m.root = m.erase(m.root, key)
	if m.size > 0 {
		m.size -= 1
	}
}

func (m *OrderedMap[K, V]) erase(n *node[K, V], key K) *node[K, V] {
	if n == nil {
		return nil
	}

	if key != n.key {
		if m.less(n.key, key) {
			n.r = m.erase(n.r, key)
			return n
		}

		n.l = m.erase(n.l, key)
		return n
	}

	if n.l == nil && n.r == nil {
		return nil
	}

	if n.r == nil {
		return n.l
	}

	min := n.r.findMin()
	if min != nil {
		n.key = min.key
		n.value = min.value
		n.r = m.erase(n.r, min.key)
	}

	return n
}

func (n *node[K, V]) findMin() *node[K, V] {
	if n == nil {
		return nil
	}
	if n.l == nil {
		return n
	}

	return n.l.findMin()
}

// проверить существование элемента в словаре
func (m *OrderedMap[K, V]) Contains(key K) bool {
	return m.contains(m.root, key)
}

func (m *OrderedMap[K, V]) contains(n *node[K, V], key K) bool {
	if n == nil {
		return false
	}

	if n.key == key {
		return true
	}

	if m.less(n.key, key) {
		return m.contains(n.r, key)
	}

	return m.contains(n.l, key)
}

// получить количество элементов в словаре
func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

// применить функцию к каждому элементу словаря от меньшего к большему
func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	m.forEach(m.root, action)
}

func (m *OrderedMap[K, V]) forEach(n *node[K, V], action func(K, V)) {
	if n != nil {
		m.forEach(n.l, action)
		action(n.key, n.value)
		m.forEach(n.r, action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int](func(a, b int) bool {
		return a < b
	})
	data.Erase(55)
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(10))
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	assert.True(t, data.root.r.l.r.key == 14)
	assert.True(t, data.root.l.l.key == 2)

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))

	assert.False(t, data.Contains(15))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	require.Equal(t, expectedKeys, keys)
}

func TestCircularQueueStringKey(t *testing.T) {
	data := NewOrderedMap[time.Time, string](func(a, b time.Time) bool {
		return a.Before(b)
	})

	t1 := time.Date(2025, 01, 31, 0, 0, 0, 0, time.Local)
	data.Insert(t1, "0")
	data.Insert(t1.Add(-25*time.Hour), "-25")
	data.Insert(t1.Add(-20*time.Hour), "-20")
	data.Insert(t1.Add(-23*time.Hour), "-23")
	data.Insert(t1.Add(-22*time.Hour), "-22")
	assert.True(t, data.root.l.r.l.r.value == "-22")

	values := make([]string, 0, data.Size())
	expected := []string{"-25", "-23", "-22", "-20", "0"}
	data.ForEach(func(_ time.Time, v string) {
		values = append(values, v)
	})
	require.Equal(t, expected, values)

	data.Erase(t1.Add(-25 * time.Hour))
	require.Equal(t, data.root.l.value, "-23")
	require.Equal(t, data.root.l.r.value, "-20")
	require.Equal(t, data.root.l.r.l.value, "-22")

	data.Erase(t1)
	require.Equal(t, data.root.value, "-23")

	data.Insert(t1.Add(1*time.Hour), "1")
	require.Equal(t, data.root.r.r.value, "1")

	values = make([]string, 0, data.Size())
	expected = []string{"-23", "-22", "-20", "1"}
	data.ForEach(func(_ time.Time, v string) {
		values = append(values, v)
	})
	require.Equal(t, expected, values)
}
