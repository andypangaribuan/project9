/*
 * Copyright (c) 2023.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Val[T any](obj interface{}, defaultValue ...T) *T {
	val, _ := SVal(obj, defaultValue...)
	return val
}

func SVal[T any](obj interface{}, defaultValue ...T) (*T, bool) {
	var objVal interface{}

	if obj != nil {
		switch v := obj.(type) {
		case *wrapperspb.BoolValue:
			if v != nil {
				objVal = v.Value
			}

		case *wrapperspb.StringValue:
			if v != nil {
				objVal = v.Value
			}

		case *wrapperspb.Int32Value:
			if v != nil {
				objVal = v.Value
				if val, ok := objVal.(T); !ok {
					switch fmt.Sprintf("%T", val) {
					case "int":
						objVal = int(v.Value)
					case "int64":
						objVal = int64(v.Value)
					}
				}
			}

		case *wrapperspb.Int64Value:
			if v != nil {
				objVal = v.Value
				if val, ok := objVal.(T); !ok {
					switch fmt.Sprintf("%T", val) {
					case "int":
						objVal = int(v.Value)
					case "int32":
						objVal = int32(v.Value)
					}
				}
			}

		case *wrapperspb.FloatValue:
			if v != nil {
				objVal = v.Value
				if val, ok := objVal.(T); !ok {
					switch fmt.Sprintf("%T", val) {
					case "float64":
						objVal = float64(v.Value)
					}
				}
			}

		case *wrapperspb.DoubleValue:
			if v != nil {
				objVal = v.Value
				if val, ok := objVal.(T); !ok {
					switch fmt.Sprintf("%T", val) {
					case "float32":
						objVal = float32(v.Value)
					}
				}
			}
		}
	}

	switch {
	case objVal == nil && len(defaultValue) > 0:
		val := defaultValue[0]
		return &val, true

	case objVal == nil:
		return nil, true
	}

	val, ok := objVal.(T)
	return &val, ok
}
