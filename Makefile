all: help

help:
	@echo '✦ grpc-generate'

grpc-generate:
	@protoc --go_out=. --go-grpc_out=. ./server/proto/src/*.proto
