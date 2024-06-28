.PHONY: test build

test:
	go test -cover -bench=. -benchmem -race ./... -coverprofile=coverage.out

build: 
	go build -o $(shell echo $$GOPATH)/bin/sesh-dev
