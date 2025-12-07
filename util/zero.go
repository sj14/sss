package util

func IfNotZero[T comparable](value T) *T {
	var zero T
	if value != zero {
		return &value
	}
	return nil
}

// func SetIfNotZero[T comparable](dst **T, value T) {
// 	var zero T
// 	if value != zero {
// 		*dst = &value
// 	}
// }
