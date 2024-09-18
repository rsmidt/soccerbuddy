package core

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
