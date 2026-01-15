# .RECIPEPREFIX := $(.RECIPEPREFIX)<space>
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

default: help

.PHONY: help
help: ## Shows help for each of the Makefile recipes.
	@echo 'Commands:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-32s\033[0m %s\n", $$1, $$2}'

.PHONY: init
init: # Initializes the project.
	cd ./apps/services && \
  make init && \
  cd ../../

.PHONY: ok
ok: ## Checks if the code is ok.
	cd ./apps/services && \
  make ok && \
  cd ../../apps/webclient && \
  deno lint && \
  deno fmt --check && \
  cd ../../

.PHONY: build
build: ## (Re)builds the containers.
	COMPOSE_BAKE=true docker compose build $(ARGS)

.PHONY: up
up: ## Starts the containers.
	COMPOSE_BAKE=true docker compose up --detach $(ARGS)

.PHONY: down
down: ## Destroys the containers.
	docker compose down $(ARGS)

.PHONY: start
start: ## Starts the containers.
	COMPOSE_BAKE=true docker compose start $(ARGS)

.PHONY: watch
watch: ## Starts the containers and watches for changes.
	COMPOSE_BAKE=true docker compose watch $(ARGS)

.PHONY: restart
restart: ## Restarts the containers.
	COMPOSE_BAKE=true docker compose restart $(ARGS)

.PHONY: stop
stop: ## Stops the containers.
	docker compose stop $(ARGS)

.PHONY: logs
logs: ## Shows the logs of the containers.
	docker compose logs --follow $(ARGS)

.PHONY: cli
cli: ## Opens a shell in the containers.
	COMPOSE_BAKE=true docker compose exec services bash $(ARGS)

%:
	@:
