package core

import "slices"

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2[T any](ret T, err error) T {
	if err != nil {
		panic(err)
	}
	return ret
}

func RemoveNils[T any](s []*T) []*T {
	return slices.DeleteFunc(s, func(t *T) bool {
		return t == nil
	})
}
