package errors_test

import (
	"testing"

	"goeasy.dev/errors"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	testCases := []struct {
		desc   string
		err    error
		target error
		wanted bool
	}{
		{
			desc:   "Simple",
			err:    errors.New("test"),
			target: errors.New("test"),
			wanted: true,
		},
		{
			desc:   "SimpleFail",
			err:    errors.New("test"),
			target: errors.New("fail"),
			wanted: false,
		},
		{
			desc:   "WithType",
			err:    errors.New("fail").WithType(errors.New("test")),
			target: errors.New("test"),
			wanted: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.wanted, errors.Is(tC.err, tC.target))
		})
	}
}

func TestWrap(t *testing.T) {
	err := errors.New("base")
	err = errors.Wrap(err, "wrap1")
	err = errors.Wrap(err, "wrap2")

	assert.Equal(t, "wrap2: wrap1: base", err.Error())

	var err2 error
	err = errors.Wrap(err2, "wrapping nil error")
	assert.Equal(t, "wrapping nil error", err.Error())
}
