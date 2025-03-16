package models

import (
	"fmt"
	"reflect"
)

type DbItemAlreadyExistsError struct {
	DbItem interface{}
}

func (err *DbItemAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists", reflect.TypeOf(err.DbItem).Name())
}
