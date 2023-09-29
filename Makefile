include .env
export

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: mod
mod: ## go mod tidy & go mod vendor
	@go mod tidy
	@go mod vendor

.PHONY: update
update: ## go modules update
	@go get -u -t ./...
	@go mod tidy
	@go mod vendor

.PHONY: up
up: ## docker compose up
	@docker compose --project-name gosql --file ./.docker/compose.yaml up -d

.PHONY: down
down: ## docker compose down
	@docker compose --project-name gosql down --volumes

.PHONY: balus
balus: ## destroy everything about docker. (containers, images, volumes, networks.)
	@docker compose --project-name gosql down --rmi all --volumes

.PHONY: migrate
migrate:
	@(cd schema && bun run prisma db push)

.PHONY: psql
psql:
	@docker exec -it postgres psql -U postgres

.PHONY: prisma
prisma:
	@(cd schema && bun run prisma studio)

.PHONY: prismafmt
prismafmt:
	@(cd schema && bun run prisma format)
