# Generate mocks
mock:
    GOFLAGS="-buildvcs=false" mockery

# Run tests with coverage
test: mock
    go test -cover -bench=. -benchmem -race ./... -coverprofile=coverage.out

# Build sesh binary to GOPATH/bin
build version="dev":
    go build -buildvcs=false -ldflags "-X 'main.version={{version}}'" -o `go env GOPATH`/bin/sesh
