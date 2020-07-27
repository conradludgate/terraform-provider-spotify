.PHONY: help build

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:
	go build -o terraform-provider-spotify
	mkdir -p ~/.terraform.d/plugins
	mv terraform-provider-spotify ~/.terraform.d/plugins