package util

func SetIfNotZero[T comparable](dst **T, value T) {
	var zero T
	if value != zero {
		*dst = &value
	}
}
