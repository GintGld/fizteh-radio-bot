package slice

import (
	"slices"
	"strings"
)

type Stringer interface {
	String() string
}

func StringSlice[T Stringer](t []T) []string {
	s := make([]string, 0, len(t))

	for _, el := range t {
		s = append(s, el.String())
	}

	return s
}

func Join[T Stringer](t []T, sep string) string {
	return strings.Join(StringSlice(t), sep)
}

func Repeat[T any](t T, n int) []T {
	s := make([]T, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, t)
	}
	return s
}

// returns slice with elements
// where filter is true.
func Filter[T any](t []T, filter []bool) []T {
	if len(t) != len(filter) {
		panic("slice and its filter must be the same size")
	}

	v := make([]T, 0, len(t))

	for i, flag := range filter {
		if flag {
			v = append(v, t[i])
		}
	}
	return slices.Clip(v)
}
