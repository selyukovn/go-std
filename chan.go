package std

func ChanReadAll[T any](ch <-chan T) []T {
	var s []T

	for e := range ch {
		s = append(s, e)
	}

	return s
}
