package error

import (
	"fmt"
)

var (
	ErrNoCategoryNeedToBePublished = fmt.Errorf("no category need to be published")
)
