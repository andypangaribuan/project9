/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

/* spell-checker: disable */
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/andypangaribuan/project9/clog"
	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/model"
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

func (slf *srFuseContext) URI() string {
	return string(slf.fiberCtx.Request().URI().RequestURI())
}

func (slf *srFuseContext) SetSendResponse(send bool) {
	slf.sendResponse = send
}

func (slf *srFuseContext) GetUnSendResponse() *FuseResponse {
	return slf.unSendResponse
}

func (slf *srFuseContext) GetUnSendResponseOpt() []FuseOpt {
	return slf.unSendResponseOpt
}

func (slf *srFuseContext) Request() FuseContextRequest {
	return slf.reqCtx
}

func (slf *srFuseContextRequest) Body() []byte {
	return slf.fuseCtx.fiberCtx.Body()
}

func (slf *srFuseContextRequest) Header() map[string]string {
	return slf.fuseCtx.fiberCtx.GetReqHeaders()
}

func (slf *srFuseContextRequest) SetHeader(key, value string) {
	if len(slf.fuseCtx.header) > 0 {
		slf.fuseCtx.header[strings.ToLower(key)] = value
	}
}

func (slf *srFuseContext) Params(key string, defaultValue ...string) string {
	switch {
	case slf.fiberCtx != nil:
		return slf.fiberCtx.Params(key, defaultValue...)

	case slf.grpcCtx != nil:
		val := ""
		if len(defaultValue) > 0 {
			val = defaultValue[0]
		}

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

func (slf *srFuseContext) Query(key string, defaultValue ...string) string {
	switch {
	case slf.fiberCtx != nil:
		return slf.fiberCtx.Query(key, defaultValue...)

	case slf.grpcCtx != nil:
		val := ""
		if len(slf.grpcCtx.queries) > 0 {
			if v, ok := slf.grpcCtx.queries[key]; ok {
				val = v
			}
		}
		return val

	default:
		panic("unimplemented")
	}
}

func (slf *srFuseContext) Parser(logc *clog.Instance, header, body interface{}) (bool, error) {
	mapHeader := slf.getHeader()

	if logc != nil {
		if v, ok := mapHeader["x-uid"]; ok {
			logc.UID = v
		}

		if v, ok := mapHeader["x-svcparent"]; ok {
			logc.SvcParent = v
		}
	}

	if header != nil {
		data, err := p9.Json.Marshal(mapHeader)
		if err == nil {
			err = p9.Json.UnMarshal(data, header)
		}

		if err != nil {
			err = p9.Err.WithStack(err, 1)
			return false, slf.r500InternalServerError(logc, err)
		}

		switch h := header.(type) {
		case *model.RequestHeader:
			if h.RFTimeRaw != "" {
				tm, err := time.Parse(time.RFC3339, h.RFTimeRaw)
				if err != nil {
					rfc3339 := "2006-01-02 15:04:05Z07:00"
					tm, err = time.Parse(rfc3339, h.RFTimeRaw)
				}

				if err == nil {
					h.RFTime = &tm
				}
			}

		default:
			v, err := p9.Util.ReflectionGet(header, "RequestHeader")
			if err == nil && v != nil {
				if h, ok := v.(model.RequestHeader); ok && h.RFTimeRaw != "" {
					tm, err := time.Parse(time.RFC3339, h.RFTimeRaw)
					if err != nil {
						rfc3339 := "2006-01-02 15:04:05Z07:00"
						tm, err = time.Parse(rfc3339, h.RFTimeRaw)
					}

					if err == nil {
						h.RFTime = &tm

						_ = p9.Util.ReflectionSet(header, map[string]interface{}{"RequestHeader": h})
					}
				}
			}
		}
	}

	if body != nil && slf.fiberCtx != nil {
		xP9BodyJson, ok := mapHeader["x-p9-body-json"]
		if ok && xP9BodyJson != "" {
			err := p9.Json.Decode(xP9BodyJson, &body)
			if err != nil {
				return false, slf.r500InternalServerError(logc, err)
			}
		} else {
			cType, ok := mapHeader["content-type"]
			if !ok {
				err := errors.New("unknown content-type")
				return false, slf.r500InternalServerError(logc, err)
			}

			if idx := strings.Index(cType, ";"); idx > -1 {
				cType = cType[0:idx]
			}

			switch cType {
			case "application/json":
				err := slf.fiberCtx.BodyParser(&body)
				if err != nil {
					err = p9.Err.WithStack(err, 1)
					return false, slf.r500InternalServerError(logc, err)
				}

				if body != nil {
					if json, err := p9.Json.Encode(body); err == nil {
						slf.jsonBody = &json
					}
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
					return false, slf.r500InternalServerError(logc, err)
				}

				if body != nil {
					if json, err := p9.Json.Encode(body); err == nil {
						slf.jsonBody = &json
					}
				}

			case "multipart/form-data":
				mf, err := slf.fiberCtx.MultipartForm()
				if err == nil {
					err = slf.mpfDecoder().Decode(body, mf.Value)
				}

				if err != nil {
					err = p9.Err.WithStack(err, 1)
					return false, slf.r500InternalServerError(logc, err)
				}

				slf.multipartFile = mf.File

				if body != nil {
					if json, err := p9.Json.Encode(body); err == nil {
						slf.jsonBody = &json
					}
				}
			}
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
				return false, slf.r500InternalServerError(logc, err)
			}
		}
	}

	return true, nil
}

func (slf *srFuseContext) GetMultipartFiles() map[string][]*multipart.FileHeader {
	return slf.multipartFile
}

func (slf *srFuseContext) ClientIP() string {
	if slf.clientIP == "" {
		slf.clientIP = cip.getClientIP(slf)
	}

	return slf.clientIP
}

func (slf *srFuseContext) Path() string {
	return slf.fiberCtx.Route().Path
}

func (slf *srFuseContext) Method() string {
	return slf.fiberCtx.Route().Method
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
	slf.isAuthSet = true
}

//endregion

// region util

func (slf *srFuseContext) SetCLog(logc *clog.Instance) {
	slf.logc = logc
}

func (slf *srFuseContext) GetCLog() *clog.Instance {
	if slf.logc == nil {
		slf.logc = clog.New()
	}

	return slf.logc
}

func (slf *srFuseContext) wrapError(err error) error {
	if err != nil {
		return fmt.Errorf("endpoint-path: %v\n%w", slf.path, err)
	}
	return fmt.Errorf("endpoint-path: %v", slf.path)
}

//endregion

//region response

func (slf *srFuseContext) GetResCode() (int, string) {
	return slf.resCode, slf.resSubCode
}

func (slf *srFuseContext) GetResObject() interface{} {
	return slf.resObject
}

func (slf *srFuseContext) SetResponse(code int, subCode string, obj interface{}) error {
	slf.resCode = code
	slf.resSubCode = subCode
	slf.resObject = obj
	return nil
}

func (slf *srFuseContext) SR200OK(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusOK, subCode, obj)
}

