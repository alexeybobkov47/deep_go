package main

import (
	"reflect"
	"testing"

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

type node struct {
	key   int
	value int

	l *node
	r *node
}

type OrderedMap struct {
	root *node
	size int
}

// создать упорядоченный словарь
func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

// добавить элемент в словарь
func (m *OrderedMap) Insert(key, value int) {
	m.root = m.root.insert(key, value)
	m.size += 1
}

func (n *node) insert(key, value int) *node {
	if n == nil {
		return &node{key: key, value: value}
	}

	if key > n.key {
		n.r = n.r.insert(key, value)
		return n
	}

	n.l = n.l.insert(key, value)
	return n
}

// удалить элемент из словари
func (m *OrderedMap) Erase(key int) {
	m.root = m.root.erase(key)
	if m.size > 0 {
		m.size -= 1
	}
}

func (n *node) erase(key int) *node {
	if n == nil {
		return nil
	}

	if key != n.key {
		if key > n.key {
			n.r = n.r.erase(key)
			return n
		}

		n.l = n.l.erase(key)
		return n
	}

	if n.l == nil && n.r == nil {
		return nil
	}

	if n.r == nil {
		return n.l
	}

	k, v := n.r.findMin()
	n.r = n.r.erase(k)
	n.key = k
	n.value = v

	return n
}

func (n *node) findMin() (key, value int) {
	if n == nil {

	}
	if n.r == nil {
		return n.key, n.value
	}

	return n.r.findMin()
}

// проверить существование элемента в словаре
func (m *OrderedMap) Contains(key int) bool {
	return m.root.contains(key)
}

func (n *node) contains(key int) bool {
	if n == nil {
		return false
	}

	if n.key == key {
		return true
	}

	if key > n.key {
		return n.r.contains(key)
	}

	return n.l.contains(key)
}

// получить количество элементов в словаре
func (m *OrderedMap) Size() int {
	return m.size
}

// применить функцию к каждому элементу словаря от меньшего к большему
func (m *OrderedMap) ForEach(action func(int, int)) {
	m.root.forEach(action)
}

func (n *node) forEach(action func(int, int)) {
	if n != nil {
		n.l.forEach(action)
		action(n.key, n.value)
		n.r.forEach(action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
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
