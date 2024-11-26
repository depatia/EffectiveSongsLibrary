package serverrors

import (
	"errors"
	"fmt"
)

func ArgErr(arg string) error {
	errStr := fmt.Sprintf("field %s is required", arg)
	ErrBadRequest := errors.New(errStr)
	return ErrBadRequest
}
