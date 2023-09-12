/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

type result struct {
	t        reflect.Type
	changed  bool
	hasIFace bool
}

type reTagProcessor func(structureType reflect.Type, fieldIndex int) reflect.StructTag

func reTagAny(p interface{}, processor reTagProcessor) (interface{}, error) {
	return convert(p, processor, true)
}

func convert(p interface{}, processor reTagProcessor, any bool) (interface{}, error) {
	strPtrVal := reflect.ValueOf(p)
	res, err := getType(strPtrVal.Type().Elem(), processor, any)
	if err != nil {
		return nil, err
	}

	newPtrVal := reflect.NewAt(res.t, unsafe.Pointer(strPtrVal.Pointer()))
	return newPtrVal.Interface(), nil
}

func getType(structType reflect.Type, processor reTagProcessor, any bool) (*result, error) {
	return makeType(structType, processor, any)
}

func makeType(t reflect.Type, processor reTagProcessor, any bool) (*result, error) {
	switch t.Kind() {
	case reflect.Struct:
		return makeStructType(t, processor, any)

	case reflect.Ptr:
		res, err := getType(t.Elem(), processor, any)
		if err != nil {
			return nil, err
		}

		if !res.changed {
			return &result{t: t, changed: false}, nil
		}

		return &result{t: reflect.PtrTo(res.t), changed: true}, nil

	case reflect.Array:
		res, err := getType(t.Elem(), processor, any)
		if err != nil {
			return nil, err
		}

		if !res.changed {
			return &result{t: t, changed: false}, nil
		}

		return &result{t: reflect.ArrayOf(t.Len(), res.t), changed: true}, nil

	case reflect.Slice:
		res, err := getType(t.Elem(), processor, any)
		if err != nil {
			return nil, err
		}

		if !res.changed {
			return &result{t: t, changed: false}, nil
		}

		return &result{t: reflect.SliceOf(res.t), changed: true}, nil

	case reflect.Map:
		resKey, err := getType(t.Key(), processor, any)
		if err != nil {
			return nil, err
		}

		resElem, err := getType(t.Elem(), processor, any)
		if err != nil {
			return nil, err
		}

		if !resKey.changed && !resElem.changed {
			return &result{t: t, changed: false}, nil
		}

		return &result{t: reflect.MapOf(resKey.t, resElem.t), changed: true}, nil

	case reflect.Interface:
		if any {
			return &result{t: t, changed: false, hasIFace: true}, nil
		}
		fallthrough

	case
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer:
		return nil, errors.New("tags.Map: Unsupported type: " + t.Kind().String())

	default:
		// don't modify type in another case
		return &result{t: t, changed: false}, nil
	}
}

func makeStructType(structType reflect.Type, processor reTagProcessor, any bool) (*result, error) {
	if structType.NumField() == 0 {
		return &result{t: structType, changed: false}, nil
	}

	changed := false
	hasPrivate := false
	hasIFace := false
	fields := make([]reflect.StructField, 0, structType.NumField())

	for i := 0; i < structType.NumField(); i++ {
		strField := structType.Field(i)

		if isExported(strField.Name) {
			oldType := strField.Type
			new, err := getType(oldType, processor, any)
			if err != nil {
				return nil, err
			}

			strField.Type = new.t
			if oldType != new.t {
				changed = true
			}

			if new.hasIFace {
				hasIFace = true
			}

			oldTag := strField.Tag
			newTag := processor(structType, i)
			strField.Tag = newTag

			if oldTag != newTag {
				changed = true
			}
		} else {
			hasPrivate = true
			if !structTypeConstructorBugWasFixed {
				strField.PkgPath = ""
				strField.Name = ""
			}
		}

		fields = append(fields, strField)
	}

	if !changed {
		return &result{t: structType, changed: false, hasIFace: hasIFace}, nil
	} else if hasPrivate {
		return nil, fmt.Errorf("unable to change tags for type %s, because it contains unexported fields", structType)
	}

	newType := reflect.StructOf(fields)
	err := compareStructTypes(structType, newType)
	if err != nil {
		return nil, err
	}

	return &result{t: newType, changed: true, hasIFace: hasIFace}, nil
}

func isExported(name string) bool {
	b := name[0]
	return !('a' <= b && b <= 'z') && b != '_'
}

func compareStructTypes(source, result reflect.Type) error {
	if source.Size() != result.Size() {
		return errors.New("tags.Map: Unexpected case - type has a size different from size of original type")
	}
	return nil
}

var structTypeConstructorBugWasFixed bool

func init() {
	switch {
	case strings.HasPrefix(runtime.Version(), "go1.7"):
		// there is bug in reflect.StructOf
	default:
		structTypeConstructorBugWasFixed = true
	}
}
