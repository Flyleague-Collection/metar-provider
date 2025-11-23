// Package utils
package utils

func ReverseForEach[T any](slice []T, f func(index int, value T)) {
	for i := len(slice) - 1; i >= 0; i-- {
		f(i, slice[i])
	}
}

func Any[T comparable](src []T, comparator func(element T) bool) bool {
	for _, v := range src {
		if comparator(v) {
			return true
		}
	}
	return false
}

func Find[T comparable](src []T, comparator func(element T) bool) T {
	for _, v := range src {
		if comparator(v) {
			return v
		}
	}
	var zero T
	return zero
}

func Filter[T comparable](src []T, filter func(element T) bool) (result []T) {
	result = make([]T, 0, len(src))
	for _, v := range src {
		if filter(v) {
			result = append(result, v)
		}
	}
	return
}

func Map[T comparable](src []T, mapper func(element T) T) {
	for i := range src {
		src[i] = mapper(src[i])
	}
}

func ForEach[T comparable](src []T, callback func(index int, element T)) {
	for i, v := range src {
		callback(i, v)
	}
}
