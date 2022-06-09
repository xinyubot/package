package errors

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

var (
	_ error        = (*Error)(nil) // make sure Error implements error interface
	_ fmt.Stringer = (*Error)(nil) // make sure Error implements fmt.Stringer
)

// ErrorIface defines a set of methods that a gin-awesome-style error should implement.
type ErrorIface interface {
	GetCode() int
	GetReason() string
	GetMessage() string
	GetData() any
	GetStack() string
	GetFileLine() string
	GetExtraDataMap() map[string]any
	Is(error) bool
	Equal(error) bool
	error
}

// Option ...
type Option func(*Error)

// Error is a trivial implementation of ErrorIface, hence of error.
type Error struct {
	code         int
	reason       string
	message      string
	stack        string
	fileLine     string
	data         any
	extraDataMap map[string]any
}

// GetCode ...
func (e *Error) GetCode() int {
	if e != nil {
		return e.code
	}
	return 0
}

// GetReason ...
func (e *Error) GetReason() string {
	if e != nil {
		return e.reason
	}
	return ""
}

// GetMessage ...
func (e *Error) GetMessage() string {
	if e != nil {
		return e.message
	}
	return ""
}

// GetStack ...
func (e *Error) GetStack() string {
	if e != nil {
		return e.stack
	}
	return ""
}

// GetFileLine ...
func (e *Error) GetFileLine() string {
	if e != nil {
		return e.fileLine
	}
	return ""
}

// GetData ...
func (e *Error) GetData() any {
	if e != nil {
		return e.data
	}
	return nil
}

// GetExtraDataMap ...
func (e *Error) GetExtraDataMap() map[string]any {
	if e != nil {
		return e.extraDataMap
	}
	return nil
}

// Error returns a customized format of the entire error; it implements standard `error` interface.
func (e *Error) Error() string {
	return fmt.Sprintf(`code: %d, reason: %s, message: %s, stack: %s, file: %s, data: %#v, extraDataMap: %#v`,
		e.code, e.reason, e.message, e.stack, e.fileLine, e.data, e.extraDataMap,
	)
}

// String returns a customized format of the entire error; it implements `fmt.Stringer` as well.
func (e *Error) String() string {
	return fmt.Sprintf(`code: %d, reason: %s, message: %s, stack: %s, file: %s, data: %#v, extraDataMap: %#v`,
		e.code, e.reason, e.message, e.stack, e.fileLine, e.data, e.extraDataMap,
	)
}

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.reason == e.reason
	}
	return false
}

// DeepEqual compares code and message.
func (e *Error) Equal(err error) bool {
	fe := FromError(err)
	if e.GetCode() == fe.GetCode() && e.GetMessage() == fe.GetMessage() {
		return true
	}
	return false
}

// DeepEqual compares code, reason, message, and data.
func (e *Error) DeepEqual(err error) bool {
	fe := FromError(err)
	if e.GetCode() == fe.GetCode() && e.GetReason() == fe.GetReason() &&
		e.GetMessage() == fe.GetMessage() && e.GetData() == fe.GetData() && reflect.DeepEqual(e.GetExtraDataMap(), fe.GetExtraDataMap()) {
		return true
	}
	return false
}

// New returns a new error with its underlying type being Error. If data is nil, error.data will be assigned with, by default, a pointer to an empty struct{} slice.
func New(code int, reason, message string, data any, options ...Option) ErrorIface {
	if data == nil {
		data = struct{}{}
	}
	_, file, line, _ := runtime.Caller(1)
	fileLine := fmt.Sprintf("%s:%d", file, line)

	Error := Error{
		code:     code,
		reason:   reason,
		message:  message,
		fileLine: fileLine,
		data:     data,
	}
	for _, opt := range options {
		opt(&Error)
	}

	return &Error
}

// NewWithSkip returns a new error with its underlying type being Error. If data is nil, error.data will be assigned with, by default, a pointer to an empty struct{} slice.
// The parameter skip is the number of stack layers to be skipped when logging the file and line info, i.e. the location where the standard error actually occurred.
// Use 1 as the base skip number and increase it by 1 for each layer of encapsulation.
func NewWithSkip(skip, code int, reason, message string, data any, options ...Option) ErrorIface {
	if data == nil {
		data = struct{}{}
	}
	_, file, line, _ := runtime.Caller(skip)
	fileLine := fmt.Sprintf("%s:%d", file, line)

	Error := Error{
		code:     code,
		reason:   reason,
		message:  message,
		fileLine: fileLine,
		data:     data,
	}
	for _, opt := range options {
		opt(&Error)
	}

	return &Error
}

// FromError tries to convert an error to ErrorIface.
// It supports wrapped errors.
func FromError(err error) ErrorIface {
	if err == nil {
		return nil
	}

	if se := new(Error); errors.As(err, &se) {
		return se
	}

	return New(500, "", err.Error(), nil)
}

// FromErrorPro tries to convert an error to ErrorIface with passed in parameters.
// It supports wrapped errors.
func FromErrorPro(err error, code int, reason string, data any, options ...Option) ErrorIface {
	if err == nil {
		return nil
	}

	if se := new(Error); errors.As(err, &se) {
		return se
	}

	return New(code, reason, err.Error(), data, options...)
}

// RecastError tries to add addition information to an ErrorIface.
// It DOES NOT support standard errors by design.
func RecastError(err error, options ...Option) ErrorIface {
	if err == nil {
		return nil
	}

	if se := new(Error); errors.As(err, &se) {
		for _, opt := range options {
			opt(se)
		}
		return se
	}

	panic("RecastError DOES NOT support std errors")
}

// WithStackTrace loads the current stack trace information into the error of given length, max capped at 65536 [2 << 15] (which in most cases should be enough).
// Notice that this function will intensely negatively affect the performance.
func WithStackTrace(sizeInByte int) Option {
	return func(err *Error) {
		var buf [2 << 15]byte
		if sizeInByte > len(buf) {
			sizeInByte = len(buf)
		}
		length := runtime.Stack(buf[:], false)
		if length < sizeInByte {
			err.stack = string(buf[:length])
		} else {
			err.stack = string(buf[:sizeInByte])
		}
	}
}

// WithExtraData loads the passed-in key-val pair into the error.
func WithExtraData(key string, val any) Option {
	return func(err *Error) {
		if err.extraDataMap == nil {
			err.extraDataMap = map[string]any{key: val}
			return
		}
		err.extraDataMap[key] = val
	}
}

// WithExtraDataMap loads the passed-in map into the error.
func WithExtraDataMap(m map[string]any) Option {
	return func(err *Error) {
		if err.extraDataMap == nil {
			err.extraDataMap = m
			return
		}
		for k, v := range m {
			err.extraDataMap[k] = v
		}

	}
}

// WithCode loads the passed-in function from the error.
func WithCode(code int) Option {
	return func(err *Error) {
		err.code = code
	}
}

// WithReason loads the passed-in reason into the error.
func WithReason(reason string) Option {
	return func(err *Error) {
		err.reason = reason
	}
}

// WithMessage loads the passed-in message into the error.
func WithMessage(message string) Option {
	return func(err *Error) {
		err.message = message
	}
}

// WithData loads the passed-in data into the error.
func WithData(data any) Option {
	return func(err *Error) {
		err.data = data
	}
}
