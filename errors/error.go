package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"goeasy.dev/observability/log"
	"goeasy.dev/util"
)

const (
	maxStackDepth       = 50
	newErrorStackOffset = 4
)

type ErrorFlag int

const (
	Notifiable ErrorFlag = 1 << iota
	DataLoss
)

type Error struct {
	data       map[string]interface{}
	messages   []string
	baseError  error
	errorType  error
	flags      ErrorFlag
	stackTrace []runtime.Frame
	reason     string
}

func New(msg string, args ...interface{}) *Error {
	return new(fmt.Errorf(msg, args...), newErrorStackOffset, Notifiable)
}

func Wrap(err error, args ...interface{}) *Error {
	if err == nil {
		caller := util.GetCaller()
		log.Warn("nil error wrapped in %s:%d", caller.File, caller.Line)

		msg := "nil error"
		if len(args) > 0 {
			tmp, ok := args[0].(string)
			if ok {
				msg = tmp
			}
		}

		return New(msg, args[1:]...)
	}

	udoErr, ok := err.(*Error)
	if !ok {
		udoErr = new(err, newErrorStackOffset, Notifiable)
	}

	if len(args) > 0 {
		msg, ok := args[0].(string)
		if msg != "" && ok {
			udoErr.messages = append([]string{fmt.Sprintf(msg, args[1:]...)}, udoErr.messages...)
		}
	}

	return udoErr
}

func (e Error) Error() string {
	return strings.Join(
		append(e.messages, e.baseError.Error()),
		": ",
	)
}

func (e Error) Unwrap() error {
	base := e.baseError
	if e.errorType != nil {
		base = e.errorType
	}

	return base
}

func (e Error) Is(target error) bool {
	if tmp := errors.Unwrap(target); tmp != nil {
		target = tmp
	}

	var source error = e
	if tmp := errors.Unwrap(e); tmp != nil {
		source = tmp
	}

	return source.Error() == target.Error()
}

func (e Error) IsNotifiable() bool {
	return e.flags&Notifiable == Notifiable
}

func (e *Error) SetNotifiable(notifiable bool) {
	if !notifiable {
		e.flags = clearFlag(e.flags, Notifiable)
		return
	}

	e.flags = setFlag(e.flags, Notifiable)
}

func (e Error) IsDataLoss() bool {
	return e.flags&DataLoss == DataLoss
}

func (e *Error) SetDataLoss() {
	e.flags = setFlag(e.flags, DataLoss)
}

func (e Error) Stack() []string {
	stack := make([]string, 0, len(e.stackTrace))
	for _, frame := range e.stackTrace {
		stack = append(stack, fmt.Sprintf("%s:%d  -  %s", frame.File, frame.Line, frame.Function))
	}

	return stack
}

func (e Error) StackTrace() []runtime.Frame {
	return e.stackTrace
}

func (e Error) Data() map[string]interface{} {
	return e.data
}

func (e *Error) WithData(data map[string]interface{}) *Error {
	mergeMaps(e.data, data)

	return e
}

// WithReason sets a reason code on the error
func (e *Error) WithReason(reason string) *Error {
	e.reason = reason

	return e
}

// Reason returns a code that can be used to get additional details about the error
func (e Error) Reason() string {
	return e.reason
}

// SetType allows for overriding what the error will unwrap as, thus affecting `Is` comparisons.
// As an example, this allows for setting a type on a wrapped error returned from a service so it
// can be appropratly handled by calling functions.
func (e *Error) WithType(err error) *Error {
	e.errorType = err

	return e
}

func new(err error, offset int, flags ErrorFlag) *Error {
	return &Error{
		messages:   []string{},
		baseError:  err,
		data:       map[string]interface{}{},
		stackTrace: getStack(offset),
		flags:      flags,
	}
}

func getStack(skipFrames int) []runtime.Frame {
	out := make([]runtime.Frame, 0, maxStackDepth)
	callers := make([]uintptr, maxStackDepth)

	runtime.Callers(skipFrames, callers)
	frames := runtime.CallersFrames(callers)
	more := true
	var next runtime.Frame
	for more {
		next, more = frames.Next()
		out = append(out, next)
	}

	return out
}

func mergeMaps(base, top map[string]interface{}) {
	for k, v := range top {
		base[k] = v
	}
}

func setFlag(flags ErrorFlag, flag ErrorFlag) ErrorFlag {
	return flags | flag
}
func clearFlag(flags ErrorFlag, flag ErrorFlag) ErrorFlag {
	return flags &^ flag
}
