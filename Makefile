SHELL := /bin/bash
.DEFAULT_GOAL := help

GO := $(or $(GOBIN), go)
GOTEST := $(or $(GOTEST),$(GO) test)
CONFIG_FILE := $(or ${CONFIG_FILE}, "contrib/config.yaml")
BINARY_FILE := unifi-doorbell-chime

NPM := npm
NPM_PREFIX := web/frontend

dev: web/build ## start dev
	$(GO) run main.go start --config $(CONFIG_FILE)

prebuild: web/build ## pre build

build: prebuild ## build for production
	$(GO) build -o $(BINARY_FILE) .

start: build ## exec built binary
	./$(BINARY_FILE) start --config $(CONFIG_FILE)

web/build: ## build web frontend
	$(NPM) run build --prefix $(NPM_PREFIX)

clean: ## clean built and dependencies
	rm $(BINARY_FILE)
	rm -rf ./web/frontend/node_modules ./web/frontend/static

install: ## install dependencies
	$(GO) mod download -x
	$(NPM) install --prefix $(NPM_PREFIX)

test : go/test npm/test ## Run all tests

go/test: ## run go test
	$(GOTEST) -v -cover ./...

npm/test: ## run npm test
	$(NPM) run test --prefix $(NPM_PREFIX)

npm/dev: ## start npm dev server
	$(NPM) start --prefix $(NPM_PREFIX)

lint: go/lint npm/lint ## check lint

go/lint: ## check lint Go code
	golangci-lint run

npm/lint: ## check lint node code
	$(NPM) run lint --prefix $(NPM_PREFIX)

npm/format: ## format node code
	$(NPM) run format --prefix $(NPM_PREFIX)

go/version: ## print Go version
	$(GO) version

go/tidy: ## print Go version
	$(GO) mod tidy

Tag=
release: ## Release by goreleaser
	git tag ${Tag}
	git push origin ${Tag}
	goreleaser --rm-dist

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
