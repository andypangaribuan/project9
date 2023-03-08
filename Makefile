# enable make
all: help

help:
	@cat Makefile \
		| grep -B1 -E '^[a-zA-Z0-9_.-]+:.*' \
		| grep -v -- -- \
		| sed 'N;s/\n/###/' \
		| sed -n 's/^#: \(.*\)###\(.*\):.*/\2###→ \1/p' \
		| column -t -s '###' \
		| sed 's/.*→ space.*//g'

#: update module
tidy:
	@go mod tidy

#: find error
vet:
	@go vet

#: check unformatted files
format:
	@gofmt -l .

#: space
.:

#: generate grpc proto files
grpc-generate:
	@protoc --go_out=. --go-grpc_out=. ./server/proto/src/*.proto
