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

.PHONY: run-dev-admin-manager
run-dev-admin-manager: .env ## Runs api for local dev
	export ADMIN_MANAGER_LOCAL=true && go run cmd/adminmanager/main.go
