package util

func Ptr[T comparable](v T) *T {
	return &v
}
