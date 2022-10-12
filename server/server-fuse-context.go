/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
	"github.com/andypangaribuan/project9/server/proto/gen/grf"
	"github.com/gorilla/schema"
)

var (
	mpfMutex           sync.Mutex
	mpfDecoderInstance *schema.Decoder
)

// region parser
func (slf *srFuseContext) getHeader() map[string]string {
	if len(slf.header) == 0 {
		header := make(map[string]string, 0)

		switch {
		case slf.fiberCtx != nil:
			for k, v := range slf.fiberCtx.GetReqHeaders() {
				header[strings.ToLower(k)] = v
			}

		case slf.grpcCtx != nil:
			if len(slf.grpcCtx.header) > 0 {
				for k, v := range slf.grpcCtx.header {
					header[strings.ToLower(k)] = v
				}
			}

		default:
			panic("unimplemented")
		}

		slf.header = header
	}

	return slf.header
}

func (slf *srFuseContext) mpfDecoder() *schema.Decoder {
	if mpfDecoderInstance != nil {
		return mpfDecoderInstance
	}

	mpfMutex.Lock()
	defer mpfMutex.Unlock()

	if mpfDecoderInstance != nil {
		return mpfDecoderInstance
	}

	decoder := schema.NewDecoder()
	decoder.SetAliasTag("json")
	mpfDecoderInstance = decoder

	return mpfDecoderInstance
}

func (slf *srFuseContext) Params(key string) string {
	switch {
	case slf.fiberCtx != nil:
		return slf.fiberCtx.Params(key, "")

	case slf.grpcCtx != nil:
		val := ""
		if len(slf.grpcCtx.params) > 0 {
			if v, ok := slf.grpcCtx.params[key]; ok {
				val = v
			}
		}
		return val

	default:
		panic("unimplemented")
	}
}

func (slf *srFuseContext) Parser(header, body interface{}) (bool, error) {
	mapHeader := slf.getHeader()

	if header != nil {
		data, err := p9.Json.Marshal(mapHeader)
		if err == nil {
			err = p9.Json.UnMarshal(data, header)
		}

		if err != nil {
			err = p9.Err.WithStack(err, 1)
			return false, slf.r500InternalServerError(err)
		}
	}

	if body != nil && slf.fiberCtx != nil {
		cType, ok := mapHeader["content-type"]
		if !ok {
			err := errors.New("unknown content-type")
			return false, slf.r500InternalServerError(err)
		}

		if idx := strings.Index(cType, ";"); idx > -1 {
			cType = cType[0:idx]
		}

		switch cType {
		case "application/json":
			err := slf.fiberCtx.BodyParser(&body)
			if err != nil {
				err = p9.Err.WithStack(err, 1)
				return false, slf.r500InternalServerError(err)
			}

		case "application/x-www-form-urlencoded":
			var (
				err  error
				res  interface{}
				data []byte
			)

			res, err = reTagAny(body, func(structureType reflect.Type, fieldIndex int) reflect.StructTag {
				f := structureType.Field(fieldIndex)
				tag := f.Tag
				jsonTag := tag.Get("json")
				formTag := tag.Get("form")

				if jsonTag != "" && formTag == "" {
					newTag := fmt.Sprintf(`%v form:"%v"`, tag, jsonTag)
					st := reflect.StructTag(newTag)
					return reflect.StructTag(st)
				}

				return ""
			})

			if err == nil {
				err = slf.fiberCtx.BodyParser(res)
			}
			if err == nil {
				data, err = p9.Json.Marshal(res)
			}
			if err == nil {
				err = p9.Json.UnMarshal(data, body)
			}

			if err != nil {
				err = p9.Err.WithStack(err, 1)
				return false, slf.r500InternalServerError(err)
			}

		case "multipart/form-data":
			mf, err := slf.fiberCtx.MultipartForm()
			if err == nil {
				err = slf.mpfDecoder().Decode(body, mf.Value)
			}

			if err != nil {
				err = p9.Err.WithStack(err, 1)
				return false, slf.r500InternalServerError(err)
			}

			slf.multipartFile = mf.File
		}
	}

	if body != nil && slf.grpcCtx != nil {
		if len(slf.grpcCtx.payload) > 0 {
			payload := make(map[string]interface{}, 0)
			for k, v := range slf.grpcCtx.payload {
				if v != nil {
					payload[k] = v.AsInterface()
				}
			}

			data, err := p9.Json.Marshal(payload)
			if err == nil {
				err = p9.Json.UnMarshal(data, body)
			}

			if err != nil {
				err = p9.Err.WithStack(err)
				return false, slf.r500InternalServerError(err)
			}
		}
	}

	return true, nil
}

//endregion

//region auth

func (slf *srFuseContext) AuthX() interface{} {
	return slf.authX
}

func (slf *srFuseContext) AuthY() interface{} {
	return slf.authY
}

func (slf *srFuseContext) AuthZ() interface{} {
	return slf.authZ
}

func (slf *srFuseContext) SetAuth(authX, authY, authZ interface{}) {
	slf.authX = authX
	slf.authY = authY
	slf.authZ = authZ
}

//endregion

// region util
func (slf *srFuseContext) wrapError(err error) error {
	if err != nil {
		return fmt.Errorf("endpoint-path: %v\n%w", slf.path, err)
	}
	return fmt.Errorf("endpoint-path: %v", slf.path)
}

//endregion

//region response

func (slf *srFuseContext) RString(code int, data string) error {
	return slf.sendRawA(code, data)
}

func (slf *srFuseContext) RJson(code int, data interface{}) error {
	return slf.sendRawB(code, data)
}

