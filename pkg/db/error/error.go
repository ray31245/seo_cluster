package dberror

import (
	"errors"

	"gorm.io/gorm"
)

func IsNotfoundErr(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
