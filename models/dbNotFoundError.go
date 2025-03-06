package models

import (
	"fmt"
	"reflect"
)

type DbNotFoundError struct {
	DbItem interface{}
}

func (err *DbNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", reflect.TypeOf(err.DbItem).Name())
}
