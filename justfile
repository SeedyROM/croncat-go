test:
    go test -json -v ./... | gotestfmt

coverage:
    courtney ./...
    go tool cover -func=coverage.out
    go tool cover -html=coverage.out -o coverage.html

run *FLAGS:
    go run ./cli/croncatd {{FLAGS}}