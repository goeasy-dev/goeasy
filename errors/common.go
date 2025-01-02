package errors

import "errors"

const commonErrorStackOffset int = 4

type commonError struct {
	error
}

func (c commonError) Here() *Error {
	return new(c, commonErrorStackOffset, 0)
}

var (
	ErrBadRequest     commonError = commonError{errors.New("bad request")}
	ErrNotFound       commonError = commonError{errors.New("not found")}
	ErrInternal       commonError = commonError{errors.New("internal")}
	ErrUnauthorized   commonError = commonError{errors.New("unauthorized")}
	ErrForbidden      commonError = commonError{errors.New("forbidden")}
	ErrNotImplemented commonError = commonError{errors.New("not implemented")}
)
