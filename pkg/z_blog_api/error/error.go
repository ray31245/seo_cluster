package zblogerror

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrHTTPStatusCodeError = errors.New("http status code error")
	ErrHTTPInternal        = fmt.Errorf("%w: internal", ErrHTTPStatusCodeError)
	ErrHTTPBadRequest      = fmt.Errorf("%w: bad request", ErrHTTPStatusCodeError)
	ErrHTTPUnauthorized    = fmt.Errorf("%w: unauthorized", ErrHTTPStatusCodeError)
	ErrHTTPForbidden       = fmt.Errorf("%w: forbidden", ErrHTTPStatusCodeError)
	ErrHTTPNotFound        = fmt.Errorf("%w: not found", ErrHTTPStatusCodeError)
	ErrIllegalAccess       = fmt.Errorf("%w: illegal access", ErrHTTPStatusCodeError)
)

const (
	StatusIllegalAccess = 419 // this is not a standard http status code, but it is defined by z-blog
)

func NewHTTPStatusCodeError(statusCode int) error {
	// check if status code is 2xx
	if statusCode/100 == 2 { //nolint:mnd
		return nil
	}

	switch statusCode {
	case http.StatusBadRequest:
		return ErrHTTPBadRequest
	case http.StatusUnauthorized:
		return ErrHTTPUnauthorized
	case http.StatusForbidden:
		return ErrHTTPForbidden
	case http.StatusNotFound:
		return ErrHTTPNotFound
	case http.StatusInternalServerError:
		return ErrHTTPInternal
	case StatusIllegalAccess:
		return ErrIllegalAccess
	default:
		return fmt.Errorf("%w: %d", ErrHTTPStatusCodeError, statusCode)
	}
}

func NewHTTPStatusCodeErrWithMsg(statusCode int, msg string) error {
	codeErr := NewHTTPStatusCodeError(statusCode)
	if codeErr == nil {
		return nil
	}

	return fmt.Errorf("%w: with message: %s", codeErr, msg)
}

func NewHTTPStatusCodeErrorf(statusCode int, format string, args ...interface{}) error {
	if statusCode/100 == 2 { //nolint:mnd
		return nil
	}

	return fmt.Errorf("%w: "+format, append([]interface{}{NewHTTPStatusCodeError(statusCode)}, args...)...) //nolint: err113
}
