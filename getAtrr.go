package main

import (
	"reflect"
)

func getAttr(obj any, fieldName string) *reflect.Value {
	pointToStruct := reflect.ValueOf(obj)

	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}

	curField := curStruct.FieldByName(fieldName)
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return &curField
}
