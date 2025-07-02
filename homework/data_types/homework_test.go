package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func ToLittleEndian[T uint16 | uint32 | uint64](number T) (resp T) {
	inputPtr := unsafe.Pointer(&number)
	respPtr := unsafe.Pointer(&resp)
	n := int(unsafe.Sizeof(number))

	for i := 0; i < n; i++ {
		inputVal := *(*uint8)(unsafe.Add(inputPtr, i))
		ptr := (*uint8)(unsafe.Add(respPtr, (n-1)-i))
		*ptr = inputVal
	}

	return resp
}

func TestToLittleEndianUint32(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1 (all zero)": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2 (all one)": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3 (alternating bytes)": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4 (low half set)": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5 (byte sequence)": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"test case #6 (only msb set)": {
			number: 0x80000000,
			result: 0x00000080,
		},
		"test case #7 (only lsb set)": {
			number: 0x00000001,
			result: 0x01000000,
		},
		"test case #8 (alternating pattern)": {
			number: 0xAA55AA55,
			result: 0x55AA55AA,
		},
		"test case #9 (mirror pattern)": {
			number: 0x12345678,
			result: 0x78563412,
		},
		"test case #10 (max)": {
			number: 0x7FFFFFFF,
			result: 0xFFFFFF7F,
		},
		"test case #11 (all but one)": {
			number: 0xFFFFFFFE,
			result: 0xFEFFFFFF,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestToLittleEndianUint64(t *testing.T) {
	tests := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"test case #4": {
			number: 0x00000000FFFFFFFF,
			result: 0xFFFFFFFF00000000,
		},
		"test case #5": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
		"test case #6": { // only msb set
			number: 0x8000000000000000,
			result: 0x0000000000000080,
		},
		"test case #7": { // only lsb set
			number: 0x0000000000000001,
			result: 0x0100000000000000,
		},
		"test case #8": { // alternating pattern
			number: 0xAA55AA55AA55AA55,
			result: 0x55AA55AA55AA55AA,
		},
		"test case #9": { // mirror pattern
			number: 0x0123456789ABCDEF,
			result: 0xEFCDAB8967452301,
		},
		"test case #10": { // max
			number: 0x7FFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFF7F,
		},
		"test case #11": { // all but one
			number: 0xFFFFFFFFFFFFFFFE,
			result: 0xFEFFFFFFFFFFFFFF,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestToLittleEndianUint16(t *testing.T) {
	tests := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"test case #4": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #5": {
			number: 0x0102,
			result: 0x0201,
		},
		"test case #6": { // only msb set
			number: 0x8000,
			result: 0x0080,
		},
		"test case #7": { // only lsb set
			number: 0x0001,
			result: 0x0100,
		},
		"test case #8": { // alternating pattern
			number: 0xAA55,
			result: 0x55AA,
		},
		"test case #9": { // mirror pattern
			number: 0x1234,
			result: 0x3412,
		},
		"test case #10": { // max
			number: 0x7FFF,
			result: 0xFF7F,
		},
		"test case #11": { // all but one
			number: 0xFFFE,
			result: 0xFEFF,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