func (slf *srFuseContext) RJsonRaw(code int, data []byte) error {
	return slf.sendRawB(code, f9.ToJsonRaw(data))
}

func (slf *srFuseContext) R200OK(data interface{}, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:   http.StatusOK,
		Status: fuseDefaultStatus.R200OK,
		Data:   data,
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R400BadRequest(message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusBadRequest,
		Status:  fuseDefaultStatus.R400BadRequest,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R400BadRequest),
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R401Unauthorized(message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusUnauthorized,
		Status:  fuseDefaultStatus.R401Unauthorized,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R401Unauthorized),
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R403Forbidden(message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusForbidden,
		Status:  fuseDefaultStatus.R403Forbidden,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R403Forbidden),
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R404NotFound(message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusNotFound,
		Status:  fuseDefaultStatus.R404NotFound,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R404NotFound),
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R406NotAcceptable(message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusNotAcceptable,
		Status:  fuseDefaultStatus.R406NotAcceptable,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R406NotAcceptable),
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R428PreconditionRequired(message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusPreconditionRequired,
		Status:  fuseDefaultStatus.R428PreconditionRequired,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R428PreconditionRequired),
	}

	return slf.send(fo, opt...)
}

func (slf *srFuseContext) R500InternalServerError(err error, opt ...FuseOpt) error {
	return slf.r500InternalServerError(p9.Err.WithStack(slf.wrapError(err), 1), opt...)
}

func (slf *srFuseContext) r500InternalServerError(err error, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusInternalServerError,
		Status:  fuseDefaultStatus.R500InternalServerError,
		Message: fuseDefaultMessage.R500InternalServerError,
		Error:   err,
	}

	return slf.send(fo, opt...)
}

//endregion

//region send response

func (slf *srFuseContext) sendRawA(code int, data string) error {
	switch {
	case slf.fiberCtx != nil:
		return slf.fiberCtx.Status(code).SendString(data)

	case slf.grpcCtx != nil:
		ls, err := p9.Json.Marshal(data)
		if err != nil {
			return p9.Err.WithStack(err)
		}

		slf.grpcCtx.response = &grf.Response{Data: ls}
		return nil
	}

	panic("unimplemented")
}

func (slf *srFuseContext) sendRawB(code int, data interface{}) error {
	switch {
	case slf.fiberCtx != nil:
		return slf.fiberCtx.Status(code).JSON(data)

	case slf.grpcCtx != nil:
		switch cast := data.(type) {
		case json.RawMessage:
			slf.grpcCtx.response = &grf.Response{Data: cast}
			return nil

		default:
			ls, err := p9.Json.Marshal(data)
			if err != nil {
				return p9.Err.WithStack(err)
			}

			slf.grpcCtx.response = &grf.Response{Data: ls}
			return nil
		}
	}

	panic("unimplemented")
}

func (slf *srFuseContext) send(fo FuseOpt, opt ...FuseOpt) error {
	type srMeta struct {
		Code    int         `json:"code"`
		Status  string      `json:"status,omitempty"`
		Message string      `json:"message,omitempty"`
		Address string      `json:"address,omitempty"`
		Error   string      `json:"error,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	type srResponse struct {
		Meta srMeta      `json:"meta"`
		Data interface{} `json:"data,omitempty"`
	}

	if len(opt) > 0 {
		o := opt[0]
		fo.Status = f9.Ternary(o.Status != "", o.Status, fo.Status)
		fo.Message = f9.Ternary(o.Message != "", o.Message, fo.Message)
		fo.Address = f9.Ternary(o.Address != "", o.Address, fo.Address)
		fo.Error = f9.Ternary(o.Error != nil, o.Error, fo.Error)
		fo.MetaData = f9.Ternary(o.MetaData != nil, o.MetaData, fo.MetaData)
		fo.Data = f9.Ternary(o.Data != nil, o.Data, fo.Data)
	}

	response := srResponse{
		Meta: srMeta{
			Code:    fo.code,
			Status:  fo.Status,
			Message: fo.Message,
			Address: fo.Address,
			Data:    fo.MetaData,
		},
		Data: fo.Data,
	}

	if fuseWithError != nil && *fuseWithError && fo.Error != nil {
		timeNow := p9.Conv.Time.ToStrFull(f9.TimeNow())
		response.Meta.Error = fmt.Sprintf("[ERROR] %v\n---------------------------\n%+v", timeNow, fo.Error)

		fmt.Printf("\n\n%v\n", response.Meta.Error)
	}

	grpcResponse := func() error {
		var (
			err      error
			data     []byte
			metaData []byte
		)

		if response.Data != nil {
			data, err = p9.Json.Marshal(response.Data)
			if err != nil {
				return p9.Err.WithStack(err)
			}
		}

		if response.Meta.Data != nil {
			metaData, err = p9.Json.Marshal(response.Meta.Data)
			if err != nil {
				return p9.Err.WithStack(err)
			}
		}

		slf.grpcCtx.response = &grf.Response{
			Data: data,
			Meta: &grf.RMeta{
				Code:    int32(response.Meta.Code),
				Status:  response.Meta.Status,
				Message: response.Meta.Status,
				Address: response.Meta.Address,
				Error:   response.Meta.Error,
				Data:    metaData,
			},
		}

		return nil
	}

	switch {
	case slf.fiberCtx != nil:
		return slf.fiberCtx.Status(fo.code).JSON(response)
	case slf.grpcCtx != nil:
		return grpcResponse()
	}

	panic("unimplemented")
}

func (slf *srFuseContext) grpcSend(err error) (*grf.Response, error) {
	return slf.grpcCtx.response, err
}

//endregion