func (slf *srFuseContext) SR400BadRequest(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusBadRequest, subCode, obj)
}

func (slf *srFuseContext) SR401Unauthorized(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusUnauthorized, subCode, obj)
}

func (slf *srFuseContext) SR403Forbidden(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusForbidden, subCode, obj)
}

func (slf *srFuseContext) SR404NotFound(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusNotFound, subCode, obj)
}

func (slf *srFuseContext) SR406NotAcceptable(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusNotAcceptable, subCode, obj)
}

func (slf *srFuseContext) SR428PreconditionRequired(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusPreconditionRequired, subCode, obj)
}

func (slf *srFuseContext) SR500InternalServerError(subCode string, obj interface{}) error {
	return slf.SetResponse(http.StatusInternalServerError, subCode, obj)
}

func (slf *srFuseContext) RString(logc *clog.Instance, code int, data string) error {
	return slf.sendRawA(logc, code, data)
}

func (slf *srFuseContext) RJson(logc *clog.Instance, code int, data interface{}, opt ...FuseOpt) error {
	return slf.sendRawB(logc, code, data, opt...)
}

func (slf *srFuseContext) RJsonRaw(logc *clog.Instance, code int, data []byte) error {
	return slf.sendRawB(logc, code, f9.ToJsonRaw(data))
}

