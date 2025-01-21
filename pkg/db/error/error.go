package dberror

import (
	"errors"

	"gorm.io/gorm"
)

// ErrNotFound is a not found error
var ErrNotFound = errors.New("not found")

var ErrInvalidArticleCacheStatus = errors.New("invalid article cache status")

func IsNotfoundErr(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrNotFound)
}
