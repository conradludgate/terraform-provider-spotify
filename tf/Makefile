.PHONY: help build install

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the go binary
	mkdir build
	go build -o build/terraform-provider-spotify

install: build ## Install the provider binary into the terraform plugins dir
	mkdir -p ~/.terraform.d/plugins
	cp build/terraform-provider-spotify ~/.terraform.d/plugins