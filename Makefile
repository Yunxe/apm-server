ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
OUTPUT_DIR := $(CURDIR)/_output

.PHONY: build
build:
	go build -o ${OUTPUT_DIR}/apmserver -v ./cmd/main.go

.PHONY: run
run: build
	 ${OUTPUT_DIR}/apmserver -c ./configs/apm-server.yaml

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: inittopic
inittopic: build
	${OUTPUT_DIR}/main init-topic -t true -c ./configs/apm-server.yaml


