package util

func ReverseSlice[T any](in []T) []T {
	out := make([]T, 0, len(in))
	for i := len(in) - 1; i >= 0; i-- {
		out = append(out, in[i])
	}

	return out
}
