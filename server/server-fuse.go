/*
 * Copyright (c) 2022.
 * Created by Andy Pangaribuan <https://github.com/apangaribuan>.
 * All Rights Reserved.
 */

package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/andypangaribuan/project9/f9"
	"github.com/andypangaribuan/project9/p9"
	"github.com/andypangaribuan/project9/server/proto/gen/grf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"google.golang.org/grpc"
)

var (
	fuseFiberApp       *fiber.App
	fuseDefaultStatus  *FuseDefaultStatus
	fuseDefaultMessage *FuseDefaultMessage
	fuseWithError      *bool
)

func init() {
	fuseWithError = f9.Ptr(false)

	fuseDefaultStatus = &FuseDefaultStatus{
		R200OK:                   "success",
		R400BadRequest:           "bad-request",
		R401Unauthorized:         "unauthorized",
		R403Forbidden:            "forbidden",
		R404NotFound:             "not-found",
		R406NotAcceptable:        "not-acceptable",
		R428PreconditionRequired: "precondition-required",
		R500InternalServerError:  "internal-server-error",
	}

	fuseDefaultMessage = &FuseDefaultMessage{
		R400BadRequest:           "Something went wrong while processing your request.",
		R401Unauthorized:         "Bad credentials. Please login again.",
		R403Forbidden:            "You do not have permission for this request/resource",
		R404NotFound:             "The resource you are looking for is not available.",
		R406NotAcceptable:        "The resource you are looking for is not acceptable.",
		R428PreconditionRequired: "The resource you are looking for is conditional to fulfill.",
		R500InternalServerError:  "We apologize and are fixing the problem. Please try again at a later stage.",
	}
}

func SetFuseStatusMessage(call func(status *FuseDefaultStatus, message *FuseDefaultMessage)) {
	call(fuseDefaultStatus, fuseDefaultMessage)
}

func Fuse(restfulPort, grpcPort int, autoRecover bool, withErr bool, routes func(router FuseRouter)) {
	Fuse2(restfulPort, grpcPort, autoRecover, withErr, routes)
}

func Fuse2(restfulPort, grpcPort int, autoRecover bool, withErr bool, routes func(router FuseRouter), grpcRegister ...func(svc *grpc.Server)) {
	if restfulPort == grpcPort {
		panic("restfulPort and grpcPort cannot have same value")
	}

	if restfulPort != -1 && isPortUse(restfulPort) {
		panic("restfulPort already in use")
	}

	if grpcPort != -1 && isPortUse(grpcPort) {
		panic("grpcPort already in use")
	}

	fuseWithError = &withErr

	fuseFiberApp = fiber.New(fiber.Config{
		JSONEncoder: p9.Json.Marshal,
		JSONDecoder: p9.Json.UnMarshal,
	})

	if restfulPort != -1 && autoRecover {
		fuseFiberApp.Use(recover.New())
	}

	router := &srFuseRouter{
		fiberApp: fuseFiberApp,
		fuseGrpc: &srFuseGrpc{
			routes: make(map[string][]func(FuseContext) error, 0),
		},
	}

	routes(router)

	if grpcPort != -1 {
		go func() {
			var register func(svc *grpc.Server)

			if len(grpcRegister) > 0 {
				register = grpcRegister[0]
			} else {
				register = func(server *grpc.Server) {
					grf.RegisterRestfulServiceServer(server, router.fuseGrpc)
				}
			}

			time.Sleep(time.Millisecond * 10)
			p9.Server.StartGRPC(grpcPort, autoRecover, register, 3)
		}()
	}

	if restfulPort != -1 {
		log.Fatal(fuseFiberApp.Listen(fmt.Sprintf(":%v", restfulPort)))
	}
}

func isPortUse(port int) bool {
	host := "127.0.0.1"
	timeout := time.Second * 3

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)
	if err != nil {
		return false
	}

	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
		return true
	}

	return false
}
