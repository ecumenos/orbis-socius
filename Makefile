GOPRIVATE=github.com/ecumenos
SHELL=/bin/sh

.PHONY: all
all: check fmt lint test

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-short
test-short: ## Run tests, skipping slower integration tests
	go test -test.short ./...

.PHONY: test-interop
test-interop: ## Run tests, including local interop (requires services running)
	go clean -testcache && go test -tags=localinterop ./...

.PHONY: coverage-html
coverage-html: ## Generate test coverage report and open in browser
	go test ./... -coverpkg=./... -coverprofile=test-coverage.out
	go tool cover -html=test-coverage.out

.PHONY: lint
lint: ## Verify code style and run static checks
	go vet -asmdecl -assign -atomic -bools -buildtag -cgocall -copylocks -httpresponse -loopclosure -lostcancel -nilfunc -printf -shift -stdmethods -structtag -tests -unmarshal -unreachable -unsafeptr -unusedresult ./...
	test -z $(gofmt -l ./...)

.PHONY: fmt
fmt: ## Run syntax re-formatting (modify in place)
	go fmt ./...

.PHONY: check
check: ## Compile everything, checking syntax (does not output binaries)
	go build ./...

.env:
	if [ ! -f ".env" ]; then cp example.dev.env .env; fi

.PHONY: run-dev-api
run-dev-api: .env ## Runs api for local dev
	export API_LOCAL=true && go run cmd/api/main.go

.PHONY: build-api-image
build-api-image:
	docker build -t api -f cmd/api/Dockerfile .

.PHONY: run-api-image
run-api-image:
	docker run -p 9090:9090 api /api

.PHONY: run-dev-admin-manager
run-dev-admin-manager: .env ## Runs api for local dev
	export ADMIN_MANAGER_LOCAL=true && go run cmd/adminmanager/main.go

.PHONY: build-admin-manager-image
build-admin-manager-image:
	docker build -t admin-manager -f cmd/adminmanager/Dockerfile .

.PHONY: run-admin-manager-image
run-admin-manager-image:
	docker run -p 9091:9091 admin-manager /adminmanager
