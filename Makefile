SHELL := /bin/bash
GO_VERSION=$(shell go mod edit -json | jq -r '.Go')

BUILD_DIR=./build
SHELL := /bin/bash
DOCKER_BUILD_IMAGE=mexfoo/mononoke-go
DOCKER_WORKDIR=/proj
DOCKER_GO_BUILD=go build -mod=readonly -a -installsuffix cgo -ldflags "$$LD_FLAGS"
DOCKER_TEST_LEVEL ?= 0 # Optionally run a test during docker build

require-version:
	if [ -n ${VERSION} ] && [[ $$VERSION == "v"* ]]; then echo "The version may not start with v" && exit 1; fi
	if [ -z ${VERSION} ]; then echo "Need to set VERSION" && exit 1; fi;

test-coverage:
	go test --race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... ./...

format:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

check-go:
	golangci-lint run

package-zip:
	for BUILD in $(shell find ${BUILD_DIR}/*); do \
       zip -j $$BUILD.zip $$BUILD ./LICENSE; \
    done

build-docker: require-version
	docker buildx build \
		$(if $(DOCKER_BUILD_PUSH),--push) \
		-t ${DOCKER_IMAGE_REPO}${DOCKER_ORG}/mononoke-go:latest \
		-t ${DOCKER_IMAGE_REPO}${DOCKER_ORG}/mononoke-go:${VERSION} \
		-t ${DOCKER_IMAGE_REPO}${DOCKER_ORG}/mononoke-go:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ${DOCKER_IMAGE_REPO}${DOCKER_ORG}/mononoke-go:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		--build-arg RUN_TESTS=$(DOCKER_TEST_LEVEL) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg LD_FLAGS="$$LD_FLAGS" \
		--platform linux/amd64 \
		-f Dockerfile .

build-linux-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$$LD_FLAGS" -o ${BUILD_DIR}/mononoke-go

build-linux-arm-7:
	CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$$LD_FLAGS" -o ${BUILD_DIR}/mononoke-go-linux-arm-7

build-linux-arm64:
	CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags "$$LD_FLAGS" -o ${BUILD_DIR}/mononoke-go-linux-arm64

build-windows-amd64:
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "$$LD_FLAGS" -o ${BUILD_DIR}/mononoke-go-windows-amd64.exe

build: build-linux-arm-7 build-linux-amd64 build-linux-arm64 build-windows-amd64

.PHONY: test-coverage check-go package-zip build-docker build