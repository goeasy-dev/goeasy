package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"goeasy.dev/errors"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc   string
		err    error
		target error
	}{
		{
			desc:   "CommonHere",
			err:    errors.ErrForbidden.Here(),
			target: errors.ErrForbidden,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			fmt.Println(errors.Wrap(tC.err).Stack())
			assert.True(t, errors.Is(tC.err, tC.target))
			t.Fail()
		})
	}
}
