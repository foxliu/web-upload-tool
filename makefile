# include .env
PROJECTNAME=$(shell basename "$(PWD)")

BINARY_PATH=./bin
BINARY=${BINARY_PATH}/webUpload
VERSION=0.0.1
BUILD=`date +%FT%T%z`

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

list:
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}

fmt:
	@gofmt -s -w ${GOFILES}

default:
	@go build -o ${BINARY} -tags=jsoniter main.go

build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}-linux -tags=jsoniter main.go

run:
	@go run main.go

clean:
	@-if [ -d ${BINARY_PATH} ] ; then rm -f ${BINARY_PATH}/* ; fi


.PHONY: default fmt fmt-check install test vet docker clean start