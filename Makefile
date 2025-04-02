.PHONY: test build

BUILD_FLAGS="-X 'main.version=`git describe --tags --abbrev=0`'"

test:
	@go test -cover -bench=. -benchmem -race ./... -coverprofile=coverage.out

build: 
	@go build -ldflags ${BUILD_FLAGS} -o $(shell go env GOPATH)/bin/sesh
