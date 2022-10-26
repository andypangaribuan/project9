/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package f9

import (
	"strings"

	"github.com/andypangaribuan/project9/constraint"
	"google.golang.org/protobuf/types/known/structpb"
)

func IfAllNil(items ...interface{}) bool {
	allNil := true

	for _, item := range items {
		if item != nil {
			allNil = false
			break
		}
	}

	return allNil
}

func IfStrNotNilButEmpty(val *string) bool {
	if val == nil {
		return false
	}

	if strings.TrimSpace(*val) == "" {
		return true
	}

	return false
}

func IfStrNilOrEmpty(val *string) bool {
	if val == nil {
		return true
	}

	if strings.TrimSpace(*val) == "" {
		return true
	}

	return false
}

func IfPtrValueEqual[T constraint.ComparisonType, PtrT *T](left, right PtrT) bool {
	if left == nil || right == nil {
		return false
	}

	return *left == *right
}

func IfHaveEmpty(items ...string) bool {
	haveEmpty := false

	for _, item := range items {
		if strings.TrimSpace(item) == "" {
			haveEmpty = true
			break
		}
	}

	return haveEmpty
}

func IfAllStructPbNil(items ...*structpb.Value) bool {
	allNil := true

	for _, item := range items {
		if item != nil && item.AsInterface() != nil {
			allNil = false
			break
		}
	}

	return allNil
}
