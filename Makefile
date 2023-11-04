.DEFAULT_GOAL:=help

IMG_REGISTRY?=ttl.sh
IMG_REPOSITORY?=${USER}
PROJECT?=spectro-rapid-agent
IMG_TAG?=latest
IMG=${IMG_REGISTRY}/${IMG_REPOSITORY}/${PROJECT}:${IMG_TAG}
TARGETARCH ?= amd64
TARGETOS ?= linux

BIN_DIR ?= ./bin
BUILD_DIR ?= ./_build
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell go env GOPATH)

.PHONY: build
build:
	go build -o bin/rapid-agent ./cmd/ra

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: test-toolchain
test-toolchain:
	GOTOOLCHAIN=go1.21.0 go test -v ./... -cover

# TODO: configure golangci-lint
.PHONY: lint 
lint: 
	go fmt ./...
	go vet ./...
	go mod tidy
	go mod verify
	golangci-lint run ./...

.PHONY: run
run:
	go run ./cmd/ra

.PHONY: docker-build
docker-build:
	docker buildx build --platform ${TARGETOS}/${TARGETARCH} -t ${IMG} -f deploy/Dockerfile --no-cache .

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: docker
docker: docker-build docker-push

.PHONY: scan
scan: scan-fs scan-image

.PHONY: scan-fs
scan-fs:
	trivy fs \
		--scanners vuln \
		--exit-code=1 \
		--severity="CRITICAL,HIGH,MEDIUM" \
		--ignore-unfixed \
		--ignorefile .trivyignore \
		./

.PHONY: scan-image
scan-image:
	trivy image \
		--scanners vuln \
		--exit-code=1 \
		--severity="CRITICAL,HIGH,MEDIUM" \
		--ignore-unfixed \
		--ignorefile .trivyignore \
		${IMG}

# TODO: configure swagger
.PHONY: swagger-spec
swagger-spec: export SHELL = /bin/sh
swagger-spec:
	set -e
	mkdir -p ./spec/

	swagger generate spec ./pkg/handlers -o ./spec/ra.json
	swagger validate ./spec/ra.json
