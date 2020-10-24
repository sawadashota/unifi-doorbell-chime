SHELL := /bin/bash
.DEFAULT_GOAL := help

GO := go
GOTEST := $(or $(GOTEST),$(GO) test)
CONFIG_FILE := $(or ${CONFIG_FILE}, ".unifi-doorbell-chime.yaml")
BINARY_FILE := unifi-doorbell-chime

NPM := npm
NPM_PREFIX := web/frontend

dev: web/build ## start dev
	$(GO) run main.go --config $(CONFIG_FILE)

build: web/build ## build for production
	$(GO) build -o $(BINARY_FILE) .

start: build ## exec built binary
	./$(BINARY_FILE) --config $(CONFIG_FILE)

web/build: ## build web frontend
	$(NPM) run build --prefix $(NPM_PREFIX)

clean: ## clean built and dependencies
	rm $(BINARY_FILE)
	rm -rf ./web/node_modules ./web/build

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

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
