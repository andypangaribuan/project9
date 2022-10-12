/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (slf *srFuseRouter) Single(path string, handlers ...func(sc FuseContext) error) {
	if len(handlers) < 1 && len(handlers) > 2 {
		panic("First handler must be AUTH_CONTROLLER, then the second is LOGIC_CONTROLLER, or LOGIC_CONTROLLER only")
	}

	panicMsg := "Please use ▶︎ GET, POS, DEL, PUT or PAT"
	index := strings.Index(path, ":")
	if index == -1 {
		panic(panicMsg)
	}

	handler := handlers[0]
	ep := strings.TrimSpace(path[index+1:])
	auth := make([]func(FuseContext) error, 0)

	if len(handlers) == 2 {
		auth = append(auth, handlers[1])
	}

	switch path[0:index] {
	case "GET":
		slf.fiberApp.Get(ep, slf.restfulProcess(ep, handler, auth...))
	case "POS":
		slf.fiberApp.Post(ep, slf.restfulProcess(ep, handler, auth...))
	case "DEL":
		slf.fiberApp.Delete(ep, slf.restfulProcess(ep, handler, auth...))
	case "PUT":
		slf.fiberApp.Put(ep, slf.restfulProcess(ep, handler, auth...))
	case "PAT":
		slf.fiberApp.Patch(ep, slf.restfulProcess(ep, handler, auth...))
	default:
		panic(panicMsg)
	}

	slf.grpcProcess(path, handler, auth...)
}

func (slf *srFuseRouter) Group(endpoints map[string][]func(sc FuseContext) error) {
	for path, handlers := range endpoints {
		slf.Single(path, handlers...)
	}
}

func (slf *srFuseRouter) restfulProcess(path string, handler func(FuseContext) error, auth ...func(FuseContext) error) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		context := &srFuseContext{fiberCtx: ctx, path: path}

		if len(auth) == 1 {
			err := auth[0](context)
			if err != nil {
				return err
			}
		}

		return handler(context)
	}
}

func (slf *srFuseRouter) grpcProcess(path string, handler func(FuseContext) error, auth ...func(FuseContext) error) {
	handlers := make([]func(FuseContext) error, 0)
	handlers = append(handlers, handler)
	if len(auth) > 0 {
		handlers = append(handlers, auth...)
	}

	slf.fuseGrpc.routes[path] = handlers
}
