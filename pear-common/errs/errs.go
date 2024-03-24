package errs

import "fmt"


type ErrCode int

type BError struct {
	Code ErrCode
	Msg string
}

func (e *BError) Error () string{
	return fmt.Sprintf("code: %d, msg :%s", int(e.Code), e.Msg)
}

func NewError(code ErrCode, msg string) *BError {
	return &BError{
		Code: code,
		Msg: msg,
	}
}

