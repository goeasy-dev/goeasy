package container

import (
	"reflect"
)

var services = map[reflect.Type]interface{}{}

func Resolve[T any]() T {
	val := services[reflect.TypeOf(new(T))]
	if t, ok := val.(T); ok {
		return t
	}

	if resolver, ok := val.(func() T); ok {
		t := resolver()
		Set(t)

		return t
	}

	panic("type not registered")
}

func Set[T any](service T) {
	services[reflect.TypeOf(new(T))] = service
}

func SetResolver[T any](resolver func() T) {
	services[reflect.TypeOf(new(T))] = resolver
}
