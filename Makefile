.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

.PHONY: oapi-gen
oapi-gen: ## Generate open api
	oapi-codegen -old-config-style -generate "chi-server" -package api openapi.yaml > user-service/api/server.go
	oapi-codegen -old-config-style -generate "spec" -package api openapi.yaml > user-service/api/spec.go
	oapi-codegen -old-config-style -generate "types" -package api openapi.yaml > user-service/api/model.go

.PHONY: vendor
vendor: ## Download Vendor packages
	go mod vendor

.PHONY: lint-install
lint-install:
	GO111MODULE=off go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: lint
lint: lint-install## Run lint
	golangci-lint run ./...

.PHONY: unit-tests
unit-tests: ## Run tests
	go test ./...

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t user-service .

.PHONY: docker-run
docker-run: docker-build ## Run containerized application
	docker run -p 8080:8080 user-service

.PHONY: docker-build-debug
docker-build-debug: ## Build Docker image with debugging port
	docker build -t user-service-debug -f ./Dockerfile.debug .

.PHONY: docker-run-debug
docker-run-debug: docker-build-debug ## Run containerized application with debugging port
	docker run -p 40000:40000 --name user-service-debug --cap-add SYS_PTRACE --network database-network --security-opt apparmor=unconfined user-service-debug:latest

.PHONY: docker-push
docker-push: docker-build ## Push Docker image
	docker login -u $DOCKER_USER -p $DOCKER_PASS
	docker push user-service:latest
