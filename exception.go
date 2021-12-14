package exception

import (
	"fmt"
)

type ExceptionCode uint32
type StatusCode int
type Exception struct {
	ExceptionCode     ExceptionCode
	StatusCode        StatusCode
	ExceptionMessage  string
	InternalException error
	Context           map[string]string
}

var mapping = Mapping{}

type Mapping map[ExceptionCode]struct {
	StatusCode StatusCode
	Message    string
}

func Init(m Mapping) {
	mapping = m
}

func (err Exception) Error() string {
	if err.InternalException != nil {
		return fmt.Sprintf("%s: %s", err.ExceptionMessage, err.InternalException.Error())
	}

	return err.ExceptionMessage
}

func New(exceptionCode ExceptionCode) error {
	return Newf(exceptionCode)
}

func Newf(exceptionCode ExceptionCode, args ...interface{}) error {
	if e, ok := mapping[exceptionCode]; ok {
		return Exception{
			ExceptionCode:    exceptionCode,
			StatusCode:       e.StatusCode,
			ExceptionMessage: fmt.Sprintf(e.Message, args...),
		}
	}

	return Exception{
		ExceptionCode:    exceptionCode,
		StatusCode:       0,
		ExceptionMessage: "Unknown error",
	}
}

func NewWith(exceptionCode ExceptionCode, msg string) error {
	return NewfWith(exceptionCode, msg)
}

func NewfWith(exceptionCode ExceptionCode, msg string, args ...interface{}) error {
	return Exception{
		ExceptionCode:    exceptionCode,
		StatusCode:       mapping[exceptionCode].StatusCode,
		ExceptionMessage: fmt.Sprintf(msg, args...),
	}
}

func Wrap(err error, exceptionCode ExceptionCode) error {
	return Wrapf(err, exceptionCode)
}

func Wrapf(err error, exceptionCode ExceptionCode, args ...interface{}) error {
	if e, ok := mapping[exceptionCode]; ok {
		return Exception{
			ExceptionCode:     exceptionCode,
			StatusCode:        e.StatusCode,
			ExceptionMessage:  fmt.Sprintf(e.Message, args...),
			InternalException: err,
		}
	}

	return Exception{
		ExceptionCode:     exceptionCode,
		StatusCode:        mapping[exceptionCode].StatusCode,
		ExceptionMessage:  "Unknown error",
		InternalException: err,
	}
}

func WrapWith(err error, exceptionCode ExceptionCode, msg string) error {
	return WrapfWith(err, exceptionCode, msg)
}

func WrapfWith(err error, exceptionCode ExceptionCode, msg string, args ...interface{}) error {
	return Exception{
		ExceptionCode:     exceptionCode,
		StatusCode:        mapping[exceptionCode].StatusCode,
		ExceptionMessage:  fmt.Sprintf(msg, args...),
		InternalException: err,
	}
}

func GetExceptionMessage(err error) string {
	if e, ok := err.(Exception); ok {
		return e.ExceptionMessage
	}

	return err.Error()
}

func GetExceptionCode(err error) ExceptionCode {
	if e, ok := err.(Exception); ok {
		return e.ExceptionCode
	}

	return 0xFFFFFFFF
}

func GetStatusCode(err error) StatusCode {
	if e, ok := err.(Exception); ok {
		return e.StatusCode
	}

	return 500
}

func GetContext(err error) map[string]string {
	if e, ok := err.(Exception); ok {
		return e.Context
	}

	return nil
}
