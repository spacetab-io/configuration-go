deps: ## Get the dependencies
	@go mod vendor

## Test stuff
get_lint_config:
	@[ -f ./.golangci.yml ] && echo ".golangci.yml exists" || ( echo "getting .golangci.yml" && curl -O https://raw.githubusercontent.com/spacetab-io/docker-images-golang/master/linter/.golangci.yml )
.PHONY: get_lint_config

lint: get_lint_config
	golangci-lint run -v
.PHONY: lint

# ----
## TEST stuff start

test-unit:
	go test -cover -race -count=1 -timeout 1s -coverprofile=c.out -v ./... && go tool cover -func=c.out
.PHONY: test-unit

coverage-html:
	go tool cover -html=c.out -o coverage.html
.PHONY: coverage-html

test: deps test-unit coverage-html
.PHONY: test

## TEST stuff end
# ----

tests_in_docker: ## Testing code with unit tests in docker container
	docker run --rm -v $(shell pwd):/app -i spacetabio/docker-test-golang:1.14-1.0.2 make test
.PHONY: tests_in_docker