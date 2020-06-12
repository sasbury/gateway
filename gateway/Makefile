ifndef GOPATH
GOPATH := $(HOME)/go
endif

BUILD_DIR := build
EXECUTABLE_NAME := gateway
EXECUTABLE_DOCKER_NAME := ${EXECUTABLE_NAME}_docker

build:
	go build -a -o ../$(BUILD_DIR)/$(EXECUTABLE_NAME) -v .

build-docker:
	CGO_ENABLED=0 GOOS=linux go build -a -o ../$(BUILD_DIR)/$(EXECUTABLE_DOCKER_NAME) -v .

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
