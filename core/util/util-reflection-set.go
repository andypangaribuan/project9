/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package util

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unsafe"

	"github.com/andypangaribuan/project9/f9"
)

func (slf *srUtil) ReflectionSet(obj interface{}, bind map[string]interface{}) error {
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Ptr {
		return errors.New("obj must be a pointer")
	}

	objVal = objVal.Elem()
	objType := objVal.Type()

	if objVal.Kind() != reflect.Struct && objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
		objType = objVal.Type()
	}

	if objVal.Kind() == reflect.Struct {
		for i := 0; i < objType.NumField(); i++ {
			rs := objType.Field(i)
			rf := objVal.Field(i)
			fieldName := rs.Name

			if bindValue, ok := bind[fieldName]; ok {
				if !rf.IsValid() {
					return fmt.Errorf("invalid field: %v", fieldName)
				}

				err := slf.reflectionSet(rs, rf, bindValue)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (slf *srUtil) ReflectionGet(obj interface{}, fieldName string) (interface{}, error) {
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Ptr {
		return nil, errors.New("obj must be a pointer")
	}

	objVal = objVal.Elem()
	objType := objVal.Type()

	if objVal.Kind() != reflect.Struct && objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
		objType = objVal.Type()
	}

	if objVal.Kind() == reflect.Struct {
		for i := 0; i < objType.NumField(); i++ {
			rs := objType.Field(i)
			rf := objVal.Field(i)

			if rs.Name == fieldName {
				if rf.IsValid() {
					return rf.Interface(), nil
				}
			}
		}
	}

	return nil, errors.New("not found")
}

func (slf *srUtil) reflectionSet(sf reflect.StructField, rv reflect.Value, obj interface{}) (err error) {
	switch rv.CanSet() {
	case true:
		err = slf.reflectionPublicSet(sf, rv, obj)
	case false:
		err = slf.reflectionPrivateSet(sf, rv, obj)
	}
	return
}

func (slf *srUtil) reflectionPublicSet(rs reflect.StructField, rv reflect.Value, obj interface{}) error {
	err := f9.PanicCatcher(func() {
		rv.Set(reflect.ValueOf(obj))
	})
	return slf.reflectionSetError(rs.Name, err)
}

func (slf *srUtil) reflectionPrivateSet(rs reflect.StructField, rv reflect.Value, obj interface{}) error {
	var first rune
	for _, c := range rs.Name {
		first = c
		break
	}

	if unicode.IsUpper(first) {
		return fmt.Errorf("cannot set the field: %v", rs.Name)
	}

	ptr := unsafe.Pointer(rv.UnsafeAddr())
	newRV := reflect.NewAt(rv.Type(), ptr)
	val := newRV.Elem()
	err := f9.PanicCatcher(func() {
		val.Set(reflect.ValueOf(obj))
	})
	return slf.reflectionSetError(rs.Name, err)
}

func (*srUtil) reflectionSetError(fieldName string, err error) error {
	if err != nil {
		return fmt.Errorf("%v\nfield name: %v", err, fieldName)
	}
	return err
}
