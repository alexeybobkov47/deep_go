package main

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

/*
Идея подхода **copy-on-write** заключается в том, что при чтении данных используется общая копия данных буффера,
	но в случае изменения данных — создается новая копия данных буффера.
Для реализации такого подхода можно использовать разделяемый счетчик ссылок - если при изменении данных буффера кто-то еще ссылается на этот буффер,
	то нужно будет сначала произвести копию данных буффера, изменить счетчик ссылок и
	только затем произвести изменение (*если никто не ссылается на буффер, то копировать данные буффера не нужно при изменении данных*).

Дополнительно еще нужно реализовать метод конвертации данных буффера в строку без копирования и дополнительного выделения памяти.
*/

// Предположим, что все будут производить копирование
// буффера только с использованием метода Clone()
type COWBuffer struct {
	data []byte
	refs *int
	size int
}

// создать буффер с определенными данными
func NewCOWBuffer(data []byte) COWBuffer {
	// d := make([]byte, len(data))
	// copy(d, data)
	init := 0
	b := COWBuffer{
		data: data,
		size: len(data),
		refs: &init,
	}

	runtime.SetFinalizer(&b, func(b *COWBuffer) {
		fmt.Println("finalizer close")
		b.Close()
	})

	return b
}

// создать новую копию буфера
func (b *COWBuffer) Clone() COWBuffer {
	*b.refs = *b.refs + 1
	return *b
}

// перестать использовать копию буффера
func (b *COWBuffer) Close() {
	if *b.refs == 0 {
		return
	}

	*b.refs = *b.refs - 1

	d := make([]byte, len(b.data))
	copy(d, b.data)
	init := 0

	b.data = d
	b.refs = &init

}

// изменить определенный байт в буффере
func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index >= b.size {
		return false
	}

	if *b.refs > 0 {
		b.Close()
	}

	b.data[index] = value

	return true
}

// сконвертировать буффер в строку
func (b *COWBuffer) String() string {
	if len(b.data) > 0 {
		return unsafe.String(unsafe.SliceData(b.data), len(b.data))
	}

	return ""
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()
	assert.Equal(t, 2, *buffer.refs)

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	runtime.GC()
	time.Sleep(1 * time.Second)

	assert.Equal(t, 0, *copy2.refs)
	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))
}
