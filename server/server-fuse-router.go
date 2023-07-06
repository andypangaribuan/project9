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

	ep := strings.TrimSpace(path[index+1:])

	switch path[0:index] {
	case "GET":
		slf.fiberApp.Get(ep, slf.restfulProcess(ep, handlers...))
	case "POS":
		slf.fiberApp.Post(ep, slf.restfulProcess(ep, handlers...))
	case "DEL":
		slf.fiberApp.Delete(ep, slf.restfulProcess(ep, handlers...))
	case "PUT":
		slf.fiberApp.Put(ep, slf.restfulProcess(ep, handlers...))
	case "PAT":
		slf.fiberApp.Patch(ep, slf.restfulProcess(ep, handlers...))
	default:
		panic(panicMsg)
	}

	slf.grpcProcess(path, handlers...)
}

func (slf *srFuseRouter) Group(endpoints map[string][]func(sc FuseContext) error) {
	for path, handlers := range endpoints {
		slf.Single(path, handlers...)
	}
}

func (slf *srFuseRouter) restfulProcess(path string, handlers ...func(FuseContext) error) func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := &srFuseContext{fiberCtx: c, path: path, sendResponse: true}
		ctx.reqCtx = &srFuseContextRequest{
			fuseCtx: ctx,
		}

		if len(handlers) >= 2 {
			err := handlers[0](ctx)
			if !ctx.isAuthSet {
				return err
			}

			return handlers[1](ctx)
		}

		return handlers[0](ctx)
	}
}

func (slf *srFuseRouter) grpcProcess(path string, handlers ...func(FuseContext) error) {
	slf.fuseGrpc.routes[path] = handlers
}
