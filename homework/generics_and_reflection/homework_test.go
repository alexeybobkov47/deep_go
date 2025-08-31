package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

/*
В домашнем задании нужно реализовать сериализацию объекта структуры данных в .properties формат с использованием тегов структур и рефлексии.

При сериализации данных в некоторых случаях могут возникнуть проблемы с пустыми значениями.
Например, если у вас есть структура, в которой некоторые поля могут быть не заполнены, и вы сериализуете, то в результате получится объект с пустыми полями.
Если это не является ожидаемым поведением, то в нашей реализации можно будет использовать тег omitempty, чтобы пропустить пустые поля при сериализации.

Стуктура для сериализации в .properties формат (поподробнее с .properties форматом можно ознакомиться здесь https://ru.wikipedia.org/wiki/Properties):
*/

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](s T) string {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	out := strings.Builder{}

	fields := reflect.VisibleFields(t)
	lines := make([]string, 0, len(fields))
	for _, f := range fields {
		properties, ok := f.Tag.Lookup("properties")
		if !ok {
			continue
		}

		value := v.FieldByIndex(f.Index)

		tags := strings.Split(properties, ",")

		if len(tags) > 1 && tags[1] == "omitempty" && value.IsZero() {
			continue
		}

		lines = append(lines, fmt.Sprintf("%s=%v", tags[0], value))
	}

	out.WriteString(strings.Join(lines, "\n"))

	return out.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
