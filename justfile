# Generate mocks
mock:
    GOFLAGS="-buildvcs=false" go tool mockery

# Run tests with coverage
test: mock
    go test -cover -race ./... -coverprofile=coverage.out

# Run benchmarks (no -race: it measures the race detector, not the code)
bench: mock
    go test -run='^$' -bench=. -benchmem ./lister/... ./picker/...

# Build sesh binary to GOPATH/bin
build version="dev":
    go build -buildvcs=false -ldflags "-X 'main.version={{version}}'" -o `go env GOPATH`/bin/sesh

# Generate man page
man: build
    mkdir -p share/man/man1
    sesh man > share/man/man1/sesh.1
