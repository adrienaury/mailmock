# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

SHELL := /bin/bash # Use bash syntax

# Build variables
BUILD_DIR ?= .target
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
BUILD_BY ?= $(shell git config user.email)
LDFLAGS += -X main.version=${VERSION} -X main.commit=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE} -X main.builtBy=${BUILD_BY}

# Project variables
DOCKER_IMAGE = adrienaury/mailmock
DOCKER_TAG ?= $(shell echo -n ${VERSION} | sed -e 's/[^A-Za-z0-9_\\.-]/_/g')
RELEASE := $(shell [[ $(VERSION) =~ ^[0-9]*.[0-9]*.[0-9]*$$ ]] && echo 1 || echo 0 )
MAJOR := $(shell echo $(VERSION) | cut -f1 -d.)
MINOR := $(shell echo $(VERSION) | cut -f2 -d.)
PATCH := $(shell echo $(VERSION) | cut -f3 -d. | cut -f1 -d-)

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: info
info: ## Prints build informations
	@echo COMMIT_HASH=$(COMMIT_HASH)
	@echo VERSION=$(VERSION)
	@echo RELEASE=$(RELEASE)
ifeq (${RELEASE}, 1)
	@echo MAJOR=$(MAJOR)
	@echo MINOR=$(MINOR)
	@echo PATCH=$(PATCH)
endif
	@echo DOCKER_IMAGE=$(DOCKER_IMAGE)
	@echo DOCKER_TAG=$(DOCKER_TAG)
	@echo BUILD_BY=$(BUILD_BY)

.PHONY: clean
clean: ## Clean builds
	rm -rf ${BUILD_DIR}/

.PHONY: mkdir
mkdir:
	mkdir -p ${BUILD_DIR}

.PHONY: tidy
tidy: ## Add missing and remove unused modules
	GO111MODULE=on go mod tidy

.PHONY: lint
lint: ## Examines Go source code and reports suspicious constructs
ifeq (, $(shell which golangci-lint))
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.18.0
endif
	golangci-lint run -E misspell -E gocyclo -E gosec -E unparam -E goimports -E nakedret -E gocritic

.PHONY: build-%
build-%: mkdir
	GO111MODULE=on go build ${GOARGS} -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

.PHONY: build
build: $(patsubst cmd/%,build-%,$(wildcard cmd/*)) ## Build all binaries

.PHONY: test
test: mkdir ## Run all tests with coverage
	GO111MODULE=on go test -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic ./...

.PHONY: run-%
run-%: build-%
	${BUILD_DIR}/$*

.PHONY: run
run: $(patsubst cmd/%,run-%,$(wildcard cmd/*)) ## Build and execute a binary

.PHONY: release-%
release-%: mkdir
	GO111MODULE=on go build ${GOARGS} -ldflags "-w -s ${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

.PHONY: release
release: clean info lint $(patsubst cmd/%,release-%,$(wildcard cmd/*)) ## Build all binaries for production

.PHONY: docker
docker: info ## Build docker image locally
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} --build-arg VERSION=${VERSION} --build-arg BUILD_BY=${BUILD_BY} .
ifeq (${RELEASE}, 1)
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${MAJOR}.${MINOR}
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${MAJOR}
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
endif

.PHONY: start
start: docker ## Start a docker container with default parameters
	docker run -ti --rm -p 1080:80 -p 1025:25 ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: push
push: docker ## Push docker image on DockerHub
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
ifeq (${RELEASE}, 1)
	docker push ${DOCKER_IMAGE}:${MAJOR}.${MINOR}
	docker push ${DOCKER_IMAGE}:${MAJOR}
	docker push ${DOCKER_IMAGE}:latest
endif
