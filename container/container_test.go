package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goeasy.dev/container"
)

type Foo interface {
	Foo() string
}

type FooBar interface {
	Foo

	Bar() string
}

type myService struct{}

func (m myService) Foo() string {
	return "foo"
}

func (m myService) Bar() string {
	return "bar"
}

func TestContainer(t *testing.T) {
	container.Clear()

	container.Set[FooBar](&myService{})
	container.SetResolver(func() Foo {
		return container.Resolve[FooBar]()
	})

	service := container.Resolve[FooBar]()
	assert.Equal(t, "foo", service.Foo())
	assert.Equal(t, "bar", service.Bar())

	fooService := container.Resolve[Foo]()
	assert.Equal(t, "foo", fooService.Foo())
}

func TestImplementsMode(t *testing.T) {
	container.Clear()
	container.ImplementsMode = true

	container.Set(&myService{})

	service := container.Resolve[FooBar]()
	assert.Equal(t, "foo", service.Foo())
	assert.Equal(t, "bar", service.Bar())

	fooService := container.Resolve[Foo]()
	assert.Equal(t, "foo", fooService.Foo())

	container.ImplementsMode = false
}

func BenchmarkResolve(b *testing.B) {
	container.Clear()

	container.Set[FooBar](&myService{})
	for i := 0; i < b.N; i++ {
		container.Resolve[FooBar]()
	}
}

func BenchmarkResolveImplements(b *testing.B) {
	container.Clear()
	container.ImplementsMode = true

	container.Set(&myService{})
	for i := 0; i < b.N; i++ {
		container.Resolve[FooBar]()
	}

	container.ImplementsMode = false
}
