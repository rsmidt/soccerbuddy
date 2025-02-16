package core

func Map[T any, R any](arr []T, f func(T) R) []R {
	res := make([]R, len(arr))
	for i, e := range arr {
		res[i] = f(e)
	}
	return res
}

func MapPtr[T any, R any](arr []T, f func(*T) R) []R {
	res := make([]R, len(arr))
	for i, e := range arr {
		res[i] = f(&e)
	}
	return res
}
