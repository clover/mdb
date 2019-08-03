GO_EXECUTABLE ?= go
VERSION ?= $(shell git describe --tags --long --always --dirty)
OUTPUT_BIN = mdb
GOPATH ?= ${HOME}/go
INSTALL_DEST ?= /usr/local/bin

build:
	@# GOPATH=${GOPATH} go get -d ./... && GOPATH=${GOPATH} ${GO_EXECUTABLE} build -o ${OUTPUT_BIN} -ldflags "-X main.version=${VERSION}"
	@# GOPATH=${GOPATH} ${GO_EXECUTABLE} build -o ${OUTPUT_BIN} -ldflags "-X main.version=${VERSION}"
	${GO_EXECUTABLE} build -o ${OUTPUT_BIN} -mod=vendor -ldflags "-X main.version=${VERSION}"

run: build
	./${OUTPUT_BIN}

install: build
	sudo cp ${OUTPUT_BIN} ${INSTALL_DEST}/${OUTPUT_BIN}
