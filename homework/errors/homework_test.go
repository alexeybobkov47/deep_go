package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%d errors occured:\n", len(e.errors)))
	for _, err := range e.errors {
		b.WriteString("\t* " + err.Error())
	}
	b.WriteString("\n")
	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	switch t := err.(type) {
	case *MultiError:
		t.errors = append(t.errors, errs...)
		return t
	}

	m := MultiError{errors: make([]error, 0)}
	m.errors = append(m.errors, errs...)

	return &m
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	var ErrNotFound = errors.New("not found")
	err = Append(err, ErrNotFound)
	assert.True(t, errors.Is(err, ErrNotFound))

	err = Append(err, TestErr{})
	assert.True(t, errors.As(err, &TestErr{}))

}

type chain []error

func (e chain) Error() string {
	return e[0].Error()
}

func (e *MultiError) Unwrap() error {
	if e == nil || len(e.errors) == 0 {
		return nil
	}

	return chain(e.errors)
}

func (e chain) Unwrap() error {
	if len(e) == 1 {
		return nil
	}
	return e[1:]
}

func (e chain) Is(tagret error) bool {
	if e[0] == tagret {
		return true
	}

	return false
}

func (e chain) As(target any) bool {
	return errors.As(e[0], target)
}

type TestErr struct {
	e error
}

func (t TestErr) Error() string {
	return ""
}
