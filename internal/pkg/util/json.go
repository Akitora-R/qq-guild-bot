package util

import (
	"encoding/json"
	"io"
)

func Stringify(i any) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func MustParse[T any](data []byte) *T {
	t := new(T)
	if err := json.Unmarshal(data, t); err != nil {
		panic(err)
	}
	return t
}

func MustParseArray[T any](data []byte) []T {
	var arr []T
	if err := json.Unmarshal(data, &arr); err != nil {
		panic(err)
	}
	return arr
}

func MustParseReader[T any](reader io.Reader) *T {
	all, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return MustParse[T](all)
}

func MustParseArrayReader[T any](reader io.Reader) []T {
	all, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return MustParseArray[T](all)
}
