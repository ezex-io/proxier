BINARY_NAME = proxier
BUILD_DIR = build
CMD_DIR = cmd/proxier/main.go

########################################
### Targets needed for development

docker:
	docker build --tag ezex-gateway .

########################################
### Building

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

release:
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-s -w" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf $(BUILD_DIR)

########################################
### Testing

unit_test:
	@echo "Running unit tests..."
	@go test ./...

race_test:
	@echo "Running race condition tests..."
	@go test ./... -race

test: unit_test race_test

########################################
### Formatting the code

fmt:
	@echo "Formatting code..."
	@go tool gofumpt -l -w .

lint:
	@echo "Running lint..."
	@go tool golangci-lint  run ./... --timeout=20m0s

check: fmt lint

.PHONY: gen-graphql docker
.PHONY: build release clean
.PHONY: test unit_test race_test
.PHONY: fmt lint check
