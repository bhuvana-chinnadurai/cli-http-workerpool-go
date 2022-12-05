.DEFAULT_GOAL := build
BIN_FILE=myhttp
all: build test

build:
	go build -o "${BIN_FILE}" cmd/main.go
clean:
	go clean
	rm ${BIN_FILE}
test:
	go test -v ./...
test_coverage: 
	go test ./... -coverprofile=coverage.out

test_coverage_html: 
	go tool cover -html=coverage.out
