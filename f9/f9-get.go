/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
)

 func GetStructPbValue[R any](val *structpb.Value) R {
	defaultValue := *new(R)
	resType := fmt.Sprintf("%T", defaultValue)

	switch resType {
	case "string", "*string":
		if val != nil && val.AsInterface() != nil {
			return toR[string, R](resType, val.GetStringValue())
		}

	case "int", "*int":
		if val != nil && val.AsInterface() != nil {
			return toR[int, R](resType, int(val.GetNumberValue()))
		}

	case "int64", "*int64":
		if val != nil && val.AsInterface() != nil {
			return toR[int64, R](resType, int64(val.GetNumberValue()))
		}

	case "float32", "*float32":
		if val != nil && val.AsInterface() != nil {
			return toR[float32, R](resType, float32(val.GetNumberValue()))
		}

	case "float64", "*float64":
		if val != nil && val.AsInterface() != nil {
			return toR[float64, R](resType, val.GetNumberValue())
		}
	}

	return defaultValue
}

func toR[T any, R any](typeData string, val T) R {
	var v interface{}
	if typeData[0:1] == "*" {
		v = &val
	} else {
		v = val
	}
	return v.(R)
}


func GetKeyFromMap[K comparable, V any](m map[K]V) []K {
	arr := make([]K, 0)

	for k := range m {
		arr = append(arr, k)
	}

	return arr
}

func GetMapStructPbValue(m map[string]interface{}) map[string]*structpb.Value {
	dictio := make(map[string]*structpb.Value, 0)

	for k, v := range m {
		s, _ := structpb.NewValue(v)
		dictio[k] = s
	}

	return dictio
}
