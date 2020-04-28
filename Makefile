deps: ## Get the dependencies
	@go mod vendor

## Test stuff
get_lint_config:
	@[ -f ./.golangci.yml ] && echo ".golangci.yml exists" || ( echo "getting .golangci.yml" && curl -O https://raw.githubusercontent.com/spacetab-io/docker-images-golang/master/linter/.golangci.yml )
.PHONY: get_lint_config

lint: get_lint_config
	golangci-lint run -v
.PHONY: lint

test-unit:
	go test ./... --race --cover -count=1 -timeout 1s -coverprofile=c.out -v
.PHONY: test-unit

coverage-html:
	go tool cover -html=c.out -o coverage.html
.PHONE: coverage-html

test: deps test-unit coverage-html
.PHONY: test

tests_in_docker: ## Testing code with unit tests in docker container
	docker run --rm -v $(shell pwd):/app -i spacetabio/docker-test-golang:1.14-1.0.2 make test
.PHONY: tests_in_docker