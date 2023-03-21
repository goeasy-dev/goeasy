package util

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

const (
	stackOffset  = 2
	callerOffset = 2
	stackDepth   = 50

	callerPackageIndex = 0
	callerTypeIndex    = 1
	callerNameIndex    = 2
)

type CallStack []runtime.Frame

func (c CallStack) PrettyPrint(w io.Writer) {
	for _, frame := range c {
		fmt.Fprintln(w, frame.Function)
		fmt.Fprintf(w, "\t%s:%d\n", frame.File, frame.Line)
	}
}

func GetStack(offset ...int) CallStack {
	o := stackOffset
	if len(offset) >= 1 {
		o += offset[0]
	}

	callers := make([]uintptr, stackDepth)
	entryCount := runtime.Callers(o, callers)

	frames := runtime.CallersFrames(callers)
	out := make(CallStack, 0, entryCount)
	for {
		frame, more := frames.Next()
		out = append(out, frame)

		if !more {
			break
		}
	}

	return out
}

type Caller struct {
	Package string
	Type    string
	Name    string
	File    string
	Line    int
}

func GetCaller() Caller {
	pc, _, _, ok := runtime.Caller(callerOffset)
	if !ok {
		return Caller{}
	}

	return CallerFromFunc(runtime.FuncForPC(pc))
}

func CallerFromFunc(f *runtime.Func) Caller {
	if f == nil {
		return Caller{}
	}

	parts := strings.Split(f.Name(), ".")
	if len(parts) > 3 {
		parts[1] = fmt.Sprintf("%s.%s", parts[0], parts[1])
		parts = parts[1:]
	} else if strings.ContainsRune(parts[1], '/') {
		parts[0] = fmt.Sprintf("%s.%s", parts[0], parts[1])
		parts[1] = "unknown"
	}

	file, line := f.FileLine(f.Entry())

	c := Caller{
		Name:    parts[callerNameIndex],
		Type:    parts[callerTypeIndex],
		Package: parts[callerPackageIndex],
		File:    file,
		Line:    line,
	}

	return c
}
