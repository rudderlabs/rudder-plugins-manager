install:
		sh install-hooks.sh
default: build

test: install-tools test-run

test-run:
		gotestsum --format pkgname-and-test-fails -- -count=1 -shuffle=on  -coverprofile=coverage.txt -vet=all ./...

build:
		go build -o bin/$(NAME) ./cmd/$(NAME).go

test-with-coverage: test coverage

coverage:
	go tool cover -html=coverage.txt -o coverage.html

install-tools:
	go install mvdan.cc/gofumpt@latest
	go install gotest.tools/gotestsum@v1.8.2

.PHONY: lint
lint: fmt ## Run linters on all go files
	docker run --rm -v $(shell pwd):/app:ro -w /app golangci/golangci-lint:v1.51.1 bash -e -c \
		'golangci-lint run -v --timeout 5m'

.PHONY: fmt
fmt: install-tools ## Formats all go files
	gofumpt -l -w -extra  .
