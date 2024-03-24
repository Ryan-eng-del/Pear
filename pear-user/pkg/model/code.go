package model

import (
	"cyan.com/pear-common/errs"
)


var (
	NoLegalMobile *errs.BError = errs.NewError(2001, "invalid mobile")
)