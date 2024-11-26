package validator

import (
	"SongsLibrary/internal/services/serverrors"
	"errors"
	"reflect"
)

var ErrNotStruct = errors.New("enter the struct")

func ValidateStruct(s any) error {
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() {
			return serverrors.ArgErr(v.Type().Field(i).Name)
		}
	}

	return nil
}
