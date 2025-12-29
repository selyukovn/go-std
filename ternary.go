package std

func Ternary[T any](condition bool, rTrue T, rFalse T) T {
	if condition {
		return rTrue
	} else {
		return rFalse
	}
}
