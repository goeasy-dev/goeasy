package util_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"goeasy.dev/util"
)

func stackRetriever(offset ...int) util.CallStack {
	return util.GetStack(offset...)
}

func callerWrapper() util.Caller {
	return util.GetCaller()
}

func TestCallStack(t *testing.T) {
	testCases := []struct {
		desc         string
		f            func(...int) util.CallStack
		offset       int
		expectedFunc string
	}{
		{
			desc:         "NoOffset",
			f:            stackRetriever,
			expectedFunc: "goeasy.dev/util_test.stackRetriever",
		},
		{
			desc:         "Offset",
			f:            stackRetriever,
			offset:       1,
			expectedFunc: "goeasy.dev/util_test.TestCallStack.func1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			stack := tC.f(tC.offset)
			stack.PrettyPrint(os.Stderr)
			assert.Equal(t, tC.expectedFunc, stack[0].Function)
		})
	}
}

func BenchmarkGetStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		util.GetStack()
	}
}

type callerWrapperStruct struct{}

func (w callerWrapperStruct) caller() util.Caller {
	return callerWrapper()
}

func TestGetCaller(t *testing.T) {
	f := callerWrapper()
	assert.Equal(t, "TestGetCaller", f.Name)

	ff := (callerWrapperStruct{}).caller()
	assert.Equal(t, "callerWrapperStruct", ff.Type)
}

func BenchmarkGetCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		util.GetCaller()
	}
}
