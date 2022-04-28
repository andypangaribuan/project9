/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package abs

type Json interface {
	Marshal(obj interface{}) ([]byte, error)
	UnMarshal(data []byte, out interface{}) error
	Encode(obj interface{}) (string, error)
	Decode(jsonStr string, out interface{}) error
	MapToJson(maps map[string]interface{}) (string, error)
}
