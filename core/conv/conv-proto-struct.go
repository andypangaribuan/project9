/*
 * Copyright (c) 2022, Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package conv

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func (*srProto) AnyToProtoStruct(sr any) (*structpb.Struct, error) {
	data, err := json.Marshal(sr)
	if err != nil {
		return nil, err
	}

	srv := &structpb.Struct{}
	err = protojson.Unmarshal(data, srv)
	return srv, err
}
