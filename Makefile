NAME := oauth
SRCS := $(shell find . -type d -name vendor -prune -o -type f -name "*.go" -print)
VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\"" 
DIST_DIRS := find * -type d -exec
DIST_NO_WINDOWS_DIRS := find * -type d -name "window*" -exec
DIST_WINDOWS_DIRS := find * -type d -name "window*" -exec

.DEFAULT_GOAL := bin/$(NAME)

.PHONY: test
test:
	@go test -cover -v `glide novendor`

.PHONY: install
install:
	@cp bin/$(NAME) /usr/local/bin/$(NAME)

.PHONY: uninstall
uninstall:
	@rm /usr/local/bin/$(NAME)

.PHONY: clean
clean:
	@rm -rf bin/*
	@rm -rf vendor/*
	@rm -rf dist/*

.PHONY: dist-clean
dist-clean: clean
	@rm -f $(NAME).tar.gz

.PHONY: format
format: import
	-@goimports -w .
	@gofmt -w .

.PHONY: cross-build
cross-build: deps
	@for os in darwin linux windows; do \
	    for arch in amd64 386; do \
	        GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo \
	        -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME) .; \
	    done; \
	done

.PHONY: import
ifeq ($(shell command -v goimports 2> /dev/null),)
	go get golang.org/x/tools/cmd/goimports
endif


.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: deps
deps: glide
	glide install

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) .

.PHONY: dist
dist:
	@cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) cp ../completions/zsh/_$(NAME) {} \; && \
	$(DIST_NO_WINDOWS_DIRS) tar zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_WINDOWS_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \;

