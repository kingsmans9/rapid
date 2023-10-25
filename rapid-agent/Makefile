.DEFAULT_GOAL:=help

.PHONY: build
build:
	go build -o ./bin/ra ./cmd/ra

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
	docker build -t rapid-agent -f deploy/Dockerfile --no-cache .

.PHONY: docker-build docker-compose-up
docker-compose-up:
	docker-compose -f deploy/docker-compose.yaml up -d

.PHONY: docker-compose-down
docker-compose-down:
	docker-compose -f deploy/docker-compose.yaml down

.PHONY: scan
scan: scan-local scan-docker

.PHONY: scan-local
scan-local:
	trivy fs \
		--scanners vuln \
		--exit-code=1 \
		--severity="CRITICAL,HIGH,MEDIUM" \
		--ignore-unfixed \
		--ignorefile .trivyignore \
		./

.PHONY: scan-docker
scan-docker: docker-build
	trivy image \
		--scanners vuln \
		--exit-code=1 \
		--severity="CRITICAL,HIGH,MEDIUM" \
		--ignore-unfixed \
		--ignorefile .trivyignore \
		ra:latest

# TODO: configure swagger
.PHONY: swagger-spec
swagger-spec: export SHELL = /bin/sh
swagger-spec:
	set -e
	mkdir -p ./spec/

	swagger generate spec ./pkg/handlers -o ./spec/ra.json
	swagger validate ./spec/ra.json
