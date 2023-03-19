# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY:
gomock:
	@echo "setting up gomock..."
	@go install github.com/golang/mock/mockgen@latest

.PHONY:
generate: gomock
	@echo "generating code..."
	@go generate ./...

.PHONY:
tests: mocks
	@echo "running tests..."
	@go test ./...