func (slf *srFuseContext) Redirect(logc *clog.Instance, code int, url string) error {
	return slf.sendRawA(logc, code, url)
}

func (slf *srFuseContext) RXXX(logc *clog.Instance, code int, status string, data interface{}, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:   code,
		Status: status,
		Data:   data,
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R200OK(logc *clog.Instance, data interface{}, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:   http.StatusOK,
		Status: fuseDefaultStatus.R200OK,
		Data:   data,
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R400BadRequest(logc *clog.Instance, message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusBadRequest,
		Status:  fuseDefaultStatus.R400BadRequest,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R400BadRequest),
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R401Unauthorized(logc *clog.Instance, message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusUnauthorized,
		Status:  fuseDefaultStatus.R401Unauthorized,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R401Unauthorized),
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R403Forbidden(logc *clog.Instance, message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusForbidden,
		Status:  fuseDefaultStatus.R403Forbidden,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R403Forbidden),
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R404NotFound(logc *clog.Instance, message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusNotFound,
		Status:  fuseDefaultStatus.R404NotFound,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R404NotFound),
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R406NotAcceptable(logc *clog.Instance, message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusNotAcceptable,
		Status:  fuseDefaultStatus.R406NotAcceptable,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R406NotAcceptable),
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R428PreconditionRequired(logc *clog.Instance, message string, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusPreconditionRequired,
		Status:  fuseDefaultStatus.R428PreconditionRequired,
		Message: f9.Ternary(message != "", message, fuseDefaultMessage.R428PreconditionRequired),
	}

	return slf.send(logc, fo, opt...)
}

func (slf *srFuseContext) R500InternalServerError(logc *clog.Instance, err error, opt ...FuseOpt) error {
	return slf.r500InternalServerError(logc, p9.Err.WithStack(slf.wrapError(err), 1), opt...)
}

func (slf *srFuseContext) r500InternalServerError(logc *clog.Instance, err error, opt ...FuseOpt) error {
	fo := FuseOpt{
		code:    http.StatusInternalServerError,
		Status:  fuseDefaultStatus.R500InternalServerError,
		Message: fuseDefaultMessage.R500InternalServerError,
		Error:   err,
	}

	return slf.send(logc, fo, opt...)
}

//endregion

//region send response

func (slf *srFuseContext) sendRawA(logc *clog.Instance, code int, data string) error {
	doSaveLog := func(resCode int, response interface{}, execFunc, execPath string, header, params map[string]string, endpoint, clientIp string) {
		var (
			severity   = clog.Info
			message    *string
			reqHeader  *string
			reqBody    = slf.jsonBody
			reqParam   *string
			resData    *string
			data       *string
			err        *string
			stackTrace *string
		)

		resCodeOne := fmt.Sprintf("%v", resCode)[:1]
		switch {
		case resCodeOne == "2":
			severity = clog.Info
		case resCodeOne == "4":
			severity = clog.Warning
		case resCodeOne == "5":
			severity = clog.Error
		}

		if value, err := p9.Json.Encode(header); err == nil && value != "{}" {
			reqHeader = &value
		}

		if value, err := p9.Json.Encode(params); err == nil && value != "{}" {
			reqParam = &value
		}

		if value, err := p9.Json.Encode(response); err == nil {
			resData = &value
		}

		m := clog.SendServiceModel{
			Endpoint:   endpoint,
			ExecFunc:   &execFunc,
			ExecPath:   &execPath,
			Message:    message,
			ReqHeader:  reqHeader,
			ReqBody:    slf.removeExcludedFieldReqBody(logc, reqBody),
			ReqParam:   reqParam,
			ResData:    resData,
			ResCode:    &resCode,
			Data:       data,
			Error:      err,
			StackTrace: stackTrace,
			ClientIP:   clientIp,
		}

		clog.SendService(0, *logc, severity, m, false)
	}

	saveLog := func(resCode int, response interface{}) {
		if logc != nil {
			depth := 3
			execFunc, execPath := p9.Util.GetExecutionInfo(depth)

			for {
				if !strings.Contains(execPath, "/project9/server/server-fuse-context.go") {
					break
				}

				depth++
				execFunc, execPath = p9.Util.GetExecutionInfo(depth)
			}

			header := slf.getHeader()
			params := slf.fiberCtx.AllParams()
			clientIp := f9.TernaryFnB(slf.clientIP != "", slf.clientIP, func() string { return cip.getClientIP(slf) })
			endpoint := strings.ToLower(fmt.Sprintf("%v:%v", slf.fiberCtx.Route().Method, slf.fiberCtx.Route().Path))
			go doSaveLog(resCode, response, execFunc, execPath, header, params, endpoint, clientIp)
		}
	}

	switch {
	case slf.fiberCtx != nil:
		saveLog(code, data)
		codeOne := fmt.Sprintf("%v", code)[:1]
		if codeOne == "3" {
			return slf.fiberCtx.Redirect(data, code)
		}
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

func (slf *srFuseContext) sendRawB(logc *clog.Instance, code int, data interface{}, opt ...FuseOpt) error {
	doSaveLog := func(resCode int, response interface{}, execFunc, execPath string, header, params map[string]string, endpoint, clientIp string) {
		var (
			severity   = clog.Info
			message    *string
			reqHeader  *string
			reqBody    = slf.jsonBody
			reqParam   *string
			resData    *string
			data       *string
			err        *string
			stackTrace *string
		)

		resCodeOne := fmt.Sprintf("%v", resCode)[:1]
		switch {
		case resCodeOne == "2":
			severity = clog.Info
		case resCodeOne == "4":
			severity = clog.Warning
		case resCodeOne == "5":
			severity = clog.Error
		}

		if value, err := p9.Json.Encode(header); err == nil && value != "{}" {
			reqHeader = &value
		}

		if value, err := p9.Json.Encode(params); err == nil && value != "{}" {
			reqParam = &value
		}

		if value, err := p9.Json.Encode(response); err == nil {
			resData = &value
		}

		m := clog.SendServiceModel{
			Endpoint:   endpoint,
			ExecFunc:   &execFunc,
			ExecPath:   &execPath,
			Message:    message,
			ReqHeader:  reqHeader,
			ReqBody:    slf.removeExcludedFieldReqBody(logc, reqBody),
			ReqParam:   reqParam,
			ResData:    resData,
			ResCode:    &resCode,
			Data:       data,
			Error:      err,
			StackTrace: stackTrace,
			ClientIP:   clientIp,
		}

		clog.SendService(0, *logc, severity, m, false)
	}

	saveLog := func(resCode int, response interface{}) {
		if logc != nil {
			depth := 3
			if len(opt) > 0 {
				for _, v := range opt {
					if v.LogDepthAdd > 0 {
						depth += v.LogDepthAdd
					}
				}
			}

			execFunc, execPath := p9.Util.GetExecutionInfo(depth)

			for {
				if !strings.Contains(execPath, "/project9/server/server-fuse-context.go") {
					break
				}

				depth++
				execFunc, execPath = p9.Util.GetExecutionInfo(depth)
			}

			header := slf.getHeader()
			params := slf.fiberCtx.AllParams()
			clientIp := f9.TernaryFnB(slf.clientIP != "", slf.clientIP, func() string { return cip.getClientIP(slf) })
			endpoint := strings.ToLower(fmt.Sprintf("%v:%v", slf.fiberCtx.Route().Method, slf.fiberCtx.Route().Path))
			go doSaveLog(resCode, response, execFunc, execPath, header, params, endpoint, clientIp)
		}
	}

	switch {
	case slf.fiberCtx != nil:
		saveLog(code, data)
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

func (slf *srFuseContext) send(logc *clog.Instance, fo FuseOpt, opt ...FuseOpt) error {
	if len(opt) > 0 {
		o := opt[0]
		fo.SubCode = f9.Ternary(o.SubCode != "", o.SubCode, fo.SubCode)
		fo.Status = f9.Ternary(o.Status != "", o.Status, fo.Status)
		fo.Message = f9.Ternary(o.Message != "", o.Message, fo.Message)
		fo.Address = f9.Ternary(o.Address != "", o.Address, fo.Address)
		fo.Error = f9.Ternary(o.Error != nil, o.Error, fo.Error)
		fo.MetaData = f9.Ternary(o.MetaData != nil, o.MetaData, fo.MetaData)
		fo.Data = f9.Ternary(o.Data != nil, o.Data, fo.Data)
	}

	response := FuseResponse{
		Meta: FuseResponseMeta{
			Code:    fo.code,
			SubCode: fo.SubCode,
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

		log.Printf("\n\n%v\n", response.Meta.Error)
	}

	doSaveLog := func(resCode int, response interface{}, execFunc, execPath string, header, params map[string]string, endpoint, clientIp string) {
		var (
			severity   = clog.Info
			message    *string
			reqHeader  *string
			reqBody    = slf.jsonBody
			reqParam   *string
			resData    *string
			data       *string
			err        *string
			stackTrace *string
		)

		resCodeOne := fmt.Sprintf("%v", resCode)[:1]
		switch {
		case resCodeOne == "2":
			severity = clog.Info
		case resCodeOne == "4":
			severity = clog.Warning
		case resCodeOne == "5":
			severity = clog.Error
		}

		if len(opt) > 0 {
			o := opt[0]
			if o.LogMessage != "" {
				message = &o.LogMessage
			}

			if o.LogData != "" {
				data = &o.LogData
			}
		}

		if value, err := p9.Json.Encode(header); err == nil && value != "{}" {
			reqHeader = &value
		}

		if value, err := p9.Json.Encode(params); err == nil && value != "{}" {
			reqParam = &value
		}

		if len(opt) > 0 && opt[0].StringCustomOutput {
			switch v := response.(type) {
			case string:
				resData = &v
			case *string:
				resData = v
			}
		}

		if resData == nil {
			if value, err := p9.Json.Encode(response); err == nil {
				resData = &value
			}
		}

		if fo.Error != nil {
			err = f9.Ptr(fo.Error.Error())
			stackTrace = f9.Ptr(fmt.Sprintf("%+v", fo.Error))
		}

		m := clog.SendServiceModel{
			Endpoint:   endpoint,
			ExecFunc:   &execFunc,
			ExecPath:   &execPath,
			Message:    message,
			ReqHeader:  reqHeader,
			ReqBody:    slf.removeExcludedFieldReqBody(logc, reqBody),
			ReqParam:   reqParam,
			ResData:    resData,
			ResCode:    &resCode,
			Data:       data,
			Error:      err,
			StackTrace: stackTrace,
			ClientIP:   clientIp,
		}

		clog.SendService(0, *logc, severity, m, false)
	}

	saveLog := func(resCode int, response interface{}) {
		if logc != nil {
			depth := 4
			execFunc, execPath := p9.Util.GetExecutionInfo(depth)

			for {
				if !strings.Contains(execPath, "/project9/server/server-fuse-context.go") {
					break
				}

				depth++
				execFunc, execPath = p9.Util.GetExecutionInfo(depth)
			}

			header := slf.getHeader()
			params := slf.fiberCtx.AllParams()
			clientIp := f9.TernaryFnB(slf.clientIP != "", slf.clientIP, func() string { return cip.getClientIP(slf) })
			endpoint := strings.ToLower(fmt.Sprintf("%v:%v", slf.fiberCtx.Route().Method, slf.fiberCtx.Route().Path))
			go doSaveLog(resCode, response, execFunc, execPath, header, params, endpoint, clientIp)
		}
	}

	restResponseRaw := func() error {
		saveLog(fo.code, fo.Data)

		if len(opt) > 0 && opt[0].StringCustomOutput {
			switch v := fo.Data.(type) {
			case string:
				return slf.fiberCtx.Status(fo.code).SendString(v)
			case *string:
				return slf.fiberCtx.Status(fo.code).SendString(*v)
			}
		}

		return slf.fiberCtx.Status(fo.code).JSON(fo.Data)
	}

	restResponse := func() error {
		if len(opt) == 0 || (len(opt) > 0 && len(opt[0].NewHeader) == 0 && len(opt[0].NewMeta) == 0) {
			saveLog(fo.code, response)

			if slf.sendResponse {
				return slf.fiberCtx.Status(fo.code).JSON(response)
			}

			slf.unSendResponse = &response
			return nil
		}

		var newMeta interface{} = response.Meta

		if len(opt[0].NewMeta) > 0 {
			meta, _ := p9.Conv.AnyToMap(response.Meta)
			for k, v := range opt[0].NewMeta {
				meta[k] = v
			}

			newMeta = meta
		}

		newResponse := map[string]interface{}{
			"meta": newMeta,
		}

		switch v := response.Data.(type) {
		case string:
			if len(v) > 0 {
				newResponse["data"] = v
			}

		default:
			if v != nil {
				newResponse["data"] = v
			}
		}

		for k, v := range opt[0].NewHeader {
			newResponse[k] = v
		}

		saveLog(fo.code, newResponse)
		if slf.sendResponse {
			return slf.fiberCtx.Status(fo.code).JSON(newResponse)
		}

		slf.unSendResponse = &response
		slf.unSendResponseOpt = opt
		return nil
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
				Message: response.Meta.Message,
				Address: response.Meta.Address,
				Error:   response.Meta.Error,
				Data:    metaData,
			},
		}

		return nil
	}

	switch {
	case slf.fiberCtx != nil && len(opt) > 0 && (opt[0].JsonCustomOutput || opt[0].StringCustomOutput):
		return restResponseRaw()
	case slf.fiberCtx != nil:
		return restResponse()
	case slf.grpcCtx != nil:
		return grpcResponse()
	}

	panic("unimplemented")
}

func (slf *srFuseContext) grpcSend(err error) (*grf.Response, error) {
	return slf.grpcCtx.response, err
}

func (slf *srFuseContext) removeExcludedFieldReqBody(logc *clog.Instance, jsonReqBody *string) *string {
	if logc != nil && len(logc.ReqBodySaveExcluded) > 0 && jsonReqBody != nil && len(*jsonReqBody) > 0 {
		out := make(map[string]interface{}, 0)
		err := p9.Json.Decode(*jsonReqBody, &out)
		if err == nil {
			haveDeleted := false

			for _, field := range logc.ReqBodySaveExcluded {
				data := out[field]
				if data != nil {
					switch val := data.(type) {
					case *string:
						if len(*val) > 0 {
							haveDeleted = true
							delete(out, field)
						}

					case string:
						if len(val) > 0 {
							haveDeleted = true
							delete(out, field)
						}
					}
				}
			}

			if haveDeleted {
				jsonBody, err := p9.Json.Encode(out)
				if err == nil {
					return &jsonBody
				}
			}
		}
	}

	return jsonReqBody
}

//endregion
