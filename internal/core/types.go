package core

func PTR[T any](val T) *T {
	return &val
}
