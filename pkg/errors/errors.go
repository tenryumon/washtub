package errors

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
)

type Kind int

const (
	KindBadRequest Kind = iota + 1
	KindUnauthorized
	KindInternalError
	KindInvalidForm

	MsgErrInvalidForm = "Maaf, terdapat kesalahan dalam form. Harap melakukan pengecekan ulang."
)

func (k Kind) HTTPCode() int {
	switch k {
	case KindBadRequest:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindInternalError:
		return http.StatusInternalServerError
	case KindInvalidForm:
		return 400001
	default:
		return http.StatusInternalServerError
	}
}

type Errors struct {
	err    error
	kind   Kind
	fields Fields
	forms  Forms
}

func (e *Errors) Error() string {
	return e.err.Error()
}

func (e *Errors) Kind() Kind {
	return e.kind
}

func (e *Errors) Fields() Fields {
	return e.fields
}

func (e *Errors) Forms() Forms {
	return e.forms
}

// New creates a completely new error.
func New(s string, v ...any) *Errors {
	e := errors.New(s)
	return Wrap(e, v...)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

// InternalHTTPCode implements internal http code to accomodate custom error code number
// from internal implementation.
type InternalHTTPCode int

func (i InternalHTTPCode) String() string {
	switch i {
	case http.StatusOK:
		return "200000"
	case http.StatusUnauthorized:
		return "403000"
	case http.StatusBadRequest:
		return "400000"
	case 400001:
		return "400001"
	default:
		return "500000"
	}
}

// HTTPCode returns the http code based on the errors.
func HTTPCode(err error) InternalHTTPCode {
	if err == nil {
		return InternalHTTPCode(http.StatusOK)
	}
	var errs *Errors
	if As(err, &errs) {
		return InternalHTTPCode(errs.Kind().HTTPCode())
	}
	return InternalHTTPCode(http.StatusInternalServerError)
}

// Wrap the error.
func Wrap(err error, v ...any) *Errors {
	// Don't return a wrapped error if the error is nil, we should return nil to ensure the usual
	// error checking flow working as intended.
	if err == nil {
		return nil
	}

	e := &Errors{
		err: err,
	}

	var errs *Errors
	if As(err, &errs) {
		e.forms = errs.Forms()
		e.kind = errs.Kind()
	}

	for idx := range v {
		switch t := v[idx].(type) {
		case error:
			e.err = errors.Join(e.err, t)
		case Kind:
			e.kind = t
		case Fields:
			e.fields = t
		case Forms:
			e.forms = t
		}
	}
	return e
}

// Join wraps multiple error into one error. This functionality is added in Go 1.20.
func Join(err error, v ...error) error {
	if len(v) == 0 {
		return err
	}

	// If we have the internal error type here, we should wrap the internal error with the incoming errors.
	var errs *Errors
	if errors.As(err, &errs) {
		v = append([]error{errs.err}, v...)
		errs.err = errors.Join(v...)
		return errs
	}
	v = append([]error{err}, v...)
	return errors.Join(v...)
}

// Is shadows the errors.Is function to check whether the internal error type is the same
// with the one we want to compare.
func (e *Errors) Is(err error) bool {
	return errors.Is(e.err, err)
}

// As shadows the errors.As function to check whether the internal error type is the same
// with the one we want to compare.
func (e *Errors) As(target any) bool {
	return errors.As(e.err, target)
}

type Fields []any
type Forms []string

// NewFields is for safely creating error fields becauase the error fields format
// is a key value to add more context to the error.
func NewFields(kv ...any) (f Fields) {
	if kv == nil {
		return nil
	}

	kvlen := len(kv)
	if kvlen%2 == 0 {
		f = kv
		return
	}

	// Ensure that the Fields is never 'odd'. If they 'key' is not available then we
	// should replace the 'key' with 'unknown?'.
	newKV := make([]interface{}, kvlen+1)
	for i := 0; i < kvlen; i++ {
		if i == kvlen-1 {
			newKV[i] = "unknown?"
			// We will always know this is safe to do because we previously
			// set the array capacity to kv length + 1.
			newKV[i+1] = kv[i]
			break
		}
		newKV[i] = kv[i]
	}

	f = newKV
	return
}

// NewForms is for safely creating error forms becauase the error forms format is key-value.
func NewForms(kv ...string) (f Forms) {
	if kv == nil {
		return nil
	}

	kvlen := len(kv)
	if kvlen%2 == 0 {
		f = kv
		return
	}

	// Ensure that the Forms is never 'odd'. If they 'key' is not available then we
	// should replace the 'key' with 'unknown?'.
	newKV := make([]string, kvlen+1)
	for i := 0; i < kvlen; i++ {
		if i == kvlen-1 {
			newKV[i] = "unknown?"
			// We will always know this is safe to do because we previously
			// set the array capacity to kv length + 1.
			newKV[i+1] = kv[i]
			break
		}
		newKV[i] = kv[i]
	}

	f = newKV
	return
}

func (f Forms) ToMapStringArray() map[string][]string {
	kvlen := len(f)
	// We should not convert the forms because it doesn't fulfill the
	// fields k/v criteria. Possibly, the fields is created without
	// using the NewFields function.
	if kvlen%2 != 0 {
		f = NewForms(f...)
		kvlen = len(f)
	}
	result := make(map[string][]string)
	for i := 0; i < kvlen; i += 2 {
		if _, ok := result[f[i]]; ok {
			result[f[i]] = append(result[f[i]], f[i+1])
		} else {
			result[f[i]] = []string{f[i+1]}
		}
	}
	return result
}

// ToSlogAttributes safely converts fields([]any) to slog.Attr. The function will
// return an empty []slog.Attr if the fields is empty.
//
// Please note that []Fields are expected to form a []Field{"string", "value"} where
// the first key is always a 'string' and 'any' for the value. If the key is not a
// 'string', the function will convert the key to 'string' to be compatible with the
// slog.Attr standard.
func (f Fields) ToSlogAttributes() []slog.Attr {
	kvlen := len(f)
	// We should not convert the fields because it doesn't fulfill the
	// fields k/v criteria. Possibly, the fields is created without
	// using the NewFields function.
	if kvlen%2 != 0 {
		f = NewFields(f...)
		kvlen = len(f)
	}

	var attrs = make([]slog.Attr, kvlen/2)
	for i := 0; i < kvlen; i += 2 {
		var slogKey string
		key, isString := f[i].(string)
		// If the key is not string then we will convert the key using fmt.Sprinf to a string.
		// This will allocates but only if the key is not string, so we think the trade-off
		// is worth-it.
		if !isString {
			// We use f[i] here because the value of 'key' will be empty if the field is not
			// a string.
			slogKey = fmt.Sprintf("%v", f[i])
		} else {
			slogKey = key
		}

		// The counter for attribute array need to be divided by 2 because we always increase
		// i by 2.
		attrs[i/2] = slog.Attr{
			Key:   slogKey,
			Value: slog.AnyValue(f[i+1]),
		}
		if i+2 == kvlen {
			break
		}
	}
	return attrs
}
