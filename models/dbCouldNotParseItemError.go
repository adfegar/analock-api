package models

import (
	"fmt"
	"reflect"
)

type DbCouldNotParseItemError struct {
	DbItem interface{}
}

func (err *DbCouldNotParseItemError) Error() string {
	return fmt.Sprintf("Could not parse %s", reflect.TypeOf(&err.DbItem).Name())
}
