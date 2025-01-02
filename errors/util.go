package errors

import "errors"

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func IsNotifiable(err error) bool {
	var uerr *Error
	if !As(err, &uerr) {
		return true
	}

	return uerr.IsNotifiable()
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
