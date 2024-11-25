package wordpresserror

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
	default:
		return fmt.Errorf("%w: %d", ErrHTTPStatusCodeError, statusCode)
	}
}
