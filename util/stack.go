package util

import (
	"fmt"
	"io"
	"regexp"
	"runtime"
)

const (
	stackOffset  = 3
	callerOffset = 2
	stackDepth   = 50

	callerNameRegex    = `^(?P<Package>.*\/[^.]*)\.(?:(?P<Type>.*)\.)?(?P<Name>.*)$`
	callerPackageIndex = 1
	callerTypeIndex    = 2
	callerNameIndex    = 3
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

	return getStack(o, stackDepth)
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

var callerNameRegexExp = regexp.MustCompile(callerNameRegex)

func CallerFromFunc(f *runtime.Func) Caller {
	if f == nil {
		return Caller{}
	}

	matches := callerNameRegexExp.FindStringSubmatch(f.Name())
	if matches == nil {
		matches = make([]string, 4)
	}

	file, line := f.FileLine(f.Entry())

	c := Caller{
		Name:    matches[callerNameIndex],
		Type:    matches[callerTypeIndex],
		Package: matches[callerPackageIndex],
		File:    file,
		Line:    line,
	}

	return c
}

func getStack(offset, depth int) CallStack {
	callers := make([]uintptr, depth)
	entryCount := runtime.Callers(offset, callers)

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
