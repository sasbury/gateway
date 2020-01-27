
.PHONY: build test

build: compile

compile:
	go build ./...

install:
	go install ./...

cover: test
	go tool cover -html=./cover.out

test:
	go vet ./...
	rm -rf ./cover.out
	go test -race -coverpkg=./... -coverprofile=./cover.out ./...

fast:
	go test -race -failfast ./...
