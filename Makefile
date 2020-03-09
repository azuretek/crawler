# go option
GO             ?= go
APP_VERSION    ?= v0.0.0-dev
LDFLAGS        := -w -s
GOFLAGS        :=
GO_EXTRA_FLAGS := -v
TAGS           :=
BINDIR         := $(CURDIR)/bin
PKGDIR         := github.com/azuretek/crawler
CGO_ENABLED    := 0

.PHONY: pre-commit
pre-commit: fmt

.PHONY: build
build: gobuild

.PHONY: test
test:
	go test -coverprofile cp.out ./...

.PHONY: dep
dep:
	@echo " ===> Installing dependencies <=== "
	go mod vendor

.PHONY: gobuild
gobuild:
	@echo " ===> building releases in ./bin/... <=== "
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BINDIR)/worker -ldflags "$(LDFLAGS)" $(GO_EXTRA_FLAGS) $(PKGDIR)/cmd/worker
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BINDIR)/apiserver -ldflags "$(LDFLAGS)" $(GO_EXTRA_FLAGS) $(PKGDIR)/cmd/apiserver
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BINDIR)/testserver -ldflags "$(LDFLAGS)" $(GO_EXTRA_FLAGS) $(PKGDIR)/internal/testserver

.PHONY: docker-build
docker-build:
	@echo " ===> building docker image <==="
	@DOCKER_BUILDKIT=1 docker build -t crawler:${APP_VERSION} -f Dockerfile . --build-arg VERSION=${APP_VERSION}

.PHONY: fmt
fmt:
	@echo " ===> Running goimports <==="
	find  . -path ./vendor -prune -o -type f -name '*.go' -print | xargs -n 1 goimports -w
