package container

import (
	"fmt"
	"reflect"
)

var ImplementsMode = false

var services = map[reflect.Type]interface{}{}

func Clear() {
	services = map[reflect.Type]interface{}{}
}

func Resolve[T any]() T {
	targetType := reflect.TypeOf(new(T))

	val, ok := services[targetType]
	if !ok && !ImplementsMode {
		panic(fmt.Sprintf("type %s not registered", targetType))
	} else if !ok {
		return handleImplementsMode[T]()
	}

	t, _ := resolve[T](val)
	return t
}

func Set[T any](service T) {
	services[reflect.TypeOf(new(T))] = service
}

func SetResolver[T any](resolver func() T) {
	services[reflect.TypeOf(new(T))] = resolver
}

func handleImplementsMode[T any]() T {
	iface := reflect.TypeOf((*T)(nil)).Elem()

	for setType, val := range services {
		if setType.Kind() == reflect.Ptr {
			setType = setType.Elem()
		}

		if setType.Implements(iface) {
			t, usedResolver := resolve[T](val)
			if !usedResolver {
				Set(t)
			}

			return t
		}
	}

	panic("type not registered")
}

func resolve[T any](set interface{}) (t T, usedResolver bool) {
	if t, ok := set.(T); ok {
		return t, false
	}

	if resolver, ok := set.(func() T); ok {
		t := resolver()
		Set(t)

		return t, true
	}

	panic("unable to resolve")
}
