package errors

// Must ensures error is nil, otherwise panic.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
