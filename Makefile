# note: call scripts from /scripts
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

all: ## Compile server stack and binary client for Linux et Windows systems
	bash scripts/build.sh

linux: ## Compile binary client for Linux system
	bash scripts/build.sh linux

linux_test: ## Compile binary client test for Linux system
	bash scripts/build.sh linux_test

windows: ## Compile binary client for Windows system
	bash scripts/build.sh windows

server: ## Compile Autorace's server stack
	bash scripts/build.sh server

run: ## Start Autorace docker compose stack
	docker-compose --file deployments/docker-compose.yaml up -d

down: ## Stop Autorace docker compose stack
	bash scripts/build.sh down

clean: ## Delete docker build images and server stack
	bash scripts/build.sh clean

clean_linux: ## Delete docker build image for linux release
	bash scripts/build.sh clean linux

clean_linux_test: ## Delete docker build image for linux test release
	bash scripts/build.sh clean linux_test

clean_windows: ## Delete docker build image for windows release
	bash scripts/build.sh clean windows

clean_server: ## Delete docker image for server stack
	bash scripts/build.sh clean server
