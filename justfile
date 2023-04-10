test:
    go test -json -v ./... | gotestfmt

coverage:
    courtney ./...
    go tool cover -func=coverage.out
    go tool cover -html=coverage.out -o coverage.html

build *FLAGS:
    go build -o bin/croncatd -ldflags="-s -w" {{FLAGS}} ./cmd/croncatd

run *FLAGS:
    go run ./cmd/croncatd {{FLAGS}